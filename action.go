package goinside

import "errors"

var errReportResultFalseWithEmptyCause = errors.New("Report Result False With Empty Cause")

func (a *Article) thumbsUp(s Session) error {
	return a.action(s, recommendUpAPI)
}

func (a *Article) thumbsDown(s Session) error {
	return a.action(s, recommendDownAPI)
}

func (a *Article) action(s Session, URL string) error {
	form, contentType := s.actionForm(a)
	resp, err := api(s, URL, form, contentType)
	if err != nil {
		return err
	}
	return _CheckResponse(resp)
}

func (a *Article) report(s Session, memo string) error {
	form, contentType := s.reportForm(a.URL, memo)
	resp, err := api(s, reportAPI, form, contentType)
	if err != nil {
		return err
	}
	return _CheckResponse(resp)
}
