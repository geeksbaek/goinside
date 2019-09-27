package goinside

import "testing"

func TestMemberWrite(t *testing.T) {
	s, err := getTestMemberSession()
	if err != nil {
		t.Fatal(err)
	}

	// test write article
	ad := NewArticleDraft(testTargetGallID, "ㅋㅋㅋㅋ", "ㅋㅋㅋㅋㅋㅋ", "test.jpg")
	if err := s.Write(ad); err != nil {
		t.Error(err)
	}

	// test write comment to ListItem
	l, err := FetchBestList(testTargetGallID, 1)
	if err != nil {
		t.Fatal(err)
	}
	cd := NewCommentDraft(l.Items[0], "..")
	if err := s.Write(cd); err != nil {
		t.Error(err)
	}

	// test write comment to Article
	a, err := l.Items[0].Fetch()
	if err != nil {
		t.Error(err)
	}
	cd = NewCommentDraft(a, ".")
	if err := s.Write(cd); err != nil {
		t.Error(err)
	}
}

func TestGuestWrite(t *testing.T) {
	s, err := getTestGuestSession()
	if err != nil {
		t.Fatal(err)
	}

	// test write article
	ad := NewArticleDraft(testTargetGallID, "ㅋㅋㅋㅋ", "ㅋㅋㅋㅋㅋㅋ", "test.jpg")
	if err := s.Write(ad); err != nil {
		t.Error(err)
	}

	// // test write comment to ListItem
	// l, err := FetchBestList(testTargetGallID, 1)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// cd := NewCommentDraft(l.Items[0], "..")
	// if err := s.Write(cd); err != nil {
	// 	t.Error(err)
	// }

	// // test write comment to Article
	// a, err := l.Items[0].Fetch()
	// if err != nil {
	// 	t.Error(err)
	// }
	// cd = NewCommentDraft(a, ".")
	// if err := s.Write(cd); err != nil {
	// 	t.Error(err)
	// }
}
