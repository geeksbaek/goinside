package goinside

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

var (
	errInvalidIDorPW = errors.New("invalid ID or PW")
)

var dummyGuest = RandomGuest()

// GuestSession 구조체는 유동닉 세션을 나타냅니다.
type GuestSession struct {
	id   string
	pw   string
	conn *Connection
	app  *App
}

// Guest 함수는 유동닉 세션을 반환합니다.
func Guest(id, pw string) (gs *GuestSession, err error) {
	if len(id) == 0 || len(pw) == 0 {
		err = errInvalidIDorPW
		return
	}
	gs = &GuestSession{
		id:   id,
		pw:   pw,
		conn: &Connection{timeout: time.Second * 5},
	}
	valueToken, appID, err := fetchAppID(gs)
	if err != nil {
		return nil, err
	}
	gs.app = &App{Token: valueToken, ID: appID}
	return
}

// RandomGuest 함수는 임의의 유동닉 세션을 반환합니다.
func RandomGuest() *GuestSession {
	gs := &GuestSession{
		id:   fmt.Sprint(rand.Intn(100)),
		pw:   fmt.Sprint(rand.Intn(100)),
		conn: &Connection{timeout: time.Second * 5},
	}
	valueToken, appID, _ := fetchAppID(gs)
	gs.app = &App{Token: valueToken, ID: appID}
	return gs
}

// Connection 메소드는 해당 세션의 Connection 구조체를 반환합니다.
func (gs *GuestSession) Connection() *Connection {
	if gs.conn == nil {
		gs.conn = &Connection{}
	}
	return gs.conn
}

// Write 메소드는 쓰기 가능한 객체를 전달받아 작성 요청을 보냅니다.
func (gs *GuestSession) Write(w writable) error {
	return w.write(gs)
}

func (gs *GuestSession) articleWriteForm(id, s, c string, is ...string) (io.Reader, string) {
	return multipartForm(map[string]string{
		"app_id":   gs.getAppID(),
		"mode":     "write",
		"name":     gs.id,
		"password": gs.pw,
		"id":       id,
		"subject":  s,
		"content":  c,
	}, is...)
}

func (gs *GuestSession) commentWriteForm(id, n, c string, is ...string) (io.Reader, string) {
	return multipartForm(map[string]string{
		"app_id":       gs.getAppID(),
		"comment_nick": gs.id,
		"comment_pw":   gs.pw,
		"id":           id,
		"no":           n,
		"comment_memo": c,
		"mode":         "com_write",
	}, is...)
}

// Delete 메소드는 삭제 가능한 객체를 전달받아 삭제 요청을 보냅니다.
func (gs *GuestSession) Delete(d deletable) error {
	return d.delete(gs)
}

func (gs *GuestSession) articleDeleteForm(id, n string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":   gs.getAppID(),
		"mode":     "board_del",
		"write_pw": gs.pw,
		"id":       id,
		"no":       n,
	}), defaultContentType
}

func (gs *GuestSession) commentDeleteForm(id, n, cn string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":     gs.getAppID(),
		"comment_pw": gs.pw,
		"id":         id,
		"no":         n,
		"mode":       "comment_del",
		"comment_no": cn,
	}), defaultContentType
}

// ThumbsUp 메소드는 해당 글에 추천 요청을 보냅니다.
func (gs *GuestSession) ThumbsUp(a actionable) error {
	return a.thumbsUp(gs)
}

// ThumbsDown 메소드는 해당 글에 비추천 요청을 보냅니다.
func (gs *GuestSession) ThumbsDown(a actionable) error {
	return a.thumbsDown(gs)
}

func (gs *GuestSession) actionForm(id, n string) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id": gs.getAppID(),
		"id":     id,
		"no":     n,
	}), nonCharsetContentType
}

func (gs *GuestSession) getAppID() string {
	valueToken, err := generateValueToken()
	if err != nil {
		return ""
	}
	if gs.app.Token == valueToken {
		return gs.app.ID
	}
	valueToken, appID, err := fetchAppID(gs)
	if err != nil {
		return ""
	}
	gs.app = &App{Token: valueToken, ID: appID}
	return gs.app.ID
}
