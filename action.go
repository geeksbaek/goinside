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
	categoryNoRe = regexp.MustCompile(`query \+= "&category_no=(.*)"`)
	gIPRe        = regexp.MustCompile(`query \+= "&ip=(.*)"`)
)

// ThumbsUp 함수는 인자로 전달받은 글에 대해 추천을 보냅니다.
func (s *Session) ThumbsUp(a *Article) error {
	form, err := s.commonRecommendForm(a)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_recomPrev_%s", a.Gall.ID, a.Number): "done",
	})
	_, err = s.post(recommendURL, cookies, form, defaultContentType)
	return err
}

// ThumbsDown 함수는 인자로 전달받은 글에 대해 비추천을 보냅니다.
func (s *Session) ThumbsDown(a *Article) error {
	form, err := s.commonRecommendForm(a)
	if err != nil {
		return err
	}
	cookies := cookies(map[string]string{
		fmt.Sprintf("%s_nonrecomPrev_%s", a.Gall.ID, a.Number): "done",
	})
	_, err = s.post(norecommendURL, cookies, form, defaultContentType)
	return err
}

func (s *Session) commonRecommendForm(a *Article) (io.Reader, error) {
	detail := a.Gall.detail
	if detail == nil {
		var err error
		if detail, err = s.PrefetchDetail(a); err != nil {
			return nil, err
		}
	}
	ip := detail.ip
	if s.ip != "" {
		ip = s.ip
	}
	return form(map[string]string{
		"no":          a.Number,
		"gall_id":     a.Gall.ID,
		"ip":          ip,
		"ko_name":     detail.koName,
		"gserver":     detail.gServer,
		"gno":         detail.gNo,
		"category_no": detail.categoryNo,
	}), nil
}

// PrefetchDetail 함수는 추천, 비추천에 필요한 세부 값을 미리 가져옵니다.
func (s *Session) PrefetchDetail(a *Article) (*gallInfoDetail, error) {
	resp, err := s.get(a.URL)
	if err != nil {
		return nil, err
	}
	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bytesBody)
	koName := koNameRe.FindStringSubmatch(body)
	gServer := gServerRe.FindStringSubmatch(body)
	gNo := gNoRe.FindStringSubmatch(body)
	categoryNo := categoryNoRe.FindStringSubmatch(body)
	gIP := gIPRe.FindStringSubmatch(body)
	if len(koName) != 2 || len(gServer) != 2 || len(gNo) != 2 || len(categoryNo) != 2 || len(gIP) != 2 {
		return nil, errors.New("Make Recommend Form Fail")
	}
	return &gallInfoDetail{
		koName[1], gServer[1], gNo[1], categoryNo[1], gIP[1],
	}, nil
}
