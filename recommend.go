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

// Recommend 함수는 인자로 전달받은 글에 대해 추천을 보냅니다.
func (a *Auth) Recommend(at *Article) error {
	form, err := a.commonRecommendForm(at)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_recomPrev_%s", at.GallID, at.Number): "done",
	})
	_, err = a.post(recommend, cookies, form, defaultContentType)
	return err
}

// Norecommend 함수는 인자로 전달받은 글에 대해 비추천을 보냅니다.
func (a *Auth) Norecommend(at *Article) error {
	form, err := a.commonRecommendForm(at)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_nonrecomPrev_%s", at.GallID, at.Number): "done",
	})
	_, err = a.post(recommend, cookies, form, defaultContentType)
	return err
}

func (a *Auth) commonRecommendForm(at *Article) (io.Reader, error) {
	resp, err := a.get(at.URL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	koName := koNameRe.FindSubmatch(body)
	gServer := gServerRe.FindSubmatch(body)
	gNo := gNoRe.FindSubmatch(body)
	categoryNo := categoryNoRe.FindSubmatch(body)
	if len(koName) != 2 || len(gServer) != 2 || len(gNo) != 2 || len(categoryNo) != 2 {
		return nil, errors.New("Make Recommend Form Fail")
	}
	return form(map[string]string{
		"no":          at.Number,
		"gall_id":     at.GallID,
		"ip":          a.ip,
		"ko_name":     string(koName[1]),
		"gserver":     string(gServer[1]),
		"gno":         string(gNo[1]),
		"category_no": string(categoryNo[1]),
	}), nil
}
