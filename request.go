package goinside

import (
	"io"
	"net/http"
)

const (
	UploadImageURL = "http://upload.dcinside.com/upload_imgfree_mobile.php"
	GWriteURL      = "http://upload.dcinside.com/g_write.php"
	OptionWriteURL = "http://m.dcinside.com/_option_write.php"
	RecommendURL   = "http://m.dcinside.com/_recommend_join.php"
	NorecommendURL = "http://m.dcinside.com/_nonrecommend_join.php"
	CommentURL     = "http://m.dcinside.com/_option_write.php"
	AccessTokenURL = "http://m.dcinside.com/_access_token.php"
	GallTotalURL   = "http://m.dcinside.com/category_gall_total.html"

	DefaultContentType = "application/x-www-form-urlencoded; charset=UTF-8"
)

var (
	defaultRequestHeader = map[string]string{
		"User-Agent":       "Linux Android",
		"Referer":          "http://m.dcinside.com",
		"X-Requested-With": "XMLHttpRequest",
	}
)

func (s *Session) post(URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return s.do("POST", URL, cookies, form, contentType)
}

func (s *Session) get(URL string) (*http.Response, error) {
	return s.do("GET", URL, nil, nil, DefaultContentType)
}

func (s *Session) do(method, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
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
