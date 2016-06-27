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
	"os"
	"regexp"
)

var (
	flDataRe  = regexp.MustCompile(`\('FL_DATA'\).value ?= ?'(.*)'`)
	oflDataRe = regexp.MustCompile(`\('OFL_DATA'\).value ?= ?'(.*)'`)
	urlRe     = regexp.MustCompile(`url="?(.*?)"?>`)
	idRe      = regexp.MustCompile(`id=([^&]*)`)
	numberRe  = regexp.MustCompile(`no=(\d+)`)
)

// WriteArticle 함수는 글을 작성합니다.
func (s *Session) WriteArticle(gallID, subject, content string, images ...string) (*Article, error) {
	return (&articleWriter{
		Session: s,
		gall:    &GallInfo{ID: gallID},
		subject: subject,
		content: content,
		images:  images,
	}).write()
}

func (a *articleWriter) write() (*Article, error) {
	// get cookies and block key
	cookies, authKey, err := a.getCookiesAndAuthKey(map[string]string{
		"id":        "programming",
		"w_subject": a.subject,
		"w_memo":    a.content,
		"w_filter":  "1",
		"mode":      "write_verify",
	}, optionWriteURL)
	if err != nil {
		return nil, err
	}

	// upload images and get FL_DATA, OFL_DATA string
	var flData, oflData string
	if len(a.images) > 0 {
		flData, oflData, err = a.uploadImages(a.gall.ID, a.images)
		if err != nil {
			return nil, err
		}
	}

	// wrtie article
	ret := &Article{Gall: &GallInfo{}}
	form, contentType := multipartForm(nil, map[string]string{
		"name":       a.id,
		"password":   a.pw,
		"subject":    a.subject,
		"memo":       a.content,
		"mode":       "write",
		"id":         a.gall.ID,
		"mobile_key": "mobile_nomember",
		"FL_DATA":    flData,
		"OFL_DATA":   oflData,
		"Block_key":  authKey,
		"filter":     "1",
	})
	resp, err := a.post(gWriteURL, cookies, form, contentType)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)
	URL := urlRe.FindStringSubmatch(body)
	gallID := idRe.FindStringSubmatch(body)
	number := numberRe.FindStringSubmatch(body)
	if len(URL) != 2 || len(gallID) != 2 || len(number) != 2 {
		return nil, errors.New("Write Article Fail")
	}
	ret.URL, ret.Gall.ID, ret.Number = URL[1], gallID[1], number[1]
	return ret, nil
}

func (a *Article) delete(s *Session) error {
	// get cookies and con key
	m := map[string]string{}
	if s.nomember {
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

func (as Articles) deleteAll(s *Session) error {
	done := make(chan error)
	defer close(done)
	for _, a := range as {
		a := a
		go func() {
			done <- a.delete(s)
		}()
	}
	for _ = range as {
		if err := <-done; err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) uploadImages(gall string, images []string) (string, string, error) {
	form, contentType := multipartForm(images, map[string]string{
		"imgId":   gall,
		"mode":    "write",
		"img_num": "11", // ?
	})
	resp, err := s.post(uploadImageURL, nil, form, contentType)
	if err != nil {
		return "", "", err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	body := string(bodyBytes)
	fldata := flDataRe.FindStringSubmatch(body)
	ofldata := oflDataRe.FindStringSubmatch(body)
	if len(fldata) != 2 || len(ofldata) != 2 {
		return "", "", errors.New("Image Upload Fail")
	}
	return fldata[1], ofldata[1], nil
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

func multipartForm(images []string, m map[string]string) (io.Reader, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if images != nil {
		multipartImages(w, images)
	}
	multipartOthers(w, m)
	return &b, w.FormDataContentType()
}

func multipartImages(w *multipart.Writer, images []string) {
	for i, image := range images {
		f, err := os.Open(image)
		if err != nil {
			return
		}
		defer f.Close()
		fw, err := w.CreateFormFile(fmt.Sprintf("upload[%d]", i), image)
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
	if s.nomember {
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

func (cs Comments) deleteAll(s *Session) error {
	done := make(chan error)
	defer close(done)
	for _, c := range cs {
		c := c
		go func() {
			done <- c.delete(s)
		}()
	}
	for _ = range cs {
		if err := <-done; err != nil {
			return err
		}
	}
	return nil
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
func (s *Session) Delete(d deletable) error {
	return d.delete(s)
}

// DeleteAll 함수는 인자로 주어진 여러 개의 글을 동시에 삭제합니다.
func (s *Session) DeleteAll(d delateAllable) error {
	return d.deleteAll(s)
}
