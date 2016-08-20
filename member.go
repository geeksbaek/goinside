package goinside

import (
	"errors"
	"io"
	"net/url"
)

var (
	errLoginFailed = errors.New("login failed")
)

// MemberSession 구조체는 고정닉의 세션을 표현합니다.
type MemberSession struct {
	id   string
	pw   string
	conn *Connection
	*MemberSessionDetail
}

type MemberSessionDetail struct {
	UserID string `json:"user_id"`
	UserNO string `json:"user_no"`
	Name   string `json:"name"`
	Stype  string `json:"stype"`
	// IsAdult    bool   `json:"is_adult"`
	// IsDormancy bool   `json:"is_dormancy"`
}

// Login 함수는 고정닉 세션을 반환합니다.
func Login(id, pw string) (ms *MemberSession, err error) {
	form := makeForm(map[string]string{
		"user_id": id,
		"user_pw": pw,
	})
	tempMS := &MemberSession{
		id:   id,
		pw:   pw,
		conn: &Connection{},
	}
	resp, err := loginAPI.post(tempMS, form, defaultContentType)
	if err != nil {
		return
	}
	tempMSD := new([]MemberSessionDetail)
	err = responseUnmarshal(resp, tempMSD)
	if err != nil {
		return
	}
	if !(*tempMSD)[0].isSucceed() {
		err = errLoginFailed
		return
	}
	tempMS.MemberSessionDetail = &((*tempMSD)[0])
	ms = tempMS
	return
}

func (msd *MemberSessionDetail) isSucceed() bool {
	switch {
	case msd.Name == "":
		return false
	case msd.UserID == "":
		return false
	case msd.UserNO == "":
		return false
	}
	return true
}

// Logout 메소드는 해당 고정닉 세션을 종료합니다.
func (ms *MemberSession) Logout() (err error) {
	ms = nil
	return
}

func (ms *MemberSession) Connection() *Connection {
	if ms.conn == nil {
		ms.conn = &Connection{}
	}
	return ms.conn
}

// Write 메소드는 글이나 댓글과 같은 쓰기 가능한 객체를 전달받아 작성 요청을 보냅니다.
func (ms *MemberSession) Write(wa writable) error {
	return wa.write(ms)
}

func (ms *MemberSession) articleWriteForm(ad *ArticleDraft) (io.Reader, string) {
	return multipartForm(map[string]string{
		"app_id":  AppID,
		"mode":    "write",
		"user_id": ms.UserID,
		"id":      ad.GallID,
		"subject": ad.Subject,
		"content": ad.Content,
	}, ad.Images...)
}

func (ms *MemberSession) commentWriteForm(cd *CommentDraft) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":       AppID,
		"user_id":      ms.UserID,
		"id":           cd.Target.Gall.ID,
		"no":           cd.Target.Number,
		"comment_memo": cd.Content,
		"mode":         "comment",
	}), defaultContentType
}

// Delete 메소드는 삭제 가능한 객체를 전달받아 삭제 요청을 보냅니다.
func (ms *MemberSession) Delete(da deletable) error {
	return da.delete(ms)
}

func (ms *MemberSession) articleDeleteForm(a *Article) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":  AppID,
		"user_id": ms.UserID,
		"no":      a.Number,
		"id":      a.Gall.ID,
		"mode":    "board_del",
	}), defaultContentType
}

func (ms *MemberSession) commentDeleteForm(c *Comment) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":     AppID,
		"user_id":    ms.UserID,
		"id":         c.Parents.Gall.ID,
		"no":         c.Parents.Number,
		"mode":       "comment_del",
		"comment_no": c.Number,
	}), defaultContentType
}

// ThumbsUp 메소드는 해당 글에 추천 요청을 보냅니다.
func (ms *MemberSession) ThumbsUp(a *Article) error {
	return a.thumbsUp(ms)
}

// ThumbsDown 메소드는 해당 글에 비추천 요청을 보냅니다.
func (ms *MemberSession) ThumbsDown(a *Article) error {
	return a.thumbsDown(ms)
}

func (ms *MemberSession) actionForm(a *Article) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id": AppID,
		"id":     a.Gall.ID,
		"no":     a.Number,
	}), nonCharsetContentType
}

// Report 메소드는 해당 글에 메모와 함께 신고 요청을 보냅니다.
func (ms *MemberSession) Report(a *Article, memo string) error {
	return a.report(ms, memo)
}

func (ms *MemberSession) reportForm(URL, memo string) (io.Reader, string) {
	_Must := func(s string, e error) string {
		if e != nil {
			return ""
		}
		return s
	}
	return makeForm(map[string]string{
		"confirm_id": ms.UserID,
		"choice":     "4",
		"memo":       _Must(url.QueryUnescape(memo)),
		"no":         articleNumber(URL),
		"id":         gallID(URL),
		"app_id":     AppID,
	}), nonCharsetContentType
}
