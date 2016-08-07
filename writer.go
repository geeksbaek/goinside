package goinside

import "time"

type writable interface {
	write(s session) error
}

type ArticleDraft struct {
	GallID  string
	Subject string
	Content string
	Images  []string
}

func NewArticleDraft(gallID, subject, content string, images ...string) *ArticleDraft {
	return &ArticleDraft{gallID, subject, content, images}
}

// Article 구조체는 작성된 글의 정보를 표현합니다.
type Article struct {
	Author       *AuthorInfo
	Gall         *GalleryInfo
	ArticleIcon  string
	HasImage     bool
	URL          string
	Number       string
	Subject      string
	Hit          int
	ThumbsUp     int
	Date         time.Time
	CommentCount int
	Detail       *ArticleDetail
}

// ArticleDetail 구조체는 글의 세부적인 정보를 표현합니다.
type ArticleDetail struct {
	Content    string
	ImageURLs  []string
	ThumbsDown int
	Comments   []*Comment
}

func (ad *ArticleDraft) write(s session) error {
	form, contentType := s.articleWriteForm(ad)
	resp, err := api(s, articleWriteAPI, form, contentType)
	if err != nil {
		return err
	}
	return _CheckResponse(resp)
}

type CommentDraft struct {
	Target  *Article
	Content string
}

func NewCommentDraft(article *Article, content string) *CommentDraft {
	return &CommentDraft{article, content}
}

// Comment 구조체는 작성된 댓글에 대한 정보를 표현합니다.
type Comment struct {
	Author  *AuthorInfo
	Gall    *GalleryInfo
	Parents *Article
	Number  string
	Content string
	Date    time.Time
}

func (cd *CommentDraft) write(s session) error {
	form, contentType := s.commentWriteForm(cd)
	resp, err := api(s, commentWriteAPI, form, contentType)
	if err != nil {
		return err
	}
	return _CheckResponse(resp)
}
