package goinside

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	iconURLPrefix       = "http://nstatic.dcinside.com/dgn/gallery/images/update/"
	gallogIconURLPrefix = "http://wstatic.dcinside.com/gallery/skin/gallog/"
)

var (
	iconURLMap = map[string]string{
		"ico_p_y": iconURLPrefix + "icon_picture.png",
		"ico_t":   iconURLPrefix + "icon_text.png",
		"ico_p_c": iconURLPrefix + "icon_picture_b.png",
		"ico_t_c": iconURLPrefix + "icon_text_b.png",
		"ico_mv":  iconURLPrefix + "icon_movie.png",
		"ico_sc":  iconURLPrefix + "sec_icon.png",
	}
	gallogIconURLMap = map[string]string{
		"fixed": gallogIconURLPrefix + "g_default.gif",
		"flow":  gallogIconURLPrefix + "g_fix.gif",
	}
)

func form(m map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range m {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}

func cookies(m map[string]string) []*http.Cookie {
	cookies := []*http.Cookie{}
	for k, v := range m {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	return cookies
}

func newMobileDoc(URL string) (*goquery.Document, error) {
	resp, err := (&Session{}).get(URL)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}

func strToTime(s string) *time.Time {
	if len(s) <= 5 {
		now := time.Now()
		s = fmt.Sprintf("%v.%v.%v %v", now.Year(), now.Month(), now.Day(), s)
	}
	if t, err := time.Parse("2006.June.02 3:04", s); err == nil {
		return &t
	}
	if t, err := time.Parse("2006.01.02", s); err == nil {
		return &t
	}
	if t, err := time.Parse("2006.01.02 3:04", s); err == nil {
		return &t
	}
	return nil
}
