package goinside

type writable interface {
	write(s session) error
}

// ArticleDraft 구조체는 작성하기 위한 글의 초안을 나타냅니다.
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
	form, contentType := s.articleWriteForm(
		ad.GallID, ad.Subject, ad.Content, ad.Images...)
	resp, err := writeArticleAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

func (i *ListItem) write(s session) error {
	form, contentType := s.articleWriteForm(i.Gall.ID, i.Subject, i.Subject)
	resp, err := writeArticleAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

// CommentDraft 구조체는 작성하기 위한 댓글의 초안을 나타냅니다.
type CommentDraft struct {
	TargetGallID        string
	TargetArticleNumber string
	Content             string
}

type commentable interface {
	articleInfo() (string, string)
}

// NewCommentDraft 함수는 댓글을 작성하기 위해 초안을 생성합니다.
// 해당 초안을 세션의 Write 함수로 전달하여 작성 요청을 보낼 수 있습니다.
func NewCommentDraft(c commentable, content string) *CommentDraft {
	id, number := c.articleInfo()
	return &CommentDraft{id, number, content}
}

func (cd *CommentDraft) write(s session) error {
	form, contentType := s.commentWriteForm(
		cd.TargetGallID, cd.TargetArticleNumber, cd.Content)
	resp, err := writeCommentAPI.post(s, form, contentType)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}
