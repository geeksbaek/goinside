package goinside

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

var iconURLsMap = map[string]string{
	"ico_p_y": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_picture.png",
	"ico_t":   "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_text.png",
	"ico_p_c": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_picture_b.png",
	"ico_t_c": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_text_b.png",
	"ico_mv":  "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_movie.png",
	"ico_sc":  "http://nstatic.dcinside.com/dgn/gallery/images/update/sec_icon.png",
}

var gallogIconURLsMap = map[string]string{
	"fixed": "http://wstatic.dcinside.com/gallery/skin/gallog/g_default.gif",
	"flow":  "http://wstatic.dcinside.com/gallery/skin/gallog/g_fix.gif",
}

// func newDocument(url string, header map[string]string) *goquery.Document {
// 	doc, err := goquery.NewDocumentFromResponse(get(url, header))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return doc
// }

// func isMobileWeb(url string) bool {
// 	re := regexp.MustCompile(`http(s)?:\/\/m\.*`)
// 	return re.MatchString(url)
// }

// func isDcconURL(url string) bool {
// 	re := regexp.MustCompile(`^http:\/\/dcimg1.dcinside.com\/dccon\.php\?no=\w+$`)
// 	return re.MatchString(url)
// }

// func splitURL(url string) (string, string) {
// 	re := regexp.MustCompile(`(\w+)\/\?(.*)`)
// 	substr := re.FindStringSubmatch(url)
// 	return substr[1], substr[2]
// }

// func respToString(resp *http.Response) string {
// 	body, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return string(body)
// }

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
