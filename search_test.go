package goinside

import "testing"

func TestSearch(t *testing.T) {
	l, err := Search("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}
	if len(l.Items) == 0 {
		t.Error("검색 결과가 없습니다.")
	}

	_, err = SearchBySubject("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}

	_, err = SearchByContent("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}

	_, err = SearchBySubjectAndContent("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}

	_, err = SearchByAuthor("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}
}
