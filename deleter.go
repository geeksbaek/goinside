package goinside

import "fmt"

type Deletable interface {
	delete(s Session) error
}

func (a *Article) delete(s Session) error {
	form, contentType := s.articleDeleteForm(a)
	resp, err := api(s, articleDeleteAPI, form, contentType)
	if err != nil {
		return err
	}
	return _CheckResponse(resp)
}

func (c *Comment) delete(s Session) error {
	form, contentType := s.commentDeleteForm(c)
	fmt.Println(form)
	resp, err := api(s, commentDeleteAPI, form, contentType)
	if err != nil {
		return err
	}
	return _CheckResponse(resp)
}
