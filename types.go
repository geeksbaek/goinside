package goinside

import "net/http"

// Session 구조체는 사용자의 세션을 위해 사용됩니다.
type Session struct {
	id        string
	pw        string
	ip        string
	cookies   []*http.Cookie
	nomember  bool
	transport *http.Transport
}

// AuthorInfo 구조체는 글쓴이에 대한 정보를 표현합니다.
type AuthorInfo struct {
	Name       string
	IP         string
	IsGuest    bool
	GallogID   string
	GallogURL  string
	GallogIcon string
}

// GallInfo 구조체는 갤러리에 대한 정보를 표현합니다.
type GallInfo struct {
	URL  string
	ID   string
	Name string
}

// Comment 구조체는 작성된 댓글에 대한 정보를 표현합니다.
type Comment struct {
	*AuthorInfo
	Gall    *GallInfo
	Parents *Article
	Number  string
	Content string
	Date    string
}

// Article 구조체는 작성된 글에 대한 정보를 표현합니다.
// 댓글을 달거나 추천, 비추천 할 때 사용합니다.
type Article struct {
	*AuthorInfo
	Gall         *GallInfo
	Icon         string
	URL          string
	Number       string
	Subject      string
	Content      string
	Hit          int
	ThumbsUp     int
	ThumbsDown   int
	Date         string
	Comments     []*Comment
	CommentCount int
}

// List 구조체는 특정 갤러리의 글 묶음입니다.
type List struct {
	*GallInfo
	Articles []*Article
}

// ArticleWriter 구조체는 글 작성에 필요한 정보를 전달하기 위한 구조체입니다.
type articleWriter struct {
	*Session
	gall    *GallInfo
	subject string
	content string
	images  []string
}

// CommentWriter 구조체는 댓글 작성에 필요한 정보를 전달하기 위한 구조체입니다.
type commentWriter struct {
	*Session
	*Article
	content string
}
