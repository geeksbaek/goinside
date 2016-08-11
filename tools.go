package goinside

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// regex
var (
	urlRe   = regexp.MustCompile(`id=([^&\s]+)(?:&no=([^&\s]+))?(?:&page=([^&\s]+))?`)
	imageRe = regexp.MustCompile(`img[^>]+src="([^"]+)"`)
)

// errors
var (
	errUnknownCause = errors.New("result false with empty cause")
)

// formatting
const (
	mobileGallURLFormat        = "http://m.dcinside.com/list.php?id=%v"
	mobileGallURLPageFormat    = "http://m.dcinside.com/list.php?id=%v&page=%v"
	mobileArticleURLFormat     = "http://m.dcinside.com/view.php?id=%v&no=%v"
	mobileArticleURLPageFormat = "http://m.dcinside.com/view.php?id=%v&no=%v&page=%v"
	mobileCommentPageFormat    = "http://m.dcinside.com/comment_more_new.php?id=%v&no=%v&com_page=%v"
)

var (
	ArticleIconURLMap = map[string]string{
		"ico_p_y": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_picture.png",
		"ico_t":   "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_text.png",
		"ico_p_c": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_picture_b.png",
		"ico_t_c": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_text_b.png",
		"ico_mv":  "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_movie.png",
		"ico_sc":  "http://nstatic.dcinside.com/dgn/gallery/images/update/sec_icon.png",
	}
	GallogIconURLMap = map[string]string{
		"fixed": "http://wstatic.dcinside.com/gallery/skin/gallog/g_default.gif",
		"flow":  "http://wstatic.dcinside.com/gallery/skin/gallog/g_fix.gif",
	}
)

// ToMobileURL 함수는 주어진 디시인사이드 URL을 모바일 URL로 포맷팅하여 반환합니다.
func ToMobileURL(URL string) string {
	if urlRe.MatchString(URL) {
		matched := urlRe.FindStringSubmatch(URL)
		id, number, page := matched[1], matched[2], matched[3]
		switch {
		case id != "" && number == "" && page != "": // id, page
			return fmt.Sprintf(mobileGallURLPageFormat, id, page)
		case id != "" && number != "" && page == "": // id, number
			return fmt.Sprintf(mobileArticleURLFormat, id, number)
		case id != "" && number == "" && page == "": // id
			return fmt.Sprintf(mobileGallURLFormat, id)
		case id != "" && number != "" && page != "": // id, number, page
			return fmt.Sprintf(mobileArticleURLPageFormat, id, number, page)
		}
	}
	return URL
}

func imageElements(body string) []string {
	images := []string{}
	matched := imageRe.FindAllStringSubmatch(body, -1)
	for _, v := range matched {
		if len(v) >= 2 {
			images = append(images, v[1])
		}
	}
	return images
}

func mobileCommentPageURL(gallID, number string, page int) string {
	return fmt.Sprintf(mobileCommentPageFormat, gallID, number, page)
}

func gallID(URL string) string {
	if urlRe.MatchString(URL) == false {
		return ""
	}
	return urlRe.FindStringSubmatch(URL)[1]
}

func articleNumber(URL string) string {
	if urlRe.MatchString(URL) == false {
		return ""
	}
	return urlRe.FindStringSubmatch(URL)[2]
}

func newMobileDocument(URL string) (*goquery.Document, error) {
	resp, err := get(&GuestSession{}, URL)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}

func timeFormatting(s string) time.Time {
	if len(s) <= 5 {
		now := time.Now()
		s = fmt.Sprintf("%04d.%02d.%02d %v", now.Year(), int(now.Month()), now.Day(), s)
	}
	if t, err := time.Parse("2006.01.02", s); err == nil {
		return t
	}
	if t, err := time.Parse("2006.01.02 15:04", s); err == nil {
		return t
	}
	return time.Time{}
}

func checkResponse(resp *http.Response) error {
	jsonResponse := &_JSONResponse{}
	if err := responseUnmarshal(jsonResponse, resp); err != nil {
		return err
	}
	return checkJSONResult(jsonResponse)
}

func responseUnmarshal(data interface{}, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body = []byte(strings.Trim(string(body), "[]"))
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

func checkJSONResult(jsonResponse *_JSONResponse) error {
	if jsonResponse.Result == false {
		if jsonResponse.Cause != "" {
			return errors.New(jsonResponse.Cause)
		}
		return errUnknownCause
	}
	return nil
}

func makeForm(m map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range m {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}

func multipartForm(m map[string]string, images ...string) (io.Reader, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if len(images) != 0 {
		multipartImages(w, images...)
		for i := range images {
			k := fmt.Sprintf("memo_block[%d]", i)
			m[k] = fmt.Sprintf("Dc_App_Img_%d", i+1)
		}
	}
	k := fmt.Sprintf("memo_block[%d]", len(images))
	m[k] = m["content"]
	delete(m, "content")
	multipartOthers(w, m)
	return &b, w.FormDataContentType()
}

func multipartImages(w *multipart.Writer, images ...string) {
	for i, image := range images {
		h := textproto.MIMEHeader{}
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="upload[%d]"; filename="%s"`, i, image))
		h.Set("Content-Type", "image/jpeg")
		fw, err := w.CreatePart(h)
		if err != nil {
			return
		}
		f, err := os.Open(image)
		if err != nil {
			return
		}
		if _, err = io.Copy(fw, f); err != nil {
			return
		}
	}
}

func multipartOthers(w *multipart.Writer, m map[string]string) {
	for k, v := range m {
		if fw, err := w.CreateFormField(k); err != nil {
			continue
		} else if _, err := fw.Write([]byte(v)); err != nil {
			continue
		}
	}
}
