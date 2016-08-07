package goinside

import (
	"errors"
	"io"
	"net/url"
)

var (
	errInvalidIDorPW = errors.New("invalid ID or PW")
)

// GuestSession 구조체는 유동닉의 세션을 위해 사용됩니다.
type GuestSession struct {
	id   string
	pw   string
	conn *Connection
}

// Guest 함수는 전달받은 ID, PASSWORD로 생성한 유동닉 세션을 반환합니다.
func Guest(id, pw string) (gs *GuestSession, err error) {
	if len(id) == 0 || len(pw) == 0 {
		err = errInvalidIDorPW
		return
	}
	gs = &GuestSession{id: id, pw: pw, conn: &Connection{}}
	return
}

func (gs *GuestSession) connection() *Connection {
	return gs.conn
}

// Write 메소드는 쓰기 가능한 객체를 전달받아 작성 요청을 보내고,
// 삭제 가능한 작성된 객체를 반환합니다.
func (gs *GuestSession) Write(wa Writable) error {
	return wa.write(gs)
}

func (gs *GuestSession) articleWriteForm(ad *ArticleDraft) (io.Reader, string) {
	return _MultipartForm(map[string]string{
		"app_id":   appID,
		"mode":     "write",
		"name":     gs.id,
		"password": gs.pw,
		"id":       ad.GallID,
		"subject":  ad.Subject,
		"content":  ad.Content,
	}, ad.Images...)
}

func (gs *GuestSession) commentWriteForm(cd *CommentDraft) (io.Reader, string) {
	return _Form(map[string]string{
		"app_id":       appID,
		"comment_nick": gs.id,
		"comment_pw":   gs.pw,
		"id":           cd.Target.Gall.ID,
		"no":           cd.Target.Number,
		"comment_memo": cd.Content,
		"mode":         "comment_nonmember",

		// "best_chk":"N",
		// "board_id":"whiteking",
		// "best_comno":"0",
	}), defaultContentType
}

// Delete 메소드는 삭제 가능한 객체를 전달받아 삭제 요청을 보냅니다.
func (gs *GuestSession) Delete(da Deletable) error {
	return da.delete(gs)
}

func (gs *GuestSession) articleDeleteForm(a *Article) (io.Reader, string) {
	return _Form(map[string]string{
		"app_id":   appID,
		"mode":     "board_del",
		"write_pw": gs.pw,
		"id":       a.Gall.ID,
		"no":       a.Number,
	}), defaultContentType
}

func (gs *GuestSession) commentDeleteForm(c *Comment) (io.Reader, string) {
	return _Form(map[string]string{
		"app_id":     appID,
		"comment_pw": gs.pw,
		"id":         c.Parents.Gall.ID,
		"no":         c.Parents.Number,
		"mode":       "comment_del",
		"comment_no": c.Number,
		// "board_id":   c.Parents.Author.Detail.GallogID,
		// "best_chk":   "N",
		// "best_comno": "0",
	}), nonCharsetContentType
}

// ThumbsUp 메소드는 해당 글에 추천 요청을 보냅니다.
func (gs *GuestSession) ThumbsUp(a *Article) error {
	return a.thumbsUp(gs)
}

// ThumbsDown 메소드는 해당 글에 비추천 요청을 보냅니다.
func (gs *GuestSession) ThumbsDown(a *Article) error {
	return a.thumbsDown(gs)
}

func (gs *GuestSession) actionForm(a *Article) (io.Reader, string) {
	return _Form(map[string]string{
		"app_id": appID,
		"id":     a.Gall.ID,
		"no":     a.Number,
	}), nonCharsetContentType
}

// Report 메소드는 해당 글에 메모와 함께 신고 요청을 보냅니다.
func (gs *GuestSession) Report(a *Article, memo string) error {
	return a.report(gs, memo)
}

func (gs *GuestSession) reportForm(URL, memo string) (io.Reader, string) {
	_Must := func(s string, e error) string {
		if e != nil {
			panic(e)
		}
		return s
	}
	return _Form(map[string]string{
		"name":     _Must(url.QueryUnescape(gs.id)),
		"password": _Must(url.QueryUnescape(gs.pw)),
		"choice":   "4",
		"memo":     _Must(url.QueryUnescape(memo)),
		"no":       _ParseArticleNumber(URL),
		"id":       _ParseGallID(URL),
		"app_id":   appID,
	}), nonCharsetContentType
}
