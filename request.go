package goinside

import (
	"io"
	"net/http"
	"regexp"
)

const (
	uploadImageURL  = "http://upload.dcinside.com/upload_imgfree_mobile.php"
	gWriteURL       = "http://upload.dcinside.com/g_write.php"
	optionWriteURL  = "http://m.dcinside.com/_option_write.php"
	recommendURL    = "http://m.dcinside.com/_recommend_join.php"
	norecommendURL  = "http://m.dcinside.com/_nonrecommend_join.php"
	commentURL      = "http://m.dcinside.com/_option_write.php"
	accessTokenURL  = "http://m.dcinside.com/_access_token.php"
	gallTotalURL    = "http://m.dcinside.com/category_gall_total.html"
	commentMoreURL  = "http://m.dcinside.com/comment_more_new.php"
	loginPageURL    = "http://m.dcinside.com/login.php"
	mobileLoginURL  = "https://dcid.dcinside.com/join/mobile_login_ok.php"
	logoutURL       = "http://m.dcinside.com/logout.php"
	gallogPrefixURL = "http://m.dcinside.com/gallog/home.php"

	defaultContentType    = "application/x-www-form-urlencoded; charset=UTF-8"
	nonCharsetContentType = "application/x-www-form-urlencoded"
)

var (
	defaultRequestHeader = map[string]string{
		"User-Agent":       "Linux Android",
		"Referer":          "http://m.dcinside.com",
		"X-Requested-With": "XMLHttpRequest",
	}
	desktopURLRe = regexp.MustCompile(`(?:http:\/\/)?gall\.dcinside\.com.*id=([^&]+)(?:&no=(\d+))?`)
)

func (s *Session) post(URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return s.do("POST", URL, cookies, form, contentType)
}

func (s *Session) get(URL string) (*http.Response, error) {
	return s.do("GET", URL, nil, nil, defaultContentType)
}

func (s *Session) do(method, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	URL = convertToMobileDcinside(URL)
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		return nil, err
	}
	cookies = append(cookies, s.cookies...)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	for k, v := range defaultRequestHeader {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)
	client := func() *http.Client {
		if s.proxy != nil {
			return &http.Client{Transport: &http.Transport{Proxy: s.proxy}}
		}
		return &http.Client{}
	}()
	if s.timeout != 0 {
		client.Timeout = s.timeout
	}
	return client.Do(req)
}
