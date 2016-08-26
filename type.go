package goinside

import "time"

// 디시인사이드에서 회원의 단계를 구분하는데 사용되는 값입니다.
const (
	UnknownLevel Level = ""
	Level8       Level = "8"
	Level9       Level = "9"
	Level10      Level = "10"
)

// 멤버의 타입을 구분하는 상수입니다.
// FullMemberType은 고정닉, HalfMemberType 반고정닉, GuestMemberType은 유동닉입니다.
const (
	UnknownMemberType MemberType = iota
	FullMemberType
	HalfMemberType
	GuestMemberType
)

// 글의 타입을 구분하는 상수입니다.
// TextArticleType은 텍스트만 있는 글, TextBestArticleType은 텍스트만 있는 개념글,
// ImageArticleType은 이미지만 있는 글, ImageBestArticleType은 이미지만 있는 개념글,
// MovieArticleType은 동영상이 있는 글, SuperBestArticleType은 초개념글입니다.
const (
	UnknownArticleType ArticleType = iota
	TextArticleType
	TextBestArticleType
	ImageArticleType
	ImageBestArticleType
	MovieArticleType
	SuperBestArticleType
)

// 댓글의 타입을 구분하는 상수입니다.
// TextCommentType은 일반 댓글, DCconCommentType은 디시콘,
// VoiceCommentType은 보이스 리플입니다.
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

// MemberType 은 멤버 타입을 나타내는 상수입니다.
type MemberType int

// Level 메소드는 해당 멤버 타입의 레벨을 반환합니다.
// HalfMemberType은 Level8, FullMemberType은 Level9, GuestMemberType은 Level10을
// 반환합니다. 알 수 없는 경우 UnknownLevel을 반환합니다.
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

// Level 은 멤버 레벨을 나타내는 문자열입니다.
type Level string

// Type 메소드는 해당 레벨의 멤버 타입을 반환합니다.
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

// IconURL 메소드는 해당 레벨의 멤버 아이콘 URL을 반환합니다.
func (lv Level) IconURL() string {
	return GallogIconURLMap[lv.Type()]
}

// Gall 구조체는 갤러리의 가장 기본적인 정보, ID와 이름을 나타냅니다.
type Gall struct {
	ID  string
	URL string
}

// List 구조체는 갤러리의 정보와 해당 갤러리의 특정 페이지에 있는 글의 슬라이스를 나타냅니다.
type List struct {
	Info  *ListInfo
	Items []*ListItem
}

// ListInfo 구조체는 갤러리의 정보를 나타냅니다.
// 기본 정보를 나타내는 Gall을 제외한 나머지 정보들은 디시인사이드 갤러리 목록 API가 반환하는 값들입니다.
// CategoryName은 해당 갤러리의 카테고리 이름을 나타냅니다. FileCount와 FileSize는 정확히 무엇을 나타내는지 알 수 없습니다.
type ListInfo struct {
	*Gall
	CategoryName string
	FileCount    string
	FileSize     string
}

// ListItem 구조체는 갤러리 페이지에서 확인할 수 있는 글의 정보를 나타냅니다.
// Article 구조체와는 다릅니다. Article 구조체는 글 본문과 댓글들을 포함하지만
// ListItem 구조체는 그럴 수 없습니다. 갤러리 목록에서 확인할 수 있는 정보만을 나타냅니다.
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
	Number             string
	Date               time.Time
}

// Article 구조체는 글의 정보를 나타냅니다. 여기에는 해당 글이 작성된 갤러리의 정보와
// 해당 글에 달린 댓글들을 모두 포함합니다.
type Article struct {
	*Gall
	URL           string
	Subject       string
	Content       string
	ThumbsUp      int
	ThumbsDown    int
	Name          string
	Number        string
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

// ArticleType 은 글의 타입을 나타내는 상수입니다.
type ArticleType int

// IconURL을 메소드는 해당 글 타입의 IconURL을 반환합니다.
func (a ArticleType) IconURL을() string {
	return ArticleIconURLMap[a]
}

// Comment 구조체는 댓글의 정보를 나타냅니다. 여기에는 해당 댓글이 작성된 글의
// 정보도 포함됩니다. HTML은 해당 댓글의 타입에 맞는 HTML 코드입니다.
// DCcon일 경우 img element, 보이스리플일 경우 audio element 입니다.
type Comment struct {
	*Gall
	Parents   *Article
	Content   string
	HTML      string
	Type      CommentType
	Name      string
	GallogID  string
	GallogURL string
	Level     Level
	IP        string
	Number    string
	Date      time.Time
}

// CommentType 은 댓글의 타입을 나타내는 상수입니다.
type CommentType int

// MajorGallery 구조체는 일반 갤러리의 정보를 나타냅니다.
type MajorGallery struct {
	ID     string
	Name   string
	Number string
}

// MinorGallery 구조체는 마이너 갤러리의 정보를 나타냅니다.
// 마이너 갤러리는 일반 갤러리와 달리 매니저와 부매니저가 존재합니다.
// 부매니저는 여러 명일 수 있습니다. 해당 값들은 gallog ID 입니다.
type MinorGallery struct {
	ID          string
	Name        string
	Number      string
	Manager     string
	SubManagers []string
}
