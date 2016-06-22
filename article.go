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
	numberRe  = regexp.MustCompile(`no="?(.*)"?`)
)

type Article struct {
	URL    string
	GallID string
	Number string
}

type ArticleWriter struct {
	Auth
	GallID  string
	Subject string
	Content string
	Images  []string
}

func (a *Auth) WriteArticle(atw *ArticleWriter) (*Article, error) {
	// get cookies and blockkey
	cookies, authKey, err := a.getCookiesAndAuthKey(map[string]string{
		"id":        "programming",
		"w_subject": atw.Subject,
		"w_memo":    atw.Content,
		"w_filter":  "1",
		"mode":      "write_verify",
	})
	if err != nil {
		return nil, err
	}

	// upload images and get FL_DATA, OFL_DATA string
	var flData, oflData string
	if len(atw.Images) > 0 {
		flData, oflData, err = a.UploadImages(atw.Images, atw.GallID)
		if err != nil {
			return nil, err
		}
	}

	// wrtie article
	ret := &Article{}
	form, contentType := multipartForm(nil, map[string]string{
		"name":       a.id,
		"password":   a.pw,
		"subject":    atw.Subject,
		"memo":       atw.Content,
		"mode":       "write",
		"id":         atw.GallID,
		"mobile_key": "mobile_nomember",
		"FL_DATA":    flData,
		"OFL_DATA":   oflData,
		"Block_key":  authKey,
		"filter":     "1",
	})
	resp, err := a.Post(gWrite, cookies, form, contentType)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret.URL = string(urlRe.Find(body))
	ret.Number = string(numberRe.Find(body))
	return ret, nil
}

func (a *Auth) DeleteArticle(at *Article) error {
	// get cookies and conkey
	m := map[string]string{}
	if a.nomember {
		m["token_verify"] = "nonuser_del"
	} else {
		return errors.New("Need to login")
	}
	cookies, authKey, err := a.getCookiesAndAuthKey(m)
	if err != nil {
		return err
	}

	// delete article
	form := form(map[string]string{
		"id":       at.GallID,
		"write_pw": a.pw,
		"no":       at.Number,
		"mode":     "board_del2",
		"con_key":  authKey,
	})
	_, err = a.Post(optionWrite, cookies, form, defaultContentType)
	return err
}

func (a *Auth) DeleteArticles(ats []*Article) error {
	done := make(chan error)
	defer close(done)
	for _, at := range ats {
		at := at
		go func() {
			done <- a.DeleteArticle(at)
		}()
	}
	for _ = range ats {
		if err := <-done; err != nil {
			return err
		}
	}
	return nil
}

func (a *Auth) UploadImages(images []string, gall string) (string, string, error) {
	form, contentType := multipartForm(images, map[string]string{
		"imgId":   gall,
		"mode":    "write",
		"img_num": "11",
	})
	resp, err := a.Post(uploadImage, nil, form, contentType)
	if err != nil {
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	return string(flDataRe.Find(body)), string(oflDataRe.Find(body)), nil
}

func (a *Auth) getCookiesAndAuthKey(m map[string]string) ([]*http.Cookie, string, error) {
	var cookies []*http.Cookie
	var authKey string
	form := form(m)
	resp, err := a.Post(optionWrite, nil, form, defaultContentType)
	if err != nil {
		return nil, "", err
	}
	cookies = resp.Cookies()
	authKey, err = parseAuthkKey(resp)
	if err != nil {
		return nil, "", err
	}
	return cookies, authKey, nil
}

func parseAuthkKey(resp *http.Response) (string, error) {
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
