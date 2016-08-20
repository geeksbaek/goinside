package goinside

type writable interface {
	write(s session) error
}

// ArticleDraft 구조체는 작성하기 위한 글의 초안을 표현합니다.
type ArticleDraft struct {
	GallID  string
	Subject string
	Content string
	Images  []string
}

// NewArticleDraft 함수는 글을 작성하기 위해 초안을 생성합니다.
// 해당 초안을 세션의 Write 함수로 전달하여 작성 요청을 보낼 수 있습니다.
func NewArticleDraft(gallID, subject, content string, images ...string) *ArticleDraft {
	return &ArticleDraft{gallID, subject, content, images}
}

func (ad *ArticleDraft) write(s session) error {
	form, contentType := s.articleWriteForm(ad)
	resp, err := writeArticleAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

// CommentDraft 구조체는 작성하기 위한 댓글의 초안을 표현합니다.
type CommentDraft struct {
	Target  *Article
	Content string
}

// NewCommentDraft 함수는 댓글을 작성하기 위해 초안을 생성합니다.
// 해당 초안을 세션의 Write 함수로 전달하여 작성 요청을 보낼 수 있습니다.
func NewCommentDraft(article *Article, content string) *CommentDraft {
	return &CommentDraft{article, content}
}

func (cd *CommentDraft) write(s session) error {
	form, contentType := s.commentWriteForm(cd)
	resp, err := writeCommentAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}
