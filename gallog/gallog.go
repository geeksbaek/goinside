package gallog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/geeksbaek/goinside"
)

// Session 구조체는 갤로그에 접속하기 위한 세션을 표현합니다.
type Session struct {
	id      string
	pw      string
	cookies []*http.Cookie
	*goinside.MemberSessionDetail
}

// Login 함수는 로그인 된 갤로그 세션을 반환합니다.
func Login(id, pw string) (s *Session, err error) {
	form := _Form(map[string]string{
		"s_url":    "http://www.dcinside.com/",
		"ssl":      "Y",
		"user_id":  id,
		"password": pw,
	})
	resp, err := do("POST", desktopLoginURL, nil, form, desktopRequestHeader)
	if err != nil {
		return
	}
	ms, err := goinside.Login(id, pw)
	if err != nil {
		return
	}
	s = &Session{
		id:                  id,
		pw:                  pw,
		cookies:             resp.Cookies(),
		MemberSessionDetail: ms.MemberSessionDetail,
	}
	return
}

// Logout 함수는 갤로그 세션을 종료합니다.
func (s *Session) Logout() (err error) {
	_, err = do("GET", desktopLogoutURL, s.cookies, nil, desktopRequestHeader)
	return
}

// ArticleMicroInfo 구조체는 글 삭제를 위해 필요한 최소한의 정보를 표현합니다.
type ArticleMicroInfo struct {
	gid, cid, pno, logNo, mode string
}

func (a *ArticleMicroInfo) delete(s *Session) {
	// first, get gall dcinside
	gallID, _, secretKey, secretVal, err := s._FetchDetails(a)
	if err != nil && err != errParseGallogSecreyKeyPair {
		return
	}

	// second, delete real article
	api(articleDeleteAPI, _Form(map[string]string{
		"app_id":  goinside.AppID,
		"user_id": s.UserID,
		"no":      a.pno,
		"id":      gallID,
		"mode":    "board_del",
	}))

	// third, delete article log
	form := _Form(map[string]string{
		"rb":      "",
		"dTp":     "1",
		"gid":     a.gid,
		"cid":     a.cid,
		"page":    "",
		"pno":     a.pno,
		"no":      a.pno,
		"logNo":   a.logNo,
		"id":      gallID,
		"nate":    "",
		secretKey: secretVal,
	})
	do("POST", deleteArticleLogURL, s.cookies, form, gallogRequestHeader)
}

// FetchAllArticle 함수는 해당 갤로그의 모든 글을 가져옵니다.
func (s *Session) FetchAllArticle() []*ArticleMicroInfo {
	ds := s.concurrencyFetch(func(URL string) (ds []deletable) {
		doc, err := _NewGallogDocument(s, URL)
		if err != nil {
			log.Fatal(err)
		}
		q := `td[valign='top'] td[colspan='2'] table tr:not(:first-child)`
		doc.Find(q).Each(func(i int, s *goquery.Selection) {
			data, _ := s.Find(`img`).Attr(`onclick`)
			ami, err := _ParseGallogArticleURL(data)
			if err != nil {
				return
			}
			ds = append(ds, deletable(ami))
		})
		return
	}, _GallogArticlePageURL)
	as := []*ArticleMicroInfo{}
	for _, d := range ds {
		as = append(as, d.(*ArticleMicroInfo))
	}
	return as
}

// CommentMicroInfo 구조체는 댓글 삭제를 위해 필요한 최소한의 정보를 표현합니다.
type CommentMicroInfo struct {
	gid, no, cno, logNo string
}

func (c *CommentMicroInfo) delete(s *Session) {
	gallID, cid, _, _, err := s._FetchDetails(c)
	if err != nil && err != errParseGallogSecreyKeyPair {
		return
	}

	api(commentDeleteAPI, _Form(map[string]string{
		"app_id":     goinside.AppID,
		"user_id":    s.UserID,
		"no":         c.no,
		"id":         gallID,
		"comment_no": c.cno,
		"mode":       "comment_del",
	}))

	form := _Form(map[string]string{
		"rb":    "",
		"dTp":   "1",
		"gid":   c.gid,
		"cid":   cid,
		"page":  "",
		"pno":   "",
		"no":    c.no,
		"c_no":  c.cno,
		"logNo": c.logNo,
		"id":    gallID,
		"nate":  "",
		"MTg=":  "MTg=",
	})
	fmt.Println(form)
	do("POST", deleteCommentLogURL, s.cookies, form, gallogRequestHeader)
}

// FetchAllComment 함수는 해당 갤로그의 모든 댓글을 가져옵니다.
func (s *Session) FetchAllComment() []*CommentMicroInfo {
	ds := s.concurrencyFetch(func(URL string) (ds []deletable) {
		doc, err := _NewGallogDocument(s, URL)
		if err != nil {
			log.Fatal(err)
		}
		q := `td[colspan='2'][align='center'] td[colspan='2'] table tr:not(:first-child)`
		doc.Find(q).Each(func(i int, s *goquery.Selection) {
			data, _ := s.Find(`td[width='22'] span`).Attr(`onclick`)
			cmi, err := _ParseGallogCommentURL(data)
			if err != nil {
				return
			}
			ds = append(ds, deletable(cmi))
		})
		return
	}, _GallogCommentPageURL)
	cs := []*CommentMicroInfo{}
	for _, d := range ds {
		cs = append(cs, d.(*CommentMicroInfo))
	}
	return cs
}

const (
	maxConcurrentRequestCount = 10
)

func (s *Session) concurrencyFetch(fetcher func(string) []deletable, URL func(string, int) string) (ds []deletable) {
	ds = []deletable{}
Loop:
	for i := 1; ; i += maxConcurrentRequestCount {
		tempDSS := make([][]deletable, maxConcurrentRequestCount)
		wg := new(sync.WaitGroup)
		wg.Add(maxConcurrentRequestCount)
		for j := i; j < i+maxConcurrentRequestCount; j++ {
			page := j
			URL := URL(s.id, j)
			go func() {
				defer wg.Done()
				tempDSS[page-1] = fetcher(URL)
			}()
		}
		wg.Wait()

		for _, tempDS := range tempDSS {
			if len(tempDS) == 0 {
				break Loop
			}
			ds = append(ds, tempDS...)
		}
	}
	return
}

type deletable interface {
	delete(*Session)
}

// DeleteArticle 함수는 삭제 가능한 객체를 전달받아 모두 삭제합니다.
func (s *Session) DeleteArticle(as []*ArticleMicroInfo) {
	wg := new(sync.WaitGroup)
	wg.Add(len(as))
	for _, a := range as {
		go func(a *ArticleMicroInfo) {
			defer wg.Done()
			a.delete(s)
		}(a)
	}
	wg.Wait()
}

// DeleteComment 함수는 삭제 가능한 객체를 전달받아 모두 삭제합니다.
func (s *Session) DeleteComment(cs []*CommentMicroInfo) {
	wg := new(sync.WaitGroup)
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c *CommentMicroInfo) {
			defer wg.Done()
			c.delete(s)
		}(c)
	}
	wg.Wait()
}

func (s *Session) _FetchDetails(i interface{}) (id, cid, key, val string, err error) {
	var URL string
	switch t := i.(type) {
	case *ArticleMicroInfo:
		URL = _GallogArticleDetailURL(t)
	case *CommentMicroInfo:
		URL = _GallogCommentDetailURL(t)
	default:
		err = errors.New("unknown type")
		return
	}

	resp, err := do("GET", URL, s.cookies, nil, gallogRequestHeader)
	if err != nil {
		return
	}
	_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	body := string(_body)
	id, err = _ParseGallogGallID(body)
	if err != nil {
		return
	}
	cid, err = _ParseGallogCID(body)
	if err != nil {
		return
	}
	key, val, err = _ParseGallogSecretKeyPair(body)
	if err != nil {
		return
	}
	return
}
