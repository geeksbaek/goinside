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

	// AppID 는 디시인사이드 공식 API에 접근하기 위한 ID입니다.
	AppID = "blM1T09mWjRhQXlZbE1ML21xbkM3QT09"

	gallArticleWriteAPI = "http://upload.dcinside.com/_app_write_api.php"
	recommendUpAPI      = "http://m.dcinside.com/api/_recommend_up.php"
	recommendDownAPI    = "http://m.dcinside.com/api/_recommend_down.php"
	reportAPI           = "http://m.dcinside.com/api/report_upload.php"

	defaultContentType    = "application/x-www-form-urlencoded; charset=UTF-8"
	nonCharsetContentType = "application/x-www-form-urlencoded"
	imagePngContentType   = "image/png"
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
	desktopURLRe = regexp.MustCompile(`(?:http:\/\/)?gall\.dcinside\.com.*id=([^&]+)(?:&no=(\d+))?`)
)

func (s *Session) post(URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return s.do("POST", URL, cookies, form, contentType, defaultRequestHeader)
}

func (s *Session) get(URL string) (*http.Response, error) {
	return s.do("GET", URL, nil, nil, defaultContentType, defaultRequestHeader)
}

func (s *Session) api(URL string, form io.Reader, contentType string) (*http.Response, error) {
	return s.do("POST", URL, nil, form, contentType, apiRequestHeader)
}

func (s *Session) getCaptcha(URL string) (*http.Response, error) {
	return s.do("GET", URL, nil, nil, imagePngContentType, apiRequestHeader)
}

func (s *Session) do(method, URL string, cookies []*http.Cookie, form io.Reader, contentType string, requestHeader map[string]string) (*http.Response, error) {
	URL = convertToMobileDcinside(URL)
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		return nil, err
	}
	cookies = append(cookies, s.cookies...)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	for k, v := range requestHeader {
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
