package goinside

import (
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

func (a *Auth) Recommend(at *Article) error {
	form, err := a.commonRecommendForm(at)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_recomPrev_%s", at.GallID, at.Number): "done",
	})
	_, err = a.Post(recommend, cookies, form, defaultContentType)
	return err
}

func (a *Auth) Norecommend(at *Article) error {
	form, err := a.commonRecommendForm(at)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_nonrecomPrev_%s", at.GallID, at.Number): "done",
	})
	_, err = a.Post(recommend, cookies, form, defaultContentType)
	return err
}

func (a *Auth) commonRecommendForm(at *Article) (io.Reader, error) {
	resp, err := a.Get(at.URL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	koName := string(koNameRe.Find(body))
	gServer := string(gServerRe.Find(body))
	gNo := string(gNoRe.Find(body))
	categoryNo := string(categoryNoRe.Find(body))
	return form(map[string]string{
		"no":          at.Number,
		"gall_id":     at.GallID,
		"ip":          a.ip,
		"ko_name":     koName,
		"gserver":     gServer,
		"gno":         gNo,
		"category_no": categoryNo,
	}), nil
}
