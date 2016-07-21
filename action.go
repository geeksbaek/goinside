package goinside

// ThumbsUp 함수는 인자로 전달받은 글에 대해 추천을 보냅니다.
func (s *Session) ThumbsUp(a *Article) error {
	return s.action(a, recommendUpAPI)
}

// ThumbsDown 함수는 인자로 전달받은 글에 대해 비추천을 보냅니다.
func (s *Session) ThumbsDown(a *Article) error {
	return s.action(a, recommendDownAPI)
}

func (s *Session) action(a *Article, URL string) error {
	_, err := s.api(recommendUpAPI, form(map[string]string{
		"app_id": AppID,
		"id":     a.Gall.ID,
		"no":     a.Number,
	}), nonCharsetContentType)
	if err != nil {
		return err
	}
	return nil
}
