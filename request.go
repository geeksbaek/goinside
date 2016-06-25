package goinside

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	uploadImageURL = "http://upload.dcinside.com/upload_imgfree_mobile.php"
	gWriteURL      = "http://upload.dcinside.com/g_write.php"
	optionWriteURL = "http://m.dcinside.com/_option_write.php"
	recommendURL   = "http://m.dcinside.com/_recommend_join.php"
	norecommendURL = "http://m.dcinside.com/_nonrecommend_join.php"
	commentURL     = "http://m.dcinside.com/_option_write.php"
	accessTokenURL = "http://m.dcinside.com/_access_token.php"
	gallTotalURL   = "http://m.dcinside.com/category_gall_total.html"
	commentMoreURL = "http://m.dcinside.com/comment_more_new.php"

	defaultContentType = "application/x-www-form-urlencoded; charset=UTF-8"
)

var (
	defaultRequestHeader = map[string]string{
		"User-Agent":       "Linux Android",
		"Referer":          "http://m.dcinside.com",
		"X-Requested-With": "XMLHttpRequest",
	}
	desktopURLRe = regexp.MustCompile(`(?:http:\/\/)?gall\.dcinside\.com.*id=([^&]+)&no=(\d+)`)
)

func (s *Session) post(URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return s.do("POST", URL, cookies, form, contentType)
}

func (s *Session) get(URL string) (*http.Response, error) {
	return s.do("GET", URL, nil, nil, defaultContentType)
}

func (s *Session) do(method, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	if matched := desktopURLRe.FindStringSubmatch(URL); len(matched) > 0 {
		switch {
		case len(matched) == 2:
			URL = fmt.Sprintf("http://m.dcinside.com/list.php?id=%s",
				matched[1])
		case len(matched) >= 3:
			URL = fmt.Sprintf("http://m.dcinside.com/view.php?id=%s&no=%s",
				matched[1], matched[2])
		}
	}
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	for k, v := range defaultRequestHeader {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)
	client := func() *http.Client {
		if s.transport != nil {
			return &http.Client{Transport: s.transport}
		}
		return &http.Client{}
	}()
	return client.Do(req)
}
