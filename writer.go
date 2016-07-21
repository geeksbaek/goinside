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
	"os"
	"regexp"
	"strings"
)

var (
	flDataRe  = regexp.MustCompile(`\('FL_DATA'\).value ?= ?'(.*)'`)
	oflDataRe = regexp.MustCompile(`\('OFL_DATA'\).value ?= ?'(.*)'`)
	urlRe     = regexp.MustCompile(`url="?(.*?)"?>`)
	idRe      = regexp.MustCompile(`id=([^&]*)`)
	numberRe  = regexp.MustCompile(`no=(\d+)`)
	scriptRe  = regexp.MustCompile(`(?s)<script.*>(.+?)<\/script>`)
)

// WriteArticle 함수는 글을 작성합니다.
func (s *Session) WriteArticle(gallID, subject, content string, images ...string) (*Article, error) {
	return (&articleWriter{
		Session: s,
		gall:    &GallInfo{ID: gallID},
		subject: subject,
		content: trimContent(scriptRe.ReplaceAllString(content, "")),
		images:  images,
	}).write(false)
}

func (a *articleWriter) write(isCaptcha bool) (*Article, error) {
	m := map[string]string{
		"app_id":   AppID,
		"mode":     "write",
		"name":     a.id,
		"password": a.pw,
		"id":       a.gall.ID,
		"subject":  a.subject,
		"content":  a.content,
	}

	if isCaptcha {
		code := generateMD5()
		URL := fmt.Sprintf("http://m.dcinside.com/code.php?id=%s&dccode=%s", a.gall.ID, code)
		parsedCaptcha, err := ocr(a, URL)
		if err != nil {
			return nil, err
		}
		m["code"] = code
		m["dcblock"] = parsedCaptcha
	}

	f, contentType := multipartForm(m, a.images...)
	resp, err := a.api(gallArticleWriteAPI, f, contentType)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respJSON struct {
		Result bool
		Cause  string
		ID     string
	}
	body = []byte(strings.Trim(string(body), "[]"))
	json.Unmarshal(body, &respJSON)
	if respJSON.Result == false {
		if respJSON.Cause != "" {
			if regexp.MustCompile(`코드`).MatchString(respJSON.Cause) {
				return a.write(true)
			}
			return nil, errors.New(respJSON.Cause)
		}
		return nil, errors.New("writeAPI: json.Unmarshal fail")
	}
	return &Article{
		Gall: &GallInfo{
			ID: respJSON.ID, // same a.gall.ID
		},
		URL:     fmt.Sprintf("http://m.dcinside.com/view.php?id=%s&no=%s", respJSON.ID, respJSON.Cause),
		Number:  respJSON.Cause,
		Subject: a.subject,
		Content: a.content,
	}, nil
}

func (a *Article) delete(s *Session) error {
	// get cookies and con key
	m := map[string]string{}
	if s.isGuest {
		m["token_verify"] = "nonuser_del"
	} else {
		return errors.New("Need to login")
	}
	cookies, authKey, err := s.getCookiesAndAuthKey(m, accessTokenURL)
	if err != nil {
		return err
	}

	// delete article
	form := form(map[string]string{
		"id":       a.Gall.ID,
		"write_pw": s.pw,
		"no":       a.Number,
		"mode":     "board_del2",
		"con_key":  authKey,
	})
	_, err = s.post(optionWriteURL, cookies, form, defaultContentType)
	return err
}

func (s *Session) getCookiesAndAuthKey(m map[string]string, URL string) ([]*http.Cookie, string, error) {
	var cookies []*http.Cookie
	var authKey string
	form := form(m)
	resp, err := s.post(URL, nil, form, defaultContentType)
	if err != nil {
		return nil, "", err
	}
	cookies = resp.Cookies()
	authKey, err = parseAuthKey(resp)
	if err != nil {
		return nil, "", err
	}
	return cookies, authKey, nil
}

func parseAuthKey(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tempJSON struct {
		Msg  string
		Data string
	}
	json.Unmarshal(body, &tempJSON)
	if tempJSON.Data == "" {
		return "", errors.New("Block Key Parse Fail")
	}
	return tempJSON.Data, nil
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

// WriteComment 함수는 주어진 Article로 댓글을 작성합니다.
func (s *Session) WriteComment(a *Article, content string) (*Comment, error) {
	return (&commentWriter{
		Session: s,
		target:  a,
		content: content,
	}).write()
}

func (c *commentWriter) write() (*Comment, error) {
	form := form(map[string]string{
		"id":           c.target.Gall.ID,
		"no":           c.target.Number,
		"comment_nick": c.id,
		"comment_pw":   c.pw,
		"comment_memo": c.content,
		"mode":         "comment_nonmember",
	})
	resp, err := c.post(commentURL, nil, form, defaultContentType)
	if err != nil {
		return nil, err
	}
	commentNumber, err := parseCommentNumber(resp)
	if err != nil {
		return nil, err
	}
	URL := fmt.Sprintf("http://m.dcinside.com/view.php?id=%s&no=%s",
		c.target.Gall.ID, c.target.Number)
	return &Comment{
		Gall:    &GallInfo{URL: URL, ID: c.target.Gall.ID},
		Parents: &Article{Number: c.target.Number},
		Number:  commentNumber,
	}, nil
}

func (c *Comment) delete(s *Session) error {
	// get cookies and con key
	m := map[string]string{}
	if s.isGuest {
		m["token_verify"] = "nonuser_com_del"
	} else {
		return errors.New("Need to login")
	}
	cookies, authKey, err := s.getCookiesAndAuthKey(m, accessTokenURL)
	if err != nil {
		return err
	}

	// delete comment
	form := form(map[string]string{
		"id":         c.Gall.ID,
		"no":         c.Parents.Number,
		"iNo":        c.Number,
		"comment_pw": s.pw,
		"user_no":    "nonmember",
		"mode":       "comment_del",
		"con_key":    authKey,
	})
	_, err = s.post(optionWriteURL, cookies, form, defaultContentType)
	return err
}

func parseCommentNumber(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tempJSON struct {
		Msg  string
		Data string
	}
	json.Unmarshal(body, &tempJSON)
	if tempJSON.Data == "" {
		return "", errors.New("Block Key Parse Fail")
	}
	return tempJSON.Data, nil
}

// Delete 함수는 인자로 주어진 글을 삭제합니다.
func (s *Session) Delete(ds deletable) error {
	// done := make(chan error)
	// defer close(done)
	// for _, d := range ds {
	// 	d := d
	// 	go func() {
	// 		done <- d.delete(s)
	// 	}()
	// }
	// for _ = range ds {
	// 	if err := <-done; err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	return ds.delete(s)
}
