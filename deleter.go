package goinside

type deletable interface {
	delete(s session) error
}

func (a *Article) delete(s session) error {
	form, contentType := s.articleDeleteForm(a)
	resp, err := deleteArticleAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (c *Comment) delete(s session) error {
	form, contentType := s.commentDeleteForm(c)
	resp, err := deleteCommentAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}
