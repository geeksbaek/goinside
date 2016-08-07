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
	urlRe           = regexp.MustCompile(`id=([^&]+)(?:&no=(\d+))?`)
	errUnknownCause = errors.New("result false with empty cause")
)

var (
	iconURLMap = map[string]string{
		"ico_p_y": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_picture.png",
		"ico_t":   "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_text.png",
		"ico_p_c": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_picture_b.png",
		"ico_t_c": "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_text_b.png",
		"ico_mv":  "http://nstatic.dcinside.com/dgn/gallery/images/update/icon_movie.png",
		"ico_sc":  "http://nstatic.dcinside.com/dgn/gallery/images/update/sec_icon.png",
	}
	gallogIconURLMap = map[string]string{
		"fixed": "http://wstatic.dcinside.com/gallery/skin/gallog/g_default.gif",
		"flow":  "http://wstatic.dcinside.com/gallery/skin/gallog/g_fix.gif",
	}
)

func _ParseGallID(URL string) string {
	if matched := urlRe.FindStringSubmatch(URL); len(matched) >= 2 {
		return strings.TrimSpace(matched[1])
	}
	return ""
}

func _ParseArticleNumber(URL string) string {
	if matched := urlRe.FindStringSubmatch(URL); len(matched) >= 3 {
		return strings.TrimSpace(matched[2])
	}
	return ""
}

func _MobileURL(URL string) string {
	if matched := urlRe.FindStringSubmatch(URL); len(matched) > 0 {
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

func _Emptysession() *GuestSession {
	return &GuestSession{
		conn: &Connection{},
	}
}

func _NewMobiledDocument(URL string) (*goquery.Document, error) {
	resp, err := get(_Emptysession(), URL)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}

func _Time(s string) time.Time {
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

func _CheckResponse(resp *http.Response) error {
	jsonResponse := &_JSONResponse{}
	if err := _ResponseUnmarshal(jsonResponse, resp); err != nil {
		return err
	}
	return _CheckResult(jsonResponse)
}

func _ResponseUnmarshal(data interface{}, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	body = []byte(strings.Trim(string(body), "[]"))
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

func _CheckResult(jsonResponse *_JSONResponse) error {
	if jsonResponse.Result == false {
		if jsonResponse.Cause != "" {
			return errors.New(jsonResponse.Cause)
		}
		return errUnknownCause
	}
	return nil
}

func _Form(m map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range m {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}

func _CalcArticleIcon(images []string) string {
	if len(images) > 0 {
		return "ico_p_y"
	}
	return "ico_t"
}

func _MultipartForm(m map[string]string, images ...string) (io.Reader, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if len(images) != 0 {
		_MultipartImages(w, images...)
		for i := range images {
			k := fmt.Sprintf("memo_block[%d]", i)
			m[k] = fmt.Sprintf("Dc_App_Img_%d", i+1)
		}
	}
	k := fmt.Sprintf("memo_block[%d]", len(images))
	m[k] = m["content"]
	delete(m, "content")
	_MultipartOthers(w, m)
	return &b, w.FormDataContentType()
}

func _MultipartImages(w *multipart.Writer, images ...string) {
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

func _MultipartOthers(w *multipart.Writer, m map[string]string) {
	for k, v := range m {
		if fw, err := w.CreateFormField(k); err != nil {
			continue
		} else if _, err := fw.Write([]byte(v)); err != nil {
			continue
		}
	}
}
