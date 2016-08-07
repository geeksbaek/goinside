package goinside

import "io"

type session interface {
	articleWriteForm(*ArticleDraft) (form io.Reader, contentType string)
	articleDeleteForm(*Article) (form io.Reader, contentType string)
	commentWriteForm(*CommentDraft) (form io.Reader, contentType string)
	commentDeleteForm(*Comment) (form io.Reader, contentType string)
	actionForm(*Article) (form io.Reader, contentType string)
	reportForm(string, string) (form io.Reader, contentType string)
	connector
}

// AuthorInfo 구조체는 글쓴이의 정보를 표현합니다.
type AuthorInfo struct {
	Name       string
	IsGuest    bool
	GallogIcon string
	Detail     *AuthorInfoDetail
}

// AuthorInfoDetail 구조체는 글쓴이의 세부적인 정보를 표현합니다.
type AuthorInfoDetail struct {
	IP        string
	GallogID  string
	GallogURL string
}

// List 구조체는 갤러리의 페이지에 대한 정보를 표현합니다.
type List struct {
	Gall     *GalleryInfo
	Articles []*Article
}

// GalleryInfo 구조체는 갤러리의 정보를 표현합니다.
type GalleryInfo struct {
	URL    string
	ID     string
	Detail *GalleryInfoDetail
}

// GalleryInfoDetail 구조체는 갤러리의 세부적인 정보를 표현합니다.
type GalleryInfoDetail struct {
	Name string
}

type _JSONResponse struct {
	Result bool
	Cause  string
}
