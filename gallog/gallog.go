package gallog

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/geeksbaek/goinside"
)

const (
	loginChkQuery = `input[type="hidden"]:nth-child(3)`
	itemsQuery    = `#container > article > div > section > div.gallog_cont > div > ul > li`

	gallogURLFormat        = "http://gallog.dcinside.com/%v"
	gallogArticleURLFormat = `http://gallog.dcinside.com/%v/posting?p=%v`
	gallogCommentURLFormat = `http://gallog.dcinside.com/%v/comment?p=%v`
)

// Session 구조체는 갤로그 세션 정보를 나타냅니다.
type Session struct {
	id      string
	pw      string
	cookies []*http.Cookie
	*goinside.MemberSessionDetail
}

// Login 함수는 해당 ID와 PASSWORD로 로그인한 뒤 해당 세션을 반환합니다.
func Login(id, pw string) (s *Session, err error) {
	loginPageResp := do("GET", desktopLoginPageURL, nil, nil, desktopRequestHeader)
	doc, err := goquery.NewDocumentFromResponse(loginPageResp)
	if err != nil {
		return
	}

	chk := doc.Find(loginChkQuery)
	chkName, _ := chk.Attr("name")
	chkValue, _ := chk.Attr("value")

	f := map[string]string{
		"s_url":    "//www.dcinside.com/",
		"ssl":      "Y",
		"user_id":  id,
		"password": pw,
	}
	f[chkName] = chkValue

	ssoResp := do("GET", desktopSSOIframeURL, nil, nil, ssoRequestHeader)

	form := makeForm(f)
	resp := do("POST", desktopLoginURL, ssoResp.Cookies(), form, desktopRequestHeader)

	ms, err := goinside.Login(id, pw)
	if err != nil {
		return
	}

	cookies := []*http.Cookie{}
	for _, v := range resp.Cookies() {
		if v.Value != "deleted" {
			cookies = append(cookies, v)
		}
	}

	myGallogURL := fmt.Sprintf(gallogURLFormat, id)
	resp = do("GET", myGallogURL, cookies, nil, nil)

	ci_c := getCI_CCookie(resp.Cookies())
	if ci_c != nil {
		cookies = append(cookies, ci_c)
	}

	s = &Session{id, pw, cookies, ms.MemberSessionDetail}
	return
}

// Logout 메소드는 해당 세션을 종료합니다.
func (s *Session) Logout() (err error) {
	do("GET", desktopLogoutURL, s.cookies, nil, desktopRequestHeader)
	s = nil
	return
}

type GallogItems []string

func parseItems(doc *goquery.Document) GallogItems {
	if doc == nil {
		return nil
	}
	items := GallogItems{}
	doc.Find(itemsQuery).Each(func(i int, s *goquery.Selection) {
		data, _ := s.Attr(`data-no`)
		if data != "" {
			items = append(items, data)
		}
	})
	return items
}

func newGallogDocument(s *Session, URL string) *goquery.Document {
	resp := do("GET", URL, s.cookies, nil, gallogRequestHeader)
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil
	}
	return doc
}

// FetchAll 메소드는 해당 세션의 갤로그에 존재하는 모든 데이터를 가져옵니다.
func (s *Session) FetchAll(max int, progressCh chan struct{}) (data GallogItems) {
	defer close(progressCh)
	data = GallogItems{}

	// max 값만큼 동시에 수행한다.
	for i := 1; ; i += max {
		tempArticleSlice := make([]GallogItems, max)
		tempCommentSlice := make([]GallogItems, max)

		// fetching
		wg := new(sync.WaitGroup)
		wg.Add(max)
		for page := i; page < max+i; page++ {
			articleURL := fmt.Sprintf(gallogArticleURLFormat, s.id, page)
			commentURL := fmt.Sprintf(gallogCommentURLFormat, s.id, page)
			index := page - i
			go func() {
				defer wg.Done()

				articleDoc := newGallogDocument(s, articleURL)
				tempArticleSlice[index] = parseItems(articleDoc)

				commentDoc := newGallogDocument(s, commentURL)
				tempCommentSlice[index] = parseItems(commentDoc)

				progressCh <- struct{}{}
			}()
		}
		wg.Wait()

		// check end of page and append to data
		articleDone, commentDone := false, false
		for _, tempArticles := range tempArticleSlice {
			if tempArticles == nil {
				continue
			}
			if len(tempArticles) == 0 {
				articleDone = true
				break
			}
			data = append(data, tempArticles...)
		}
		for _, tempComments := range tempCommentSlice {
			if tempComments == nil {
				continue
			}
			if len(tempComments) == 0 {
				commentDone = true
				break
			}
			data = append(data, tempComments...)
		}
		if articleDone && commentDone {
			break
		}
	}
	return
}

// DeleteAll 메소드는 해당 데이터를 모두 삭제합니다.
// 데이터 삭제 상황을 확인할 수 있도록 콜백 함수를 인자로 받습니다.
// 해당 콜백 함수는 삭제된 데이터 개수 i과 총 데이터 개수 n을 인자로 받습니다.
func (s *Session) DeleteAll(max int, data GallogItems, cb func(i, n int)) {
	wg := new(sync.WaitGroup)

	progressCh := make(chan struct{})
	doneCh := make(chan struct{})
	go func() {
		i := 1
		n := len(data)
		for range progressCh {
			cb(i, n)
			i++
		}
		close(doneCh)
	}()

	for i, no := range data {
		wg.Add(1)
		go func(no string) {
			defer wg.Done()
			s.deleteLog(no)
			progressCh <- struct{}{}
		}(no)
		if i%max == 0 {
			wg.Wait()
		}
	}
	wg.Wait()

	close(progressCh)
	<-doneCh
}

func (s *Session) deleteLog(no string) {
	fmt.Println(getCI_CValue(s.cookies))
	deleteForm := makeForm(map[string]string{
		"ci_t": getCI_CValue(s.cookies),
		"no":   no,
	})
	do("POST", deleteAPI, nil, deleteForm, gallogRequestHeader)
}

func getCI_CCookie(cookies []*http.Cookie) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == "ci_c" {
			return cookie
		}
	}
	return nil
}

func getCI_CValue(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "ci_c" {
			return cookie.Value
		}
	}
	return ""
}

func makeForm(m map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range m {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}
