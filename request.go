package goinside

import (
	"io"
	"net/http"
)

const (
	gallerysURL    = "http://m.dcinside.com/category_gall_total.html"
	commentMoreURL = "http://m.dcinside.com/comment_more_new.php"

	appID = "blM1T09mWjRhQXlZbE1ML21xbkM3QT09"

	loginAPI         = "https://dcid.dcinside.com/join/mobile_app_login.php"
	articleWriteAPI  = "http://upload.dcinside.com/_app_write_api.php"
	articleDeleteAPI = "http://m.dcinside.com/api/gall_del.php"
	commentWriteAPI  = "http://m.dcinside.com/api/comment_ok.php"
	commentDeleteAPI = "http://m.dcinside.com/api/comment_del.php"
	recommendUpAPI   = "http://m.dcinside.com/api/_recommend_up.php"
	recommendDownAPI = "http://m.dcinside.com/api/_recommend_down.php"
	reportAPI        = "http://m.dcinside.com/api/report_upload.php"

	defaultContentType    = "application/x-www-form-urlencoded; charset=UTF-8"
	nonCharsetContentType = "application/x-www-form-urlencoded"
)

var (
	apiRequestHeader = map[string]string{
		"User-Agent": "dcinside.app",
		"Referer":    "http://m.dcinside.com",
		"Host":       "m.dcinside.com",
	}
	defaultRequestHeader = map[string]string{
		"User-Agent":       "Linux Android",
		"Referer":          "http://m.dcinside.com",
		"X-Requested-With": "XMLHttpRequest",
	}
)

func post(s Session, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return do(s, "POST", URL, cookies, form, contentType, defaultRequestHeader)
}

func get(s Session, URL string) (*http.Response, error) {
	return do(s, "GET", URL, nil, nil, defaultContentType, defaultRequestHeader)
}

func api(s Session, URL string, form io.Reader, contentType string) (*http.Response, error) {
	return do(s, "POST", URL, nil, form, contentType, apiRequestHeader)
}

func do(s Session, method, URL string, cookies []*http.Cookie, form io.Reader, contentType string, requestHeader map[string]string) (*http.Response, error) {
	URL = _MobileURL(URL)
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		return nil, err
	}
	for k, v := range requestHeader {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)
	client := func() *http.Client {
		proxy := s.connection().proxy
		if proxy != nil {
			return &http.Client{Transport: &http.Transport{Proxy: proxy}}
		}
		return &http.Client{}
	}()
	if s.connection().timeout != 0 {
		client.Timeout = s.connection().timeout
	}
	return client.Do(req)
}
