package goinside

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetAllGall 함수는 디시인사이드의 모든 갤러리의 정보를 가져옵니다.
// 마이너 갤러리의 정보는 가져오지 않습니다.
func GetAllGall() ([]*GallInfo, error) {
	galls := []*GallInfo{}
	doc, err := goquery.NewDocument(gallTotalURL) // This URL doesn't check mobile user-agent
	if err != nil {
		return nil, err
	}
	gallDivs := doc.Find(`.gallery_catergory1 > div`)
	gallDivs.Each(func(i int, s *goquery.Selection) {
		a := s.Find(`a`)
		if URL, ok := a.Attr(`href`); ok {
			if ID := a.Text(); ID != "" {
				galls = append(galls, &GallInfo{URL: URL, ID: ID})
			}
		}
	})
	return galls, nil
}

// GetList 함수는 해당 갤러리의 해당 페이지에 있는 모든 글의 목록을 가져옵니다.
func GetList(gallURL string, page int) (*List, error) {
	doc, err := newMobileDoc(fmt.Sprintf("%s&page=%d", gallURL, page))
	if err != nil {
		return nil, err
	}

	list := &List{}

	fnEachList := func(i int, s *goquery.Selection) {
		newArticle := &Article{
			AuthorInfo: &AuthorInfo{
				Name:       fnListGetAuthorName(s),
				IsGuest:    fnListIsAuthorGuest(s),
				GallogIcon: fnListGetAuthorGallogIcon(s),
			},
			Gall: &GallInfo{
				URL: gallURL,
			},
			Icon:         fnListGetArticleIcon(s),
			URL:          fnListGetGallURL(s),
			Number:       fnListGetArticleNumber(s),
			Subject:      fnListGetArticleSubject(s),
			Hit:          fnListGetHit(s),
			ThumbsUp:     fnListGetThumbsUp(s),
			Date:         fnListGetArticleDate(s),
			CommentCount: fnListGetCommentCount(s),
		}
		list.Articles = append(list.Articles, newArticle)
	}

	doc.Find(".article_list > .list_best > li").Each(fnEachList)
	return list, nil
}

// GetArticle 함수는 해당 글의 정보를 가져옵니다.
func GetArticle(articleURL string) (*Article, error) {
	doc, err := newMobileDoc(articleURL)
	if err != nil {
		return nil, err
	}

	s := doc.Find(`body`)

	gallInfo := &GallInfo{
		URL:  fnArticleGetGallURL(s),
		ID:   fnArticleGetGallID(s),
		Name: fnArticleGetGallName(s),
	}

	article := &Article{
		AuthorInfo: &AuthorInfo{
			Name:       fnArticleGetAuthorName(s),
			IP:         fnArticleGetAuthorIP(s),
			IsGuest:    fnArticleIsAuthorGuest(s),
			GallogID:   fnArticleGetAuthorGallogID(s),
			GallogURL:  fnArticleGetAuthorGallogURL(s),
			GallogIcon: fnArticleGetAuthorGallogIcon(s),
		},
		Gall:         gallInfo,
		Icon:         fnArticleGetArticleIcon(s),
		URL:          fnArticleGetArticleURL(s),
		Number:       fnArticleGetArticleNumber(s),
		Subject:      fnArticleGetArticleSubject(s),
		Content:      fnArticleGetArticleContent(s),
		Hit:          fnArticleGetArticleHit(s),
		ThumbsUp:     fnArticleGetArticleThumbsUp(s),
		ThumbsDown:   fnArticleGetArticleThumbsDown(s),
		Date:         fnArticleGetArticleDate(s),
		CommentCount: fnArticleGetArticleCommentCount(s),
	}

	article.Comments = fnArticleGetArticleComments(s, gallInfo, article)
	return article, nil
}

// for List Functions
func fnListGetAuthorName(s *goquery.Selection) string {
	q := `.name`
	return s.Find(q).Text()
}

func fnListIsAuthorGuest(s *goquery.Selection) bool {
	q := `.nick_comm`
	return !s.Find(q).HasClass("nick_comm")
}

func fnListGetAuthorGallogIcon(s *goquery.Selection) string {
	q := `.nick_comm`
	iconElement := s.Find(q)
	for key := range gallogIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func fnListGetArticleIcon(s *goquery.Selection) string {
	q := `.ico_pic`
	iconElement := s.Find(q)
	for key := range iconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func fnListGetGallURL(s *goquery.Selection) string {
	q := `span > a`
	href, _ := s.Find(q).Attr("href")
	return href
}

func fnListGetArticleNumber(s *goquery.Selection) string {
	re := regexp.MustCompile(`no=(\d+)`)
	matched := re.FindStringSubmatch(fnListGetGallURL(s))
	if len(matched) == 2 {
		return matched[1]
	}
	return ""
}

func fnListGetArticleSubject(s *goquery.Selection) string {
	return s.Find(".txt").Text()
}

func fnListGetHit(s *goquery.Selection) int {
	q := `.info > .bar + span > span`
	hit, _ := strconv.Atoi(s.Find(q).Text())
	return hit
}

func fnListGetThumbsUp(s *goquery.Selection) int {
	q := `.info > span:last-of-type > span`
	thumbsUp, _ := strconv.Atoi(s.Find(q).Text())
	return thumbsUp
}

func fnListGetArticleDate(s *goquery.Selection) *time.Time {
	q1 := `.name + span`
	q2 := `.nick_comm + span`
	t := s.Find(q1).Text()
	if t == "" {
		t = s.Find(q2).Text()
	}
	return strToTime(t)
}

func fnListGetCommentCount(s *goquery.Selection) int {
	q := `.txt_num`
	cnt, _ := strconv.Atoi(strings.Trim(s.Find(q).Text(), "[]"))
	return cnt
}

// for Article Functions
func fnArticleGetAuthorName(s *goquery.Selection) string {
	q := `.gall_content .info_edit > span:first-of-type > span:first-of-type`
	return s.Find(q).Text()
}

func fnArticleGetAuthorIP(s *goquery.Selection) string {
	q := `.gall_content .ip`
	return s.Find(q).Text()
}

func fnArticleIsAuthorGuest(s *goquery.Selection) bool {
	q := `.gall_content .nick_comm`
	return !s.Find(q).HasClass("nick_comm")
}

func fnArticleGetAuthorGallogID(s *goquery.Selection) string {
	re := regexp.MustCompile(`id=(\w+)`)
	matched := re.FindStringSubmatch(fnArticleGetAuthorGallogURL(s))
	if len(matched) == 2 {
		return matched[1]
	}
	return ""
}

func fnArticleGetAuthorGallogURL(s *goquery.Selection) string {
	q := `.gall_content .btn.btn_gall`
	href, _ := s.Find(q).Attr("href")
	return href
}

func fnArticleGetAuthorGallogIcon(s *goquery.Selection) string {
	q := `.gall_content .nick_comm`
	iconElement := s.Find(q)
	for key := range gallogIconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func fnArticleGetGallURL(s *goquery.Selection) string {
	q := `.section_info h3 > a`
	href, _ := s.Find(q).Attr("href")
	return href
}

func fnArticleGetGallID(s *goquery.Selection) string {
	q := `input[name="gall_id"]`
	name, _ := s.Find(q).Attr(`value`)
	return name
}

func fnArticleGetGallName(s *goquery.Selection) string {
	q := `input:not([id])[name="ko_name"]`
	name, _ := s.Find(q).Attr(`value`)
	return name
}

func fnArticleGetArticleIcon(s *goquery.Selection) string {
	q := `.article_list .on .ico_pic`
	iconElement := s.Find(q)
	for key := range iconURLMap {
		if iconElement.HasClass(key) {
			return key
		}
	}
	return ""
}

func fnArticleGetArticleURL(s *goquery.Selection) string {
	q := `.article_list .on > a`
	href, _ := s.Find(q).Attr("href")
	return href
}

func fnArticleGetArticleNumber(s *goquery.Selection) string {
	q := `input[name="content_no"]`
	name, _ := s.Find(q).Attr(`value`)
	return name
}

func fnArticleGetArticleSubject(s *goquery.Selection) string {
	q := `.gall_content .tit_view`
	return strings.TrimSpace(s.Find(q).Text())
}

func fnArticleGetArticleContent(s *goquery.Selection) string {
	q := `.gall_content .view_main`
	html, _ := s.Find(q).Html()
	return html
}

func fnArticleGetArticleHit(s *goquery.Selection) int {
	q := `.gall_content .txt_info > .num:first-of-type`
	hit, _ := strconv.Atoi(s.Find(q).Text())
	return hit
}

func fnArticleGetArticleThumbsUp(s *goquery.Selection) int {
	q := `.gall_content #recomm_btn`
	recommend, _ := strconv.Atoi(s.Find(q).Text())
	return recommend
}

func fnArticleGetArticleThumbsDown(s *goquery.Selection) int {
	q := `.gall_content #nonrecomm_btn`
	norecommend, _ := strconv.Atoi(s.Find(q).Text())
	return norecommend
}

func fnArticleGetArticleDate(s *goquery.Selection) *time.Time {
	q := `.gall_content .info_edit > span:first-of-type > span:last-of-type`
	return strToTime(s.Find(q).Text())
}

func fnArticleGetArticleComments(s *goquery.Selection, gallInfo *GallInfo, parents *Article) (cs Comments) {
	ss := []*goquery.Selection{s}
	q := `.list_best .inner_best`
	maxPage := 1
	page := s.Find(`p.paging_page`).Text()
	splited := strings.Split(page, "/")
	if len(splited) == 2 {
		maxPage, _ = strconv.Atoi(strings.TrimSpace(splited[1]))
		for i := 2; i <= maxPage; i++ {
			URL := fmt.Sprintf(`%s?id=%s&no=%s&com_page=%d`, commentMoreURL, gallInfo.ID, parents.Number, i)
			newS, err := newMobileDoc(URL)
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
				AuthorInfo: &AuthorInfo{
					Name:       strings.Trim(s.Find(`.id`).Text(), "[]"),
					IP:         s.Find(`.ip`).Text(),
					IsGuest:    s.Find(`.nick_comm`).Length() == 0,
					GallogID:   gallogID,
					GallogURL:  gallogURL,
					GallogIcon: gallogIcon,
				},
				Gall:    gallInfo,
				Parents: parents,
				Number:  number,
				Content: content,
				Date:    strToTime(s.Find(`.date`).Text()),
			})
		})
	}
	return
}

func fnArticleGetArticleCommentCount(s *goquery.Selection) int {
	q := `.gall_content #comment_dirc`
	hit, _ := strconv.Atoi(s.Find(q).Text())
	return hit
}
