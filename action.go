package goinside

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

var (
	koNameRe     = regexp.MustCompile(`query \+= "&ko_name=(.*)"`)
	gServerRe    = regexp.MustCompile(`query \+= "&gserver=(.*)"`)
	gNoRe        = regexp.MustCompile(`query \+= "&gno=(.*)"`)
	ipRe         = regexp.MustCompile(`query \+= "&ip=(.*)"`)
	categoryNoRe = regexp.MustCompile(`query \+= "&category_no=(.*)"`)
)

// ThumbsUp 함수는 인자로 전달받은 글에 대해 추천을 보냅니다.
func (s *Session) ThumbsUp(a *Article) error {
	form, err := s.commonRecommendForm(a)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_recomPrev_%s", a.GallID, a.Number): "done",
	})
	_, err = s.post(RecommendURL, cookies, form, defaultContentType)
	return err
}

// ThumbsDown 함수는 인자로 전달받은 글에 대해 비추천을 보냅니다.
func (s *Session) ThumbsDown(a *Article) error {
	form, err := s.commonRecommendForm(a)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_nonrecomPrev_%s", a.GallID, a.Number): "done",
	})
	_, err = s.post(NorecommendURL, cookies, form, defaultContentType)
	return err
}

func (s *Session) commonRecommendForm(a *Article) (io.Reader, error) {
	resp, err := s.get(a.URL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ip := ipRe.FindSubmatch(body)
	koName := koNameRe.FindSubmatch(body)
	gServer := gServerRe.FindSubmatch(body)
	gNo := gNoRe.FindSubmatch(body)
	categoryNo := categoryNoRe.FindSubmatch(body)
	if len(ip) != 2 || len(koName) != 2 || len(gServer) != 2 || len(gNo) != 2 || len(categoryNo) != 2 {
		return nil, errors.New("Make Recommend Form Fail")
	}
	return form(map[string]string{
		"no":          a.Number,
		"gall_id":     a.GallID,
		"ip":          string(ip[1]),
		"ko_name":     string(koName[1]),
		"gserver":     string(gServer[1]),
		"gno":         string(gNo[1]),
		"category_no": string(categoryNo[1]),
	}), nil
}
