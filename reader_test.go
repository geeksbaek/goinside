package goinside

import "testing"

func TestFetch(t *testing.T) {
	URL := "http://gall.dcinside.com/board/lists/?id=baseball_new5"
	page := 1

	l, err := FetchBestList(URL, page)
	if err != nil {
		t.Error(err)
	}
	if len(l.Items) != 25 {
		t.Errorf("%v 갤러리의 %v번째 페이지에서 검색된 글이 %v개 입니다. 25개여야 정상입니다.", gallID(URL), page, len(l.Items))
	}
	targetArticle := l.Items[0]
	a, err := targetArticle.Fetch()
	if err != nil {
		t.Error(err, targetArticle.URL)
	}
	if targetArticle.Subject != a.Subject {
		t.Errorf("%v 갤러리의 첫 번째 글을 정상적으로 오지 못했습니다", gallID(URL))
	}
}

func TestFetchGalleryList(t *testing.T) {
	major, err := FetchAllMajorGallery()
	if err != nil {
		t.Error(err)
	}
	if len(major) == 0 {
		t.Error("메이저 갤러리 목록을 가져올 수 없습니다.")
	}
	minor, err := FetchAllMinorGallery()
	if err != nil {
		t.Error(err)
	}
	if len(minor) == 0 {
		t.Error("마이너 갤러리 목록을 가져올 수 없습니다.")
	}
}
