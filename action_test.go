package goinside

import "testing"

func TestMemberAction(t *testing.T) {
	s, err := getTestMemberSession()
	if err != nil {
		t.Fatal(err)
	}

	l, err := FetchList(testTargetGallID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) == 0 {
		t.Errorf("empty %v gallery list.", testTargetGallID)
	}

	// test action to ListItem
	if err := s.ThumbsUp(l.Items[0]); err != nil {
		t.Error(err)
	}
	if err := s.ThumbsDown(l.Items[0]); err != nil {
		t.Error(err)
	}
	// if err := s.Report(l.Items[0], "신고"); err != nil {
	// 	t.Error(err)
	// }

	// test action to Article
	a, err := l.Items[1].Fetch()
	if err != nil {
		t.Error(err)
	}
	if err := s.ThumbsUp(a); err != nil {
		t.Error(err)
	}
	if err := s.ThumbsDown(a); err != nil {
		t.Error(err)
	}
	// if err := s.Report(a, "신고"); err != nil {
	// 	t.Error(err)
	// }
}

func TestGuestAction(t *testing.T) {
	s, err := getTestGuestSession()
	if err != nil {
		t.Fatal(err)
	}

	l, err := FetchList(testTargetGallID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) == 0 {
		t.Errorf("empty %v gallery list.", testTargetGallID)
	}

	// test action to ListItem
	if err := s.ThumbsUp(l.Items[len(l.Items)-2]); err != nil {
		t.Error(err)
	}
	if err := s.ThumbsDown(l.Items[len(l.Items)-2]); err != nil {
		t.Error(err)
	}
	// if err := s.Report(l.Items[len(l.Items)-2], "신고"); err != nil {
	// 	t.Error(err)
	// }

	// test action to Article
	a, err := l.Items[len(l.Items)-1].Fetch()
	if err != nil {
		t.Error(err)
	}
	if err := s.ThumbsUp(a); err != nil {
		t.Error(err)
	}
	if err := s.ThumbsDown(a); err != nil {
		t.Error(err)
	}
	// if err := s.Report(a, "신고"); err != nil {
	// 	t.Error(err)
	// }
}
