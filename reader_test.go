package goinside

import (
	"errors"
	"testing"
)

func TestFetch(t *testing.T) {
	URL := "http://gall.dcinside.com/board/lists/?id=programming"
	page := 1

	l, err := FetchBestList(URL, page)
	if err != nil {
		t.Error(err)
	}
	if len(l.Items) != 25 {
		t.Errorf("%v 갤러리의 %v번째 페이지에서 검색된 글이 %v개 입니다. 25개여야 정상입니다.", gallID(URL), page, len(l.Items))
	}

	for _, v := range l.Items {
		_, err := v.Fetch()
		if err != nil {
			t.Errorf("%v article fetch failed. %v", v.URL, err)
		}
	}
}

func TestImageURLTypeFetch(t *testing.T) {
	URL := "http://gall.dcinside.com/board/lists/?id=programming"
	page := 1

	l, err := FetchBestList(URL, page)
	if err != nil {
		t.Error(err)
	}
	for _, v := range l.Items {
		if !v.HasImage {
			continue
		}

		is, err := v.FetchImageURLs()
		if err != nil {
			t.Errorf("%v ListItem.FetchImageURLs() failed. %v", v.URL, err)
		}
		if len(is) == 0 {
			t.Errorf("%v ListItem.FetchImageURLs() failed. %v", v.URL, errors.New("Empty Article Image"))
		}

		for _, i := range is {
			if _, err := i.Fetch(); err != nil {
				t.Errorf("%v ImageURLType.Fetch() failed. %v", v.URL, err)
			}
			return
		}
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
