package gallog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/geeksbaek/goinside"
)

const (
	maxConcurrentRequestCount = 100

	articlesQuery = `td[valign='top'] td[colspan='2'] table tr:not(:first-child)`
	articleQuery  = `img`
	commentsQuery = `td[colspan='2'][align='center'] td[colspan='2'] table tr:not(:first-child)`
	commentQuery  = `td[width='22'] span`

	gallogURLFormat        = "http://gallog.dcinside.com/inc/_mainGallog.php?gid=%v&page=%v&rpage=%v"
	articleDetailURLFormat = "http://gallog.dcinside.com/inc/_deleteLog.php?gid=%v&cid=%v&page=&pno=%v&logNo=%v&mode=%v"
	commentDetailURLFormat = "http://gallog.dcinside.com/inc/_deleteLogRep.php?gid=%v&cid=&id=&no=%v&c_no=%v&logNo=%v&rpage="
)

var (
	gallogArticleURLRe = regexp.MustCompile(`gid=([^&]+)&cid=([^&]+)&page=.*&pno=([^&]+)&logNo=([^&]+)&mode=([^&']+)`)
	gallogCommentURLRe = regexp.MustCompile(`gid=([^&]+)&cid=.*&id=&no=([^&]+)&c_no=([^&]+)&logNo=([^&]+)&rpage=.*`)
	gallIDRe           = regexp.MustCompile(`<INPUT TYPE="hidden" NAME="id" value=(?:"|')(.+)(?:"|')>`)
	secretRe           = regexp.MustCompile(`<INPUT TYPE="hidden" NAME=".*" id=(?:"|')([^'"]+)(?:"|') value=(?:"|')([^'"]{4,})(?:"|')>`)
	cidRe              = regexp.MustCompile(`<INPUT TYPE="hidden" NAME="cid" value="([^"]+)">`)
)

type Session struct {
	id      string
	pw      string
	cookies []*http.Cookie
	*goinside.MemberSessionDetail
}

func Login(id, pw string) (s *Session, err error) {
	form := makeForm(map[string]string{
		"s_url":    "http://www.dcinside.com/",
		"ssl":      "Y",
		"user_id":  id,
		"password": pw,
	})
	resp := do("POST", desktopLoginURL, nil, form, desktopRequestHeader)
	ms, err := goinside.Login(id, pw)
	if err != nil {
		return
	}
	s = &Session{id, pw, resp.Cookies(), ms.MemberSessionDetail}
	return
}

func (s *Session) Logout() (err error) {
	do("GET", desktopLogoutURL, s.cookies, nil, desktopRequestHeader)
	return
}

type articleMicroInfo struct {
	gid, cid, pno, logNo, mode string
}

type commentMicroInfo struct {
	gid, no, cno, logNo string
}

type GallogDataSet struct {
	As []*articleMicroInfo
	Cs []*commentMicroInfo
}

func parseArticles(doc *goquery.Document) (as []*articleMicroInfo) {
	as = []*articleMicroInfo{}
	doc.Find(articlesQuery).Each(func(i int, s *goquery.Selection) {
		data, _ := s.Find(articleQuery).Attr(`onclick`)
		if data != "" {
			as = append(as, articleURLToArticleMicroInfo(data))
		}
	})
	return
}

func articleURLToArticleMicroInfo(URL string) *articleMicroInfo {
	matched := gallogArticleURLRe.FindStringSubmatch(URL)
	return &articleMicroInfo{matched[1], matched[2], matched[3], matched[4], matched[5]}
}

func parseComments(doc *goquery.Document) (cs []*commentMicroInfo) {
	cs = []*commentMicroInfo{}
	doc.Find(commentsQuery).Each(func(i int, s *goquery.Selection) {
		data, _ := s.Find(commentQuery).Attr(`onclick`)
		if data != "" {
			cs = append(cs, commentURLToCommentMicroInfo(data))
		}
	})
	return
}

func commentURLToCommentMicroInfo(URL string) *commentMicroInfo {
	matched := gallogCommentURLRe.FindStringSubmatch(URL)
	return &commentMicroInfo{matched[1], matched[2], matched[3], matched[4]}
}

func gallogURL(gid string, page int) string {
	return fmt.Sprintf(gallogURLFormat, gid, page, page)
}

func newGallogDocument(s *Session, URL string) *goquery.Document {
	resp := do("GET", URL, s.cookies, nil, gallogRequestHeader)
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (s *Session) FetchAll() (data *GallogDataSet) {
	max := maxConcurrentRequestCount
	data = &GallogDataSet{[]*articleMicroInfo{}, []*commentMicroInfo{}}

	// maxConcurrentRequestCount 값만큼 동시에 수행한다.
	for i := 1; ; i += max {
		tempASs := make([][]*articleMicroInfo, max)
		tempCSs := make([][]*commentMicroInfo, max)

		// fetching
		wg := new(sync.WaitGroup)
		wg.Add(max)
		for page := i; page < max+i; page++ {
			URL := gallogURL(s.id, page)
			index := page - i
			go func() {
				defer wg.Done()
				doc := newGallogDocument(s, URL)
				tempASs[index] = parseArticles(doc)
				tempCSs[index] = parseComments(doc)
			}()
		}
		wg.Wait()

		// check end of page and append to data
		articleDone, commentDone := false, false
		for _, tempAS := range tempASs {
			if len(tempAS) == 0 {
				articleDone = true
				break
			}
			data.As = append(data.As, tempAS...)
		}
		for _, tempCS := range tempCSs {
			if len(tempCS) == 0 {
				commentDone = true
				break
			}
			data.Cs = append(data.Cs, tempCS...)
		}
		if articleDone && commentDone {
			break
		}
	}
	return
}

func (s *Session) DeleteAll(data *GallogDataSet, cb func(i, n int)) {
	max := maxConcurrentRequestCount
	wg := new(sync.WaitGroup)

	progressCh := make(chan struct{})
	doneCh := make(chan struct{})
	go func() {
		i := 1
		n := len(data.As) + len(data.Cs)
		for _ = range progressCh {
			cb(i, n)
			i++
		}
		close(doneCh)
	}()

	for i, a := range data.As {
		wg.Add(1)
		go func(a *articleMicroInfo) {
			defer wg.Done()
			a.delete(s)
			progressCh <- struct{}{}
		}(a)
		if i%max == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	for i, c := range data.Cs {
		wg.Add(1)
		go func(c *commentMicroInfo) {
			defer wg.Done()
			c.delete(s)
			progressCh <- struct{}{}
		}(c)
		if i%max == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	close(progressCh)
	<-doneCh
}

func (a *articleMicroInfo) delete(s *Session) {
	gallID, _, key, value := s.fetchDetail(a)

	deleteArticleForm := makeForm(map[string]string{
		"app_id":  goinside.AppID,
		"user_id": s.UserID,
		"no":      a.pno,
		"id":      gallID,
		"mode":    "board_del",
	})
	deleteArticleLogForm := makeForm(map[string]string{
		"dTp":   "1",
		"gid":   a.gid,
		"cid":   a.cid,
		"pno":   a.pno,
		"no":    a.pno,
		"logNo": a.logNo,
		"id":    gallID,
		key:     value,
		// "rb":    "",
		// "page":  "",
		// "nate":  "",
	})

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		api(deleteArticleAPI, deleteArticleForm)
	}()
	go func() {
		defer wg.Done()
		do("POST", deleteArticleLogURL, s.cookies, deleteArticleLogForm, gallogRequestHeader)
	}()
	wg.Wait()
}

func (c *commentMicroInfo) delete(s *Session) {
	gallID, cid, key, value := s.fetchDetail(c)

	deleteCommentForm := makeForm(map[string]string{
		"app_id":     goinside.AppID,
		"user_id":    s.UserID,
		"no":         c.no,
		"id":         gallID,
		"comment_no": c.cno,
		"mode":       "comment_del",
	})
	deleteCommentLogForm := makeForm(map[string]string{
		"dTp":   "1",
		"gid":   c.gid,
		"cid":   cid,
		"no":    c.no,
		"c_no":  c.cno,
		"logNo": c.logNo,
		"id":    gallID,
		key:     value,
		// "rb":    "",
		// "page":  "",
		// "pno":   "",
		// "nate":  "",
	})

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		api(deleteCommentAPI, deleteCommentForm)
	}()
	go func() {
		defer wg.Done()
		do("POST", deleteCommentLogURL, s.cookies, deleteCommentLogForm, gallogRequestHeader)
	}()
	wg.Wait()
}

type detailer interface {
	fetchDetail() string
}

func (a *articleMicroInfo) fetchDetail() string {
	return fmt.Sprintf(articleDetailURLFormat, a.gid, a.cid, a.pno, a.logNo, a.mode)
}

func (c *commentMicroInfo) fetchDetail() string {
	return fmt.Sprintf(commentDetailURLFormat, c.gid, c.no, c.cno, c.logNo)

}

func (s *Session) fetchDetail(d detailer) (gallID, cid, key, val string) {
	resp := do("GET", d.fetchDetail(), s.cookies, nil, gallogRequestHeader)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// gall ID
	if matched := gallIDRe.FindSubmatch(body); len(matched) == 2 {
		gallID = string(matched[1])
	}
	// secret key, value
	if matched := secretRe.FindSubmatch(body); len(matched) == 3 {
		key, val = string(matched[1]), string(matched[2])
	}
	// cid
	if matched := cidRe.FindSubmatch(body); len(matched) == 2 {
		cid = string(matched[1])
	}
	return
}

func makeForm(m map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range m {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}
