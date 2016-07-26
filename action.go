package goinside

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ThumbsUp 함수는 인자로 전달받은 글에 대해 추천을 보냅니다.
func (s *Session) ThumbsUp(a *Article) error {
	return s.action(a, recommendUpAPI)
}

// ThumbsDown 함수는 인자로 전달받은 글에 대해 비추천을 보냅니다.
func (s *Session) ThumbsDown(a *Article) error {
	return s.action(a, recommendDownAPI)
}

func (s *Session) action(a *Article, URL string) error {
	_, err := s.api(URL, form(map[string]string{
		"app_id": AppID,
		"id":     a.Gall.ID,
		"no":     a.Number,
	}), nonCharsetContentType)
	if err != nil {
		return err
	}
	return nil
}

// Report 함수는 인자로 전달받은 게시물을 신고합니다.
func (s *Session) Report(URL, memo string) error {
	must := func(s string, e error) string {
		if e != nil {
			panic(e)
		}
		return s
	}

	resp, err := s.api(reportAPI, form(map[string]string{
		"name":     must(url.QueryUnescape(s.id)),
		"password": must(url.QueryUnescape(s.pw)),
		"choice":   "4",
		"memo":     must(url.QueryUnescape(memo)),
		"no":       parseArticleNumber(URL),
		"id":       parseGallID(URL),
		"app_id":   AppID,
	}), nonCharsetContentType)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body = []byte(strings.Trim(string(body), "[]"))

	var respJSON struct {
		Result bool
		Cause  string
	}
	json.Unmarshal(body, &respJSON)
	if respJSON.Result == false {
		return errors.New(respJSON.Cause)
	}
	return nil
}
