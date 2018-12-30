package goinside

import (
	"errors"
	"testing"
)

func TestFetch(t *testing.T) {
	l, err := FetchList(testTargetGallID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) == 0 {
		t.Fatalf("empty %v gallery list.", testTargetGallID)
	}

	for _, v := range l.Items {
		_, err := v.Fetch()
		if err != nil {
			t.Errorf("%v Article.Fetch() failed. %v", v.URL, err)
		}
	}
}

func TestImageURLTypeFetch(t *testing.T) {
	l, err := FetchBestList(testTargetGallID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) == 0 {
		t.Fatalf("empty %v gallery list.", testTargetGallID)
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
			if _, _, err := i.Fetch(); err != nil {
				t.Errorf("%v ImageURLType.Fetch() failed. %v", v.URL, err)
			}
			return
		}
	}
}

func TestFetchGalleryList(t *testing.T) {
	major, err := FetchAllMajorGallery()
	if err != nil {
		t.Fatalf("FetchAllMajorGallery() failed. %v", err)
	}
	if len(major) == 0 {
		t.Errorf("empty major gallery result. %v", err)
	}
	minor, err := FetchAllMinorGallery()
	if err != nil {
		t.Errorf("FetchAllMinorGallery() failed. %v", err)
	}
	if len(minor) == 0 {
		t.Errorf("empty minor gallery result. %v", err)
	}
}
