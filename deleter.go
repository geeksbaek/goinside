package goinside

type deletable interface {
	delete(s session) error
}

func (a *Article) delete(s session) error {
	form, contentType := s.articleDeleteForm(a.ID, a.Number)
	resp, err := deleteArticleAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (c *Comment) delete(s session) error {
	form, contentType := s.commentDeleteForm(c.ID, c.Parents.Number, c.Number)
	resp, err := deleteCommentAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (i *ListItem) delete(s session) error {
	form, contentType := s.articleDeleteForm(i.ID, i.Number)
	resp, err := deleteArticleAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}
