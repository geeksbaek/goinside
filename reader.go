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

// FetchGallerys 함수는 디시인사이드의 모든 갤러리의 정보를 가져옵니다.
// 마이너 갤러리의 정보는 가져오지 않습니다.
func FetchGallerys() (galls []*GalleryInfo, err error) {
	// This URL doesn't check mobile user-agent
	doc, err := goquery.NewDocument(gallerysURL)
	if err != nil {
		return
	}
	galls = []*GalleryInfo{}
	gallDivs := doc.Find(`.gallery_catergory1 > div`)
	gallDivs.Each(func(i int, s *goquery.Selection) {
		a := s.Find(`a`)
		if URL, ok := a.Attr(`href`); ok {
			if ID := a.Text(); ID != "" {
				// TODO.
				// GalleryInfoDetail을 가져올 수 있음.
				galls = append(galls, &GalleryInfo{URL: URL, ID: ID})
			}
		}
	})
	return
}

// FetchList 함수는 해당 갤러리의 해당 페이지에 있는 모든 글의 목록을 가져옵니다.
func FetchList(URL string, page int) (l *List, err error) {
	URL = fmt.Sprintf("http://m.dcinside.com/list.php?id=%s&page=%d", _ParseGallID(URL), page)
	doc, err := _NewMobileDocument(URL)
	if err != nil {
		return
	}
	l = &List{Gall: nil}
	eachList := func(i int, s *goquery.Selection) {
		author := &AuthorInfo{
			Name:       _ListAuthorName(s),
			IsGuest:    _ListAuthorIsGuest(s),
			GallogIcon: _ListAuthorGallogIcon(s),
			Detail:     nil,
		}
		gall := &GalleryInfo{
			URL: _MobileURL(URL),
			ID:  _ParseGallID(URL),
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
	}
	doc.Find(".article_list > .list_best > li").Each(eachList)
	return
}

func _ListAuthorName(s *goquery.Selection) string {
	q := `.name`
	return s.Find(q).Text()
}

func _ListAuthorIsGuest(s *goquery.Selection) bool {
	q := `.nick_comm`
	return !s.Find(q).HasClass("nick_comm")
}

func _ListAuthorGallogIcon(s *goquery.Selection) string {
	q := `.nick_comm`
	iconElement := s.Find(q)
	for key := range gallogIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ListArticleIcon(s *goquery.Selection) string {
	q := `.ico_pic`
	iconElement := s.Find(q)
	for key := range iconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ListArticleHasImage(s *goquery.Selection) bool {
	if icon := _ListArticleIcon(s); strings.Contains(icon, "ico_p") {
		return true
	}
	return false
}

func _ListArticleURL(s *goquery.Selection) string {
	q := `span > a`
	href, _ := s.Find(q).Attr("href")
	return href
}

func _ListArticleNumber(s *goquery.Selection) string {
	re := regexp.MustCompile(`no=(\d+)`)
	matched := re.FindStringSubmatch(_ListArticleURL(s))
	if len(matched) == 2 {
		return matched[1]
	}
	return ""
}

func _ListArticleSubject(s *goquery.Selection) string {
	return s.Find(".txt").Text()
}

func _ListArticleHit(s *goquery.Selection) int {
	q := `.info > .bar + span > span`
	hit, _ := strconv.Atoi(s.Find(q).First().Text())
	return hit
}

func _ListArticleThumbsUp(s *goquery.Selection) int {
	q := `.info > span:last-of-type > span`
	thumbsUp, _ := strconv.Atoi(s.Find(q).Text())
	return thumbsUp
}

func _ListArticleDate(s *goquery.Selection) time.Time {
	q1 := `.name + span`
	q2 := `.nick_comm + span`
	t := s.Find(q1).Text()
	if t == "" {
		t = s.Find(q2).Text()
	}
	return _Time(t)
}

func _ListArticleCommentCount(s *goquery.Selection) int {
	q := `.txt_num`
	cnt, _ := strconv.Atoi(strings.Trim(s.Find(q).Text(), "[]"))
	return cnt
}

// FetchArticle 함수는 해당 글의 정보를 가져옵니다.
func FetchArticle(URL string) (a *Article, err error) {
	doc, err := _NewMobileDocument(URL)
	if err != nil {
		return
	}
	s := doc.Find(`body`)
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
	q := `.gall_content .info_edit > span:first-of-type > span:first-of-type`
	return s.Find(q).Text()
}

func _ArticleAuthorDetailIP(s *goquery.Selection) string {
	q := `.gall_content .ip`
	return s.Find(q).Text()
}

func _ArticleAuthorIsGuest(s *goquery.Selection) bool {
	q := `.gall_content .nick_comm`
	return !s.Find(q).HasClass("nick_comm")
}

func _ArticleAuthorDetailGallogID(s *goquery.Selection) string {
	re := regexp.MustCompile(`id=(\w+)`)
	matched := re.FindStringSubmatch(_ArticleAuthorDetailGallogURL(s))
	if len(matched) == 2 {
		return matched[1]
	}
	return ""
}

func _ArticleAuthorDetailGallogURL(s *goquery.Selection) string {
	q := `.gall_content .btn.btn_gall`
	href, _ := s.Find(q).Attr("href")
	return href
}

func _ArticleAuthorGallogIcon(s *goquery.Selection) string {
	q := `.gall_content .nick_comm`
	iconElement := s.Find(q)
	for key := range gallogIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ArticleGalleryInfoURL(s *goquery.Selection) string {
	q := `.section_info h3 > a`
	href, _ := s.Find(q).Attr("href")
	return href
}

func _ArticleGalleryInfoID(s *goquery.Selection) string {
	q := `input[name="gall_id"]`
	name, _ := s.Find(q).Attr(`value`)
	return name
}

func _ArticleGalleryInfoDetailName(s *goquery.Selection) string {
	q := `input:not([id])[name="ko_name"]`
	name, _ := s.Find(q).Attr(`value`)
	return name
}

func _ArticleIcon(s *goquery.Selection) string {
	q := `.article_list .on .ico_pic`
	iconElement := s.Find(q)
	for key := range iconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func _ArticleHasImage(s *goquery.Selection) bool {
	switch _ArticleIcon(s) {
	case "ico_p_c":
		return true
	case "ico_p_y":
		return true
	}
	return false
}

func _ArticleIsBest(s *goquery.Selection) bool {
	switch _ArticleIcon(s) {
	case "ico_p_c":
		return true
	case "ico_t_c":
		return true
	case "ico_sc":
		return true
	}
	return false
}

func _ArticleURL(s *goquery.Selection) string {
	q := `.article_list .on > a`
	href, _ := s.Find(q).Attr("href")
	return href
}

func _ArticleNumber(s *goquery.Selection) string {
	q := `input[name="content_no"]`
	name, _ := s.Find(q).Attr(`value`)
	return name
}

func _ArticleSubject(s *goquery.Selection) string {
	q := `.gall_content .tit_view`
	return strings.TrimSpace(s.Find(q).Text())
}

func _ArticleContent(s *goquery.Selection) string {
	q := `.gall_content .view_main`
	body, _ := s.Find(q).Html()
	return body
}

func _ArticleImages(s *goquery.Selection) (images []string) {
	q := `.gall_content .view_main`
	body, _ := s.Find(q).Html()
	body = html.UnescapeString(body)
	images = []string{}
	imageRe := regexp.MustCompile(`img[^>]*src="([^"]+)"`)
	allMatched := imageRe.FindAllStringSubmatch(body, -1)
	for _, matched := range allMatched {
		if len(matched) >= 2 {
			images = append(images, matched[1])
		}
	}
	return
}

func _ArticleHit(s *goquery.Selection) int {
	q := `.gall_content .txt_info > .num:first-of-type`
	hit, _ := strconv.Atoi(s.Find(q).Text())
	return hit
}

func _ArticleThumbsUp(s *goquery.Selection) int {
	q := `.gall_content #recomm_btn`
	thumbsUp, _ := strconv.Atoi(s.Find(q).Text())
	return thumbsUp
}

func _ArticleThumbsDown(s *goquery.Selection) int {
	q := `.gall_content #nonrecomm_btn`
	thumbsDown, _ := strconv.Atoi(s.Find(q).Text())
	return thumbsDown
}

func _ArticleDate(s *goquery.Selection) time.Time {
	q := `.gall_content .info_edit > span:first-of-type > span:last-of-type`
	return _Time(s.Find(q).Text())
}

func _ArticleComments(s *goquery.Selection, gall *GalleryInfo, parents *Article) (cs []*Comment) {
	ss := []*goquery.Selection{s}
	q := `.list_best .inner_best`
	maxPage := 1
	page := s.Find(`p.paging_page`).Text()
	splited := strings.Split(page, "/")
	if len(splited) == 2 {
		maxPage, _ = strconv.Atoi(strings.TrimSpace(splited[1]))
		for i := 2; i <= maxPage; i++ {
			URL := fmt.Sprintf(`%s?id=%s&no=%s&com_page=%d`, commentMoreURL, gall.ID, parents.Number, i)
			newS, err := _NewMobileDocument(URL)
			if err != nil {
				continue
			}
			ss = append(ss, newS.Find(`.total`))
		}
	}
	for _, s := range ss {
		s.Find(q).Each(func(i int, s *goquery.Selection) {
			var gallogID, gallogURL, gallogIcon string
			gallogURL, _ = s.Find(`a.id`).Attr(`href`)
			idRe := regexp.MustCompile(`id=(\w+)`)
			matchedGallogID := idRe.FindStringSubmatch(gallogURL)
			if len(matchedGallogID) == 2 {
				gallogID = matchedGallogID[1]
			}
			iconElement := s.Find(`.nick_comm`)
			for key := range gallogIconURLMap {
				if iconElement.HasClass(key) {
					gallogIcon = key
					break
				}
			}
			var number, content string
			delhref, _ := s.Find(`.btn_delete`).Attr(`href`)
			numberRe := regexp.MustCompile(`\('(\d+)'`)
			matchedNumber := numberRe.FindStringSubmatch(delhref)
			if len(matchedNumber) == 2 {
				number = matchedNumber[1]
			}
			content, _ = s.Find(`.txt`).Html()
			cs = append(cs, &Comment{
				Author: &AuthorInfo{
					Name:       strings.Trim(s.Find(`.id`).Text(), "[]"),
					IsGuest:    s.Find(`.nick_comm`).Length() == 0,
					GallogIcon: gallogIcon,
					Detail: &AuthorInfoDetail{
						IP:        s.Find(`.ip`).Text(),
						GallogID:  gallogID,
						GallogURL: gallogURL,
					},
				},
				Gall:    gall,
				Parents: parents,
				Number:  number,
				Content: content,
				Date:    _Time(s.Find(`.date`).Text()),
			})
		})
	}
	return
}

func _ArticleCommentCount(s *goquery.Selection) int {
	q := `.gall_content #comment_dirc`
	cnt, _ := strconv.Atoi(s.Find(q).Text())
	return cnt
}
