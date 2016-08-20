package goinside

import (
	"io"
	"net/http"
)

// urls
const (
	gallerysURL    = "http://m.dcinside.com/category_gall_total.html"
	commentMoreURL = "http://m.dcinside.com/comment_more_new.php"
)

type session interface {
	articleWriteForm(*ArticleDraft) (form io.Reader, contentType string)
	articleDeleteForm(*Article) (form io.Reader, contentType string)
	commentWriteForm(*CommentDraft) (form io.Reader, contentType string)
	commentDeleteForm(*Comment) (form io.Reader, contentType string)
	actionForm(*Article) (form io.Reader, contentType string)
	reportForm(string, string) (form io.Reader, contentType string)
	connector
}

type dcinsideAPI string

func (api dcinsideAPI) post(c connector, form io.Reader, contentType string) (*http.Response, error) {
	return do(c, "POST", string(api), nil, form, contentType, apiRequestHeader)
}

func (api dcinsideAPI) get() (*http.Response, error) {
	return do(&GuestSession{}, "GET", string(api), nil, nil, defaultContentType, apiRequestHeader)
}

// apis
const (
	AppID = "blM1T09mWjRhQXlZbE1ML21xbkM3QT09"

	loginAPI             dcinsideAPI = "https://dcid.dcinside.com/join/mobile_app_login.php"
	writeArticleAPI      dcinsideAPI = "http://upload.dcinside.com/_app_write_api.php"
	deleteArticleAPI     dcinsideAPI = "http://m.dcinside.com/api/gall_del.php"
	writeCommentAPI      dcinsideAPI = "http://m.dcinside.com/api/comment_ok.php"
	deleteCommentAPI     dcinsideAPI = "http://m.dcinside.com/api/comment_del.php"
	recommendUpAPI       dcinsideAPI = "http://m.dcinside.com/api/_recommend_up.php"
	recommendDownAPI     dcinsideAPI = "http://m.dcinside.com/api/_recommend_down.php"
	reportAPI            dcinsideAPI = "http://m.dcinside.com/api/report_upload.php"
	redirectAPI          dcinsideAPI = "http://m.dcinside.com/api/redirect.php"
	readListAPI          dcinsideAPI = "http://m.dcinside.com/api/gall_list.php"
	readArticleAPI       dcinsideAPI = "http://m.dcinside.com/api/view2.php"
	readArticleDetailAPI dcinsideAPI = "http://m.dcinside.com/api/gall_view.php"
	readArticleImageAPI  dcinsideAPI = "http://m.dcinside.com/api/view_img.php"
	readCommentAPI       dcinsideAPI = "http://m.dcinside.com/api/comment.php"
	majorGalleryListAPI  dcinsideAPI = "http://json.dcinside.com/App/gall_name.php"
	minorGalleryListAPI  dcinsideAPI = "http://json.dcinside.com/App/gall_name_sub.php"
)

// content types
const (
	defaultContentType    = "application/x-www-form-urlencoded; charset=UTF-8"
	nonCharsetContentType = "application/x-www-form-urlencoded"
)

var (
	apiRequestHeader = map[string]string{
		"User-Agent": "dcinside.app",
		"Referer":    "http://m.dcinside.com",
		"Host":       "m.dcinside.com",
	}
	mobileRequestHeader = map[string]string{
		"User-Agent": "Linux Android",
		"Referer":    "http://m.dcinside.com",
	}
)

type connector interface {
	Connection() *Connection
}

func post(c connector, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return do(c, "POST", URL, cookies, form, contentType, mobileRequestHeader)
}

func get(c connector, URL string) (*http.Response, error) {
	return do(c, "GET", URL, nil, nil, defaultContentType, mobileRequestHeader)
}

func do(c connector, method, URL string, cookies []*http.Cookie, form io.Reader, contentType string, requestHeader map[string]string) (*http.Response, error) {
	URL = ToMobileURL(URL)
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range requestHeader {
		req.Header.Set(k, v)
	}
	client := func() *http.Client {
		proxy := c.Connection().proxy
		if proxy != nil {
			return &http.Client{Transport: &http.Transport{Proxy: proxy}}
		}
		return &http.Client{}
	}()
	if c.Connection().timeout != 0 {
		client.Timeout = c.Connection().timeout
	}
	return client.Do(req)
}
