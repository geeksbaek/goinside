package goinside

import "time"

const (
	UnknownLevel Level = ""
	Level8       Level = "8"
	Level9       Level = "9"
	Level10      Level = "10"
)

const (
	UnknownMemberType MemberType = iota
	FullMemberType
	HalfMemberType
	GuestMemberType
)

const (
	UnknownArticleType ArticleType = iota
	TextArticleType
	TextBestArticleType
	ImageArticleType
	ImageBestArticleType
	MovieArticleType
	SuperBestArticleType
)

const (
	UnknownCommentType CommentType = iota
	TextCommentType
	DCconCommentType
	VoiceCommentType
)

type jsonValidation []struct {
	Result bool   `json:"result"`
	Cause  string `json:"cause"`
}

type MemberType int

func (m MemberType) Level() Level {
	switch m {
	case HalfMemberType:
		return Level8
	case FullMemberType:
		return Level9
	case GuestMemberType:
		return Level10
	}
	return UnknownLevel
}

type Level string

func (lv Level) Type() MemberType {
	switch lv {
	case Level8:
		return HalfMemberType
	case Level9:
		return FullMemberType
	case Level10:
		return GuestMemberType
	}
	return UnknownMemberType
}

func (lv Level) IconURL() string {
	return GallogIconURLMap[lv.Type()]
}

type Gall struct {
	ID  string
	URL string
}

type List struct {
	Info  *ListInfo
	Items []*ListItem
}

type ListInfo struct {
	*Gall
	CategoryName string
	FileCount    string
	FileSize     string
}

type ListItem struct {
	*Gall
	URL                string
	Subject            string
	Name               string
	Level              Level
	HasImage           bool
	ArticleType        ArticleType
	ThumbsUp           int
	IsBest             bool
	Hit                int
	GallogID           string
	GallogURL          string
	IP                 string
	CommentLength      int
	VoiceCommentLength int
	Number             int
	Date               time.Time
}

type Article struct {
	*Gall
	URL           string
	Subject       string
	Content       string
	ThumbsUp      int
	ThumbsDown    int
	Name          string
	Number        int
	Level         Level
	IP            string
	CommentLength int
	HasImage      bool
	Hit           int
	ArticleType   ArticleType
	GallogID      string
	GallogURL     string
	IsBest        bool
	ImageURLs     []string
	Comments      []*Comment
	Date          time.Time
}

type ArticleType int

func (a ArticleType) URL() string {
	return ArticleIconURLMap[a]
}

type Comment struct {
	*Gall
	Parents     *Article
	Content     string
	HTMLContent string
	Type        CommentType
	Name        string
	GallogID    string
	GallogURL   string
	Level       Level
	IP          string
	Number      int
	Date        time.Time
}
type CommentType int
