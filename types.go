package goinside

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Session 구조체는 사용자의 세션을 위해 사용됩니다.
type Session struct {
	id      string
	pw      string
	ip      string
	cookies []*http.Cookie
	isGuest bool
	proxy   func(*http.Request) (*url.URL, error)
	timeout time.Duration
}

// type Sessions []*Session

// AuthorInfo 구조체는 글쓴이에 대한 정보를 표현합니다.
type AuthorInfo struct {
	Name       string
	IP         string
	IsGuest    bool
	GallogID   string
	GallogURL  string
	GallogIcon string
}

func (a *AuthorInfo) String() string {
	f := "Name: %v, IP: %v, IsGuest: %v, GallogID: %v, GallogURL: %v, " + "GallogIcon: %v"
	return fmt.Sprintf(f, a.Name, a.IP, a.IsGuest, a.GallogID,
		a.GallogURL, a.GallogIcon)
}

// GallInfo 구조체는 갤러리에 대한 정보를 표현합니다.
type GallInfo struct {
	URL    string
	ID     string
	Name   string
	detail *gallInfoDetail
}

func (g *GallInfo) String() string {
	f := "URL: %v, ID: %v, Name: %v\nDetail: {%v}"
	return fmt.Sprintf(f, g.URL, g.ID, g.Name, g.detail)
}

type gallInfoDetail struct {
	koName     string
	gServer    string
	gNo        string
	categoryNo string
	ip         string
}

func (g *gallInfoDetail) String() string {
	f := "koName: %v, gServer: %v, gNo: %v, categoryNo: %v, ip: %v"
	return fmt.Sprintf(f, g.koName, g.gServer, g.gNo, g.categoryNo, g.ip)
}

// Comment 구조체는 작성된 댓글에 대한 정보를 표현합니다.
type Comment struct {
	*AuthorInfo
	Gall    *GallInfo
	Parents *Article
	Number  string
	Content string
	Date    *time.Time
}

func (c *Comment) String() string {
	f := "AuthorInfo: {%v}\nGall: {%v}, Name: %v\nContent: %v\nDate: %v"
	return fmt.Sprintf(f, c.AuthorInfo, c.Gall, c.Number, c.Content, c.Date)
}

// type Comments []*Comment

// func (cs Comments) String() string {
// 	var buf bytes.Buffer
// 	for _, c := range cs {
// 		fmt.Fprintln(&buf, c)
// 	}
// 	return buf.String()
// }

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
	Date         *time.Time
	Comments     []*Comment
	CommentCount int
}

func (a *Article) String() string {
	f := "AuthorInfo: {%v}\nGall: {%v}\n" +
		"Icon: %v, URL: %v, Number: %v, Subject: %v\n" +
		"Content: %v\n" + "Hit: %v, ThumbsUp: %v, ThumbsDown: %v, Date: %v\n" +
		"Comments: {%v}\nCommentCount: %v"
	return fmt.Sprintf(f, a.AuthorInfo, a.Gall, a.Icon, a.URL, a.Number,
		a.Subject, a.Content, a.Hit, a.ThumbsUp, a.ThumbsDown, a.Date,
		a.Comments, a.CommentCount)
}

// type Articles []*Article

// func (as Articles) String() string {
// 	var buf bytes.Buffer
// 	for _, a := range as {
// 		fmt.Fprintln(&buf, a)
// 	}
// 	return buf.String()
// }

// List 구조체는 특정 갤러리의 글 묶음입니다.
type List struct {
	Gall     *GallInfo
	Articles []*Article
}

func (l *List) String() string {
	f := "Gall: {%v}\nArticles: {%v}"
	return fmt.Sprintf(f, l.Gall, l.Articles)
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
	target  *Article
	content string
}

type deletable interface {
	delete(*Session) error
}
