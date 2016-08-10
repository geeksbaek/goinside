package goinside

type deletable interface {
	delete(s session) error
}

func (a *Article) delete(s session) error {
	form, contentType := s.articleDeleteForm(a)
	resp, err := api(s, articleDeleteAPI, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (c *Comment) delete(s session) error {
	form, contentType := s.commentDeleteForm(c)
	resp, err := api(s, commentDeleteAPI, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}
