package goinside

import "testing"

func TestFetch(t *testing.T) {
	URL := "http://gall.dcinside.com/board/lists/?id=baseball_new4"
	page := 1

	l, err := FetchList(URL, page)
	if err != nil {
		t.Error(err)
	}
	if len(l.Items) != 25 {
		t.Errorf("%v 갤러리의 %v번째 페이지에서 검색된 글이 %v개 입니다. 25개여야 정상입니다.", gallID(URL), page, len(l.Items))
	}
	targetArticle := l.Items[0]
	a, err := targetArticle.Fetch()
	if err != nil {
		t.Error(err)
	}
	if targetArticle.Subject != a.Subject {
		t.Errorf("%v 갤러리의 첫 번째 글을 정상적으로 오지 못했습니다", gallID(URL))
	}
	for _, v := range l.Items {
		imageURLs, err := v.FetchImageURLs()
		if (v.HasImage && err != nil) || (v.HasImage && len(imageURLs) == 0) {
			t.Errorf("%v 에서 이미지를 정상적으로 가져오지 못했습니다.", v.URL)
		}
	}
}
