package goinside

import (
	"io"
	"net/http"
)

var (
	optionWrite = "http://m.dcinside.com/_option_write.php"
	uploadImage = "http://upload.dcinside.com/upload_imgfree_mobile.php"
	gWrite      = "http://upload.dcinside.com/g_write.php"
	recommend   = "http://m.dcinside.com/_recommend_join.php"
	comment     = "http://m.dcinside.com/_option_write.php"
	accessToken = "http://m.dcinside.com/_access_token.php"

	defaultContentType   = "application/x-www-form-urlencoded; charset=UTF-8"
	defaultRequestHeader = map[string]string{
		"User-Agent":       "Linux Android",
		"Referer":          "http://m.dcinside.com",
		"X-Requested-With": "XMLHttpRequest",
	}
)

func (a *Auth) post(URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return a.do("POST", URL, cookies, form, contentType)
}

func (a *Auth) get(URL string) (*http.Response, error) {
	return a.do("GET", URL, nil, nil, defaultContentType)
}

func (a *Auth) do(method, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
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
		if a.transport != nil {
			return &http.Client{Transport: a.transport}
		}
		return &http.Client{}
	}()
	return client.Do(req)
}
