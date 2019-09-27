package goinside

import "testing"

func TestSearch(t *testing.T) {
	l, err := Search("programming", "ㅇㅇ")
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) == 0 {
		t.Error("검색 결과가 없습니다.")
	}
	t.Log("Search Done.")

	_, err = SearchBySubject("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}
	t.Log("SearchBySubject Done.")

	_, err = SearchByContent("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}
	t.Log("SearchByContent Done.")

	_, err = SearchBySubjectAndContent("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}
	t.Log("SearchBySubjectAndContent Done.")

	_, err = SearchByAuthor("programming", "ㅇㅇ")
	if err != nil {
		t.Error(err)
	}
	t.Log("SearchByAuthor Done.")
}
