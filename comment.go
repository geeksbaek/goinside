package goinside

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Comment struct {
	URL     string
	GallID  string
	Number  string
	Cnumber string
}

type CommentWriter struct {
	Auth
	GallID  string
	Number  string
	Content string
}

func (a *Auth) WriteComment(cw *CommentWriter) (*Comment, error) {
	form := form(map[string]string{
		"id":           cw.GallID,
		"no":           cw.Number,
		"ip":           a.ip,
		"comment_nick": a.id,
		"comment_pw":   a.pw,
		"comment_memo": cw.Content,
		"mode":         "comment_nonmember",
	})
	resp, err := a.Post(comment, nil, form, defaultContentType)
	if err != nil {
		return nil, err
	}
	commentNumber, err := parseCommentNumber(resp)
	if err != nil {
		return nil, err
	}
	return &Comment{
		URL:    fmt.Sprintf("http://m.dcinside.com/view.php?id=%s&no=%s", cw.GallID, cw.Number),
		GallID: cw.GallID,
		Number: cw.Number,
		Cnumber: commentNumber,
	}, nil
}

func (a *Auth) DeleteComment(ct *Comment) error {
	// get cookies and conkey
	m := map[string]string{}
	if a.nomember {
		m["token_verify"] = "nonuser_del"
	} else {
		return errors.New("Need to login")
	}
	cookies, authKey, err := a.getCookiesAndAuthKey(m)
	if err != nil {
		return err
	}

	// delete comment
	form := form(map[string]string{
		"id":         ct.GallID,
		"no":         ct.Number,
		"iNo":        ct.Cnumber,
		"comment_pw": a.pw,
		"user_no":    "nonmember",
		"mode":       "comment_del",
		"con_key":    authKey,
	})
	_, err = a.Post(optionWrite, cookies, form, defaultContentType)
	return err
}

func (a *Auth) DeleteComments(cts []*Comment) error {
	done := make(chan error)
	defer close(done)
	for _, ct := range cts {
		ct := ct
		go func() {
			done <- a.DeleteComment(ct)
		}()
	}
	for _ = range cts {
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
