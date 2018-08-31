package goinside

import (
	"errors"
	"io"
	"net/url"
)

var (
	errLoginFailed = errors.New("login failed")
)

// MemberSession 구조체는 고정닉의 세션을 나타냅니다.
type MemberSession struct {
	id   string
	pw   string
	conn *Connection
	*MemberSessionDetail
}

// MemberSessionDetail 구조체는 해당 세션의 세부 정보를 나타냅니다.
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

// Connection 메소드는 해당 세션의 Connection 구조체를 반환합니다.
func (ms *MemberSession) Connection() *Connection {
	if ms.conn == nil {
		ms.conn = &Connection{}
	}
	return ms.conn
}

// Write 메소드는 글이나 댓글과 같은 쓰기 가능한 객체를 전달받아 작성 요청을 보냅니다.
func (ms *MemberSession) Write(w writable) error {
	return w.write(ms)
}

func (ms *MemberSession) articleWriteForm(id, s, c string, is ...string) (io.Reader, string) {
	return multipartForm(map[string]string{
		"app_id":  GetAppID(ms),
		"mode":    "write",
		"user_id": ms.UserID,
		"id":      id,
		"subject": s,
		"content": c,
	}, is...)
}

func (ms *MemberSession) commentWriteForm(id, n, c string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":       GetAppID(ms),
		"user_id":      ms.UserID,
		"id":           id,
		"no":           n,
		"comment_memo": c,
		"mode":         "comment",
	}), defaultContentType
}

// Delete 메소드는 삭제 가능한 객체를 전달받아 삭제 요청을 보냅니다.
func (ms *MemberSession) Delete(d deletable) error {
	return d.delete(ms)
}

func (ms *MemberSession) articleDeleteForm(id, n string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":  GetAppID(ms),
		"user_id": ms.UserID,
		"no":      n,
		"id":      id,
		"mode":    "board_del",
	}), defaultContentType
}

func (ms *MemberSession) commentDeleteForm(id, n, cn string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":     GetAppID(ms),
		"user_id":    ms.UserID,
		"id":         id,
		"no":         n,
		"mode":       "comment_del",
		"comment_no": cn,
	}), defaultContentType
}

// ThumbsUp 메소드는 해당 글에 추천 요청을 보냅니다.
func (ms *MemberSession) ThumbsUp(a actionable) error {
	return a.thumbsUp(ms)
}

// ThumbsDown 메소드는 해당 글에 비추천 요청을 보냅니다.
func (ms *MemberSession) ThumbsDown(a actionable) error {
	return a.thumbsDown(ms)
}

func (ms *MemberSession) actionForm(id, n string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id": GetAppID(ms),
		"id":     id,
		"no":     n,
	}), nonCharsetContentType
}

// Report 메소드는 해당 글에 메모와 함께 신고 요청을 보냅니다.
func (ms *MemberSession) Report(a actionable, memo string) error {
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
		"app_id":     GetAppID(ms),
	}), nonCharsetContentType
}
