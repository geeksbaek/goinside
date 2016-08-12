package goinside

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// selectors
const (
	gallDivsQuery = `.gallery_catergory1 > div`
	gallListQuery = `.article_list > .list_best > li`

	listAuthorNameQuery          = `.name`
	listAuthorIsGuestQuery       = `.nick_comm`
	listAuthorGallogIconQuery    = `.nick_comm`
	listArticleIconQuery         = `.ico_pic`
	listArticleURLQuery          = `span > a`
	listArticleSubjectQuery      = `.txt`
	listArticleHitQuery          = `.info > .bar + span > span`
	listArticleThumbsUpQuery     = `.info > span:last-of-type > span`
	listArticleGuestDateQuery    = `.name + span`
	listArticleMemberDateQuery   = `.nick_comm + span`
	listArticleCommentCountQuery = `.txt_num`

	articleBodyQuery                  = `body`
	articleAuthorNameQuery            = `.gall_content .info_edit > span:first-of-type > span:first-of-type`
	articleAuthorDetailIPQuery        = `.gall_content .ip`
	articleAuthorIsGuestQuery         = `.gall_content .nick_comm`
	articleAuthorDetailGallogURLQuery = `.gall_content .btn.btn_gall`
	articleAuthorGallogIconQuery      = `.gall_content .nick_comm`
	articleGalleryInfoURLQuery        = `.section_info h3 > a`
	articleGalleryInfoIDQuery         = `input[name="gall_id"]`
	articleGalleryInfoDetailNameQuery = `input:not([id])[name="ko_name"]`
	articleIconQuery                  = `.article_list .on .ico_pic`
	articleURLQuery                   = `.article_list .on > a`
	articleNumberQuery                = `input[name="content_no"]`
	articleSubjectQuery               = `.gall_content .tit_view`
	articleContentQuery               = `.gall_content .view_main`
	articleHitQuery                   = `.gall_content .txt_info > .num:first-of-type`
	articleThumbsUpQuery              = `.gall_content #recomm_btn`
	articleThumbsDownQuery            = `.gall_content #nonrecomm_btn`
	articleDateQuery                  = `.gall_content .info_edit > span:first-of-type > span:last-of-type`

	commentsQuery              = `.list_best .inner_best`
	commentPageQuery           = `p.paging_page`
	commentRestQuery           = `.total`
	commentGallogURLQuery      = `a.id`
	commentGallogIconQuery     = `.nick_comm`
	commentDeleteURLQuery      = `.btn_delete`
	commentContentQuery        = `.txt`
	commentAuthorDetailIPQuery = `.ip`
	commentAuthorNameQuery     = `.id`
	commentDateQuery           = `.date`
	commentCountQuery          = `.gall_content #comment_dirc`
)

// FetchGallerys 함수는 디시인사이드의 모든 갤러리의 정보를 가져옵니다.
// 마이너 갤러리의 정보는 가져오지 않습니다.
func FetchGallerys() (galls []*GalleryInfo, err error) {
	doc, err := newMobileDocument(gallerysURL)
	if err != nil {
		return
	}
	galls = []*GalleryInfo{}
	doc.Find(gallDivsQuery).Each(func(i int, s *goquery.Selection) {
		a := s.Find(`a`)
		if URL, ok := a.Attr(`href`); ok {
			if ID := a.Text(); ID != "" {
				// TODO.
				// GalleryInfoDetail을 가져올 수 있음.
				galls = append(galls, &GalleryInfo{
					URL: URL,
					ID:  ID,
				})
			}
		}
	})
	return
}

// FetchList 함수는 해당 갤러리의 해당 페이지에 있는 모든 글의 목록을 가져옵니다.
func FetchList(URL string, page int) (l *List, err error) {
	URL = fmt.Sprintf("%v&page=%v", URL, page)
	doc, err := newMobileDocument(URL)
	if err != nil {
		return
	}
	l = &List{Gall: nil}
	doc.Find(gallListQuery).Each(func(i int, s *goquery.Selection) {
		author := &AuthorInfo{
			Name:       _ListAuthorName(s),
			IsGuest:    _ListAuthorIsGuest(s),
			GallogIcon: _ListAuthorGallogIcon(s),
			Detail:     nil,
		}
		gall := &GalleryInfo{
			URL: ToMobileURL(URL),
			ID:  gallID(URL),
		}
		article := &Article{
			Author:       author,
			Gall:         gall,
			ArticleIcon:  _ListArticleIcon(s),
			HasImage:     _ListArticleHasImage(s),
			URL:          _ListArticleURL(s),
			Number:       _ListArticleNumber(s),
			Subject:      _ListArticleSubject(s),
			Hit:          _ListArticleHit(s),
			ThumbsUp:     _ListArticleThumbsUp(s),
			Date:         _ListArticleDate(s),
			CommentCount: _ListArticleCommentCount(s),
			Detail:       nil,
		}
		l.Articles = append(l.Articles, article)
	})
	return
}

func _ListAuthorName(s *goquery.Selection) string {
	return s.Find(listAuthorNameQuery).Text()
}

func _ListAuthorIsGuest(s *goquery.Selection) bool {
	return !s.Find(listAuthorIsGuestQuery).HasClass("nick_comm")
}

func _ListAuthorGallogIcon(s *goquery.Selection) string {
	iconElement := s.Find(listAuthorGallogIconQuery)
	for key := range GallogIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ListArticleIcon(s *goquery.Selection) string {
	iconElement := s.Find(listArticleIconQuery)
	for key := range ArticleIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ListArticleHasImage(s *goquery.Selection) bool {
	switch _ListArticleIcon(s) {
	case "ico_p_c", "ico_p_y":
		return true
	}
	return false
}

func _ListArticleURL(s *goquery.Selection) string {
	href, _ := s.Find(listArticleURLQuery).Attr("href")
	return href
}

func _ListArticleNumber(s *goquery.Selection) string {
	return articleNumber(_ListArticleURL(s))
}

func _ListArticleSubject(s *goquery.Selection) string {
	return s.Find(listArticleSubjectQuery).Text()
}

func _ListArticleHit(s *goquery.Selection) int {
	hit, _ := strconv.Atoi(s.Find(listArticleHitQuery).First().Text())
	return hit
}

func _ListArticleThumbsUp(s *goquery.Selection) int {
	thumbsUp, _ := strconv.Atoi(s.Find(listArticleThumbsUpQuery).Text())
	return thumbsUp
}

func _ListArticleDate(s *goquery.Selection) time.Time {
	if _ListAuthorIsGuest(s) {
		return timeFormatting(s.Find(listArticleGuestDateQuery).Text())
	}
	return timeFormatting(s.Find(listArticleMemberDateQuery).Text())
}

func _ListArticleCommentCount(s *goquery.Selection) int {
	commentCount := s.Find(listArticleCommentCountQuery).Text()
	count, _ := strconv.Atoi(strings.Trim(commentCount, "[]"))
	return count
}

// FetchArticle 함수는 해당 글의 정보를 가져옵니다.
func FetchArticle(URL string) (a *Article, err error) {
	doc, err := newMobileDocument(URL)
	if err != nil {
		return
	}
	s := doc.Find(articleBodyQuery)
	author := &AuthorInfo{
		Name:       _ArticleAuthorName(s),
		IsGuest:    _ArticleAuthorIsGuest(s),
		GallogIcon: _ArticleAuthorGallogIcon(s),
		Detail: &AuthorInfoDetail{
			IP:        _ArticleAuthorDetailIP(s),
			GallogID:  _ArticleAuthorDetailGallogID(s),
			GallogURL: _ArticleAuthorDetailGallogURL(s),
		},
	}
	gall := &GalleryInfo{
		URL: _ArticleGalleryInfoURL(s),
		ID:  _ArticleGalleryInfoID(s),
		Detail: &GalleryInfoDetail{
			Name: _ArticleGalleryInfoDetailName(s),
		},
	}
	a = &Article{
		Author:       author,
		Gall:         gall,
		ArticleIcon:  _ArticleIcon(s),
		HasImage:     _ArticleHasImage(s),
		IsBest:       _ArticleIsBest(s),
		URL:          _ArticleURL(s),
		Number:       _ArticleNumber(s),
		Subject:      _ArticleSubject(s),
		Hit:          _ArticleHit(s),
		ThumbsUp:     _ArticleThumbsUp(s),
		CommentCount: _ArticleCommentCount(s),
		Date:         _ArticleDate(s),
		Detail: &ArticleDetail{
			Content:    _ArticleContent(s),
			ImageURLs:  _ArticleImages(s),
			ThumbsDown: _ArticleThumbsDown(s),
		},
	}
	a.Detail.Comments = _ArticleComments(s, gall, a) // for set parents
	return
}

func _ArticleAuthorName(s *goquery.Selection) string {
	return s.Find(articleAuthorNameQuery).Text()
}

func _ArticleAuthorDetailIP(s *goquery.Selection) string {
	return s.Find(articleAuthorDetailIPQuery).Text()
}

func _ArticleAuthorIsGuest(s *goquery.Selection) bool {
	return !s.Find(articleAuthorIsGuestQuery).HasClass("nick_comm")
}

func _ArticleAuthorDetailGallogID(s *goquery.Selection) string {
	return gallID(_ArticleAuthorDetailGallogURL(s))
}

func _ArticleAuthorDetailGallogURL(s *goquery.Selection) string {
	href, _ := s.Find(articleAuthorDetailGallogURLQuery).Attr("href")
	return href
}

func _ArticleAuthorGallogIcon(s *goquery.Selection) string {
	iconElement := s.Find(articleAuthorGallogIconQuery)
	for key := range GallogIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ArticleGalleryInfoURL(s *goquery.Selection) string {
	href, _ := s.Find(articleGalleryInfoURLQuery).Attr("href")
	return href
}

func _ArticleGalleryInfoID(s *goquery.Selection) string {
	name, _ := s.Find(articleGalleryInfoIDQuery).Attr(`value`)
	return name
}

func _ArticleGalleryInfoDetailName(s *goquery.Selection) string {
	name, _ := s.Find(articleGalleryInfoDetailNameQuery).Attr(`value`)
	return name
}

func _ArticleIcon(s *goquery.Selection) string {
	iconElement := s.Find(articleIconQuery)
	for key := range ArticleIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ArticleHasImage(s *goquery.Selection) bool {
	switch _ArticleIcon(s) {
	case "ico_p_c", "ico_p_y":
		return true
	}
	return false
}

func _ArticleIsBest(s *goquery.Selection) bool {
	switch _ArticleIcon(s) {
	case "ico_p_c", "ico_t_c", "ico_sc":
		return true
	}
	return false
}

func _ArticleURL(s *goquery.Selection) string {
	href, _ := s.Find(articleURLQuery).Attr("href")
	return href
}

func _ArticleNumber(s *goquery.Selection) string {
	name, _ := s.Find(articleNumberQuery).Attr(`value`)
	return name
}

func _ArticleSubject(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(articleSubjectQuery).Text())
}

func _ArticleContent(s *goquery.Selection) string {
	body, _ := s.Find(articleContentQuery).Html()
	return body
}

func _ArticleImages(s *goquery.Selection) (images []string) {
	body, _ := s.Find(articleContentQuery).Html()
	body = html.UnescapeString(body)
	return imageElements(body)
}

func _ArticleHit(s *goquery.Selection) int {
	hit, _ := strconv.Atoi(s.Find(articleHitQuery).Text())
	return hit
}

func _ArticleThumbsUp(s *goquery.Selection) int {
	thumbsUp, _ := strconv.Atoi(s.Find(articleThumbsUpQuery).Text())
	return thumbsUp
}

func _ArticleThumbsDown(s *goquery.Selection) int {
	thumbsDown, _ := strconv.Atoi(s.Find(articleThumbsDownQuery).Text())
	return thumbsDown
}

func _ArticleDate(s *goquery.Selection) time.Time {
	return timeFormatting(s.Find(articleDateQuery).Text())
}

func _ArticleComments(s *goquery.Selection, gall *GalleryInfo, parents *Article) (cs []*Comment) {
	ss := []*goquery.Selection{s}
	maxPage := 1
	page := s.Find(commentPageQuery).Text()
	splited := strings.Split(page, "/")
	if len(splited) == 2 {
		maxPage, _ = strconv.Atoi(strings.TrimSpace(splited[1]))
		for i := 2; i <= maxPage; i++ {
			URL := mobileCommentPageURL(gall.ID, parents.Number, i)
			newS, err := newMobileDocument(URL)
			if err != nil {
				continue
			}
			ss = append(ss, newS.Find(commentRestQuery))
		}
	}
	for _, s := range ss {
		s.Find(commentsQuery).Each(func(i int, s *goquery.Selection) {
			var gallogID, gallogURL, gallogIcon string
			gallogURL, _ = s.Find(commentGallogURLQuery).Attr(`href`)
			gallogID = gallID(gallogURL)
			iconElement := s.Find(commentGallogIconQuery)
			for key := range GallogIconURLMap {
				if iconElement.HasClass(key) {
					gallogIcon = key
					break
				}
			}
			var number, content string
			delhref, _ := s.Find(commentDeleteURLQuery).Attr(`href`)
			numberRe := regexp.MustCompile(`\('(\d+)'`)
			matchedNumber := numberRe.FindStringSubmatch(delhref)
			if len(matchedNumber) == 2 {
				number = matchedNumber[1]
			}
			content, _ = s.Find(commentContentQuery).Html()
			authorDetail := &AuthorInfoDetail{
				IP:        s.Find(commentAuthorDetailIPQuery).Text(),
				GallogID:  gallogID,
				GallogURL: gallogURL,
			}
			author := &AuthorInfo{
				Name:       strings.Trim(s.Find(commentAuthorNameQuery).Text(), "[]"),
				IsGuest:    s.Find(commentGallogIconQuery).Length() == 0,
				GallogIcon: gallogIcon,
				Detail:     authorDetail,
			}
			comment := &Comment{
				Author:  author,
				Gall:    gall,
				Parents: parents,
				Number:  number,
				Content: content,
				Date:    timeFormatting(s.Find(commentDateQuery).Text()),
			}
			cs = append(cs, comment)
		})
	}
	return
}

func _ArticleCommentCount(s *goquery.Selection) int {
	cnt, _ := strconv.Atoi(s.Find(commentCountQuery).Text())
	return cnt
}
