package goinside

import "errors"

var errReportResultFalseWithEmptyCause = errors.New("Report Result False With Empty Cause")

type actionable interface {
	thumbsUp(session) error
	thumbsDown(session) error
	report(session, string) error
}

func (a *Article) thumbsUp(s session) error {
	return a.action(s, recommendUpAPI)
}

func (a *Article) thumbsDown(s session) error {
	return a.action(s, recommendDownAPI)
}

func (a *Article) action(s session, api dcinsideAPI) error {
	form, contentType := s.actionForm(a.ID, a.Number)
	resp, err := api.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (a *Article) report(s session, memo string) error {
	form, contentType := s.reportForm(a.URL, memo)
	resp, err := reportAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (i *ListItem) thumbsUp(s session) error {
	return i.action(s, recommendUpAPI)
}

func (i *ListItem) thumbsDown(s session) error {
	return i.action(s, recommendDownAPI)
}

func (i *ListItem) action(s session, api dcinsideAPI) error {
	form, contentType := s.actionForm(i.ID, i.Number)
	resp, err := api.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (i *ListItem) report(s session, memo string) error {
	form, contentType := s.reportForm(i.URL, memo)
	resp, err := reportAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}
