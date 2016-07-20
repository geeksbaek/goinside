package goinside

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
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
		s = fmt.Sprintf("%04d.%02d.%02d %v", now.Year(), int(now.Month()), now.Day(), s)
	}
	if t, err := time.Parse("2006.01.02", s); err == nil {
		return &t
	}
	if t, err := time.Parse("2006.01.02 15:04", s); err == nil {
		return &t
	}
	return nil
}

func convertToMobileDcinside(URL string) string {
	if matched := desktopURLRe.FindStringSubmatch(URL); len(matched) > 0 {
		switch {
		case len(matched) == 2 || (len(matched) >= 3 && matched[2] == ""):
			return fmt.Sprintf("http://m.dcinside.com/list.php?id=%s",
				matched[1])
		case len(matched) >= 3:
			return fmt.Sprintf("http://m.dcinside.com/view.php?id=%s&no=%s",
				matched[1], matched[2])
		}
	}
	return URL
}

func parseGallID(URL string) string {
	if matched := desktopURLRe.FindStringSubmatch(URL); len(matched) > 2 {
		return strings.TrimSpace(matched[1])
	}
	return ""
}

func trimContent(content string) string {
	out := ""
	for _, v := range strings.Split(content, "\n") {
		out += strings.TrimSpace(v)
	}
	return strings.TrimSpace(out)
}

func generateMD5() (ret string) {
	n := strconv.Itoa(rand.Intn(9999999-1000000) + 1000000)
	for _, v := range md5.Sum([]byte(n)) {
		ret += fmt.Sprintf("%02x", v&255)
	}
	return
}

func ocr(a *articleWriter, URL string) (string, error) {
	resp, err := a.getCaptcha(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonStr := []byte(fmt.Sprintf(`{"base64": "%s", "trim": "\n", "whitelist": "123456789"}`, base64.StdEncoding.EncodeToString(body)))

	resp, _ = http.Post("http://192.168.99.100:8080/base64", "text/plain;charset=UTF-8", bytes.NewBuffer(jsonStr))

	body, _ = ioutil.ReadAll(resp.Body)
	var respJSON struct {
		Result string
	}
	json.Unmarshal(body, &respJSON)
	if respJSON.Result == "" {
		return "", errors.New("ocr: parsing fail")
	}
	return respJSON.Result, nil
}
