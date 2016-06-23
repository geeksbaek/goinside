package goinside

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Comment 구조체는 댓글 정보를 표현합니다.
type Comment struct {
	URL     string
	GallID  string
	Number  string
	Cnumber string
}

// CommentWriter 구조체는 댓글 작성에 필요한 정보를 전달하기 위한 구조체입니다.
type CommentWriter struct {
	*Session
	*Article
	Content string
}

// NewComment 함수는 새로운 CommentWriter 객체를 반환합니다.
func (s *Session) NewComment(a *Article, content string) *CommentWriter {
	return &CommentWriter{
		Session: s,
		Article: a,
		Content: content,
	}
}

// WriteComment 함수는 리시버 Auth의 정보와 인자로 전달받은 CommentWriter 구조체의 정보를 조합하여 댓글을 작성합니다.
func (c *CommentWriter) Write() (*Comment, error) {
	form := form(map[string]string{
		"id":           c.GallID,
		"no":           c.Number,
		"ip":           c.ip,
		"comment_nick": c.id,
		"comment_pw":   c.pw,
		"comment_memo": c.Content,
		"mode":         "comment_nonmember",
	})
	resp, err := c.post(comment, nil, form, defaultContentType)
	if err != nil {
		return nil, err
	}
	commentNumber, err := parseCommentNumber(resp)
	if err != nil {
		return nil, err
	}
	return &Comment{
		URL:     fmt.Sprintf("http://m.dcinside.com/view.php?id=%s&no=%s", c.GallID, c.Number),
		GallID:  c.GallID,
		Number:  c.Number,
		Cnumber: commentNumber,
	}, nil
}

// DeleteComment 함수는 리시버 Auth의 정보와 인자로 전달받은 CommentWriter 구조체의 정보를 조합하여 댓글을 삭제합니다.
func (s *Session) DeleteComment(c *Comment) error {
	// get cookies and con key
	m := map[string]string{}
	if s.nomember {
		m["token_verify"] = "nonuser_com_del"
	} else {
		return errors.New("Need to login")
	}
	cookies, authKey, err := s.getCookiesAndAuthKey(m, accessToken)
	if err != nil {
		return err
	}

	// delete comment
	form := form(map[string]string{
		"id":         c.GallID,
		"no":         c.Number,
		"iNo":        c.Cnumber,
		"comment_pw": s.pw,
		"user_no":    "nonmember",
		"mode":       "comment_del",
		"con_key":    authKey,
	})
	_, err = s.post(optionWrite, cookies, form, defaultContentType)
	return err
}

// DeleteComments 함수는 인자로 전달받은 여러 개의 댓글에 대해 동시적으로 DeleteComment 함수를 호출합니다.
func (s *Session) DeleteComments(cs []*Comment) error {
	done := make(chan error)
	defer close(done)
	for _, c := range cs {
		c := c
		go func() {
			done <- s.DeleteComment(c)
		}()
	}
	for _ = range cs {
		if err := <-done; err != nil {
			return err
		}
	}
	return nil
}

func parseCommentNumber(resp *http.Response) (string, error) {
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
