package gallog

import (
	"io"
	"log"
	"net/http"
	"time"
)

// web urls
const (
	desktopLoginPageURL = "https://www.dcinside.com" // s_url 없으면 에러남
	desktopSSOIframeURL = "https://dcid.dcinside.com/join/sso_iframe.php"
	desktopLoginURL     = "https://dcid.dcinside.com/join/member_check.php"
	desktopLogoutURL    = "https://dcid.dcinside.com/join/logout.php"
	deleteArticleLogURL = "http://gallog.dcinside.com/inc/_deleteArticle.php"
	deleteCommentLogURL = "http://gallog.dcinside.com/inc/_deleteRepOk.php"
)

// apis
const (
	deleteArticleAPI = "http://m.dcinside.com/api/gall_del.php"
	deleteCommentAPI = "http://m.dcinside.com/api/comment_del.php"
)

// content types
const (
	nonCharsetContentType = "application/x-www-form-urlencoded"
)

var (
	apiRequestHeader = map[string]string{
		"User-Agent":   "dcinside.app",
		"Referer":      "http://m.dcinside.com",
		"Host":         "m.dcinside.com",
		"Content-Type": nonCharsetContentType,
	}
	gallogRequestHeader = map[string]string{
		"User-Agent":   "Mozilla/5.0",
		"Referer":      "http://gallog.dcinside.com",
		"Host":         "gallog.dcinside.com",
		"Content-Type": nonCharsetContentType,
	}
	desktopRequestHeader = map[string]string{
		"User-Agent":   "Mozilla/5.0",
		"Referer":      "https://www.dcinside.com",
		"Host":         "dcid.dcinside.com",
		"Content-Type": nonCharsetContentType,
	}
	ssoRequestHeader = map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "en-US,en;q=0.9,ko-KR;q=0.8,ko;q=0.7",
		"Cache-Control":             "no-cache",
		"Connection":                "keep-alive",
		"DNT":                       "1",
		"Host":                      "dcid.dcinside.com",
		"Pragma":                    "no-cache",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0",
	}
)

func api(URL string, form io.Reader) *http.Response {
	return do("POST", URL, nil, form, apiRequestHeader)
}

func do(method, URL string, cookies []*http.Cookie, form io.Reader, requestHeader map[string]string) *http.Response {
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		log.Fatal("http.NewRequest error :", err)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	for k, v := range requestHeader {
		req.Header.Set(k, v)
	}
	client := &http.Client{
		Timeout: time.Second * 3,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for i := 1; ; i++ {
		if resp, err := client.Do(req); err == nil {
			return resp
		}
		if i > 300 {
			// log.Fatal("디시인사이드 서버가 응답하지 않습니다.")
			return nil
		}
	}
}
