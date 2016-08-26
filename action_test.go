package goinside

// import "testing"

// func TestAction(t *testing.T) {
// 	s, err := Guest("ㅇㅇ", "123")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	URL := "http://gall.dcinside.com/board/lists/?id=baseball_new4"
// 	page := 1

// 	l, err := FetchList(URL, page)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(l.Items) != 25 {
// 		t.Errorf("%v 갤러리의 %v번째 페이지에서 검색된 글이 %v개 입니다. 25개여야 정상입니다.", gallID(URL), page, len(l.Items))
// 	}

// 	if err := s.ThumbsUp(l.Items[0]); err != nil {
// 		t.Error(err)
// 	}
// 	if err := s.ThumbsDown(l.Items[0]); err != nil {
// 		t.Error(err)
// 	}
// 	if err := s.Report(l.Items[0], "신고"); err != nil {
// 		t.Error(err)
// 	}

// 	a, err := l.Items[0].Fetch()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if a.ThumbsUp == 0 || a.ThumbsDown == 0 {
// 		t.Error("추천 혹은 비추천이 정상적으로 이루어지지 않았습니다.")
// 	}

// 	if err := s.Report(l.Items[0], "신고"); err == nil {
// 		t.Error("게시물 신고가 정상적으로 이루어지지 않았습니다.")
// 	}
// }
