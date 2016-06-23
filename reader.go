package goinside

import "github.com/PuerkitoBio/goquery"

// GetAllGall 함수는 디시인사이드의 모든 갤러리의 정보를 가져옵니다.
// 마이너 갤러리의 정보는 가져오지 않습니다.
func GetAllGall() ([]*GallInfo, error) {
	galls := []*GallInfo{}
	doc, err := goquery.NewDocument(GallTotalURL) // This URL doesn't check mobile user-agent
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
// func GetList(gallURL string, page int) ([]*Article, error) {
// 	resp, err := (&Session{}).get(gallURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	doc, err := goquery.NewDocumentFromResponse(resp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	list := []List{}

// 	_getURL := func(s *goquery.Selection) string {
// 		href, _ := s.Find("span > a").Attr("href")
// 		return href
// 	}

// 	_getPostIconURL := func(s *goquery.Selection) string {
// 		iconElement := s.Find(".ico_pic")
// 		for key, value := range iconURLsMap {
// 			if iconElement.HasClass(key) {
// 				return value
// 			}
// 		}
// 		return ""
// 	}

// 	_getGallogIconURL := func(s *goquery.Selection) string {
// 		iconElement := s.Find(".nick_comm")
// 		for key, value := range gallogIconURLsMap {
// 			if iconElement.HasClass(key) {
// 				return value
// 			}
// 		}
// 		return ""
// 	}

// 	_getSubject := func(s *goquery.Selection) string {
// 		return s.Find(".txt").Text()
// 	}

// 	_getName := func(s *goquery.Selection) string {
// 		return s.Find(".name").Text()
// 	}

// 	_getDate := func(s *goquery.Selection) string {
// 		return s.Find(".name + span").Text()
// 	}

// 	_getHit := func(s *goquery.Selection) int {
// 		hit, _ := strconv.Atoi(s.Find(".info > .bar + span > span").Text())
// 		return hit
// 	}

// 	_getComment := func(s *goquery.Selection) int {
// 		r := strings.NewReplacer("[", "", "]", "")
// 		hit, _ := strconv.Atoi(r.Replace(s.Find(".txt_num").Text()))
// 		return hit
// 	}

// 	_getRecommend := func(s *goquery.Selection) int {
// 		recommend, _ := strconv.Atoi(s.Find(".info > span:last-of-type > span").Text())
// 		return recommend
// 	}

// 	_getIsGoJungNick := func(s *goquery.Selection) bool {
// 		return s.Find(".nick_comm").HasClass("nick_comm")
// 	}

// 	_eachList := func(i int, s *goquery.Selection) {
// 		new := List{
// 			URL:           _getURL(s),
// 			PostIconURL:   _getPostIconURL(s),
// 			GallogIconURL: _getGallogIconURL(s),
// 			Subject:       _getSubject(s),
// 			Name:          _getName(s),
// 			Date:          _getDate(s),
// 			Hit:           _getHit(s),
// 			Comment:       _getComment(s),
// 			Recommend:     _getRecommend(s),
// 			IsGoJungNick:  _getIsGoJungNick(s),
// 		}

// 		list = append(list, new)
// 	}

// 	doc.Find(".article_list > .list_best > li").Each(_eachList)

// 	return nil, nil
// }

// import (
// 	"strconv"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// )

// type ArticleReader struct {
// 	Subject       string
// 	Content       string
// 	Ip            string
// 	Date          string
// 	Name          string
// 	GallogURL     string
// 	GallogIconURL string
// 	Hit           int
// 	Comment       int
// 	Recommend     int
// 	Norecommend   int
// 	IsGoJungNick  bool
// 	IsMobile      bool
// }

// func ArticleParser(url string) (article ArticleReader) {
// 	_getGallogIconURL := func(s *goquery.Selection) string {
// 		iconElement := s.Find(".nick_comm")
// 		for key, value := range gallogIconURLsMap {
// 			if iconElement.HasClass(key) {
// 				return value
// 			}
// 		}
// 		return ""
// 	}

// 	_getGallogURL := func(s *goquery.Selection) string {
// 		href, _ := s.Find(".btn.btn_gall").Attr("href")
// 		return href
// 	}

// 	_getSubject := func(s *goquery.Selection) string {
// 		return strings.TrimSpace(s.Find(".tit_view").Text())
// 	}

// 	_getName := func(s *goquery.Selection) string {
// 		return s.Find(".info_edit > span > span:first-of-type").Text()
// 	}

// 	_getDate := func(s *goquery.Selection) string {
// 		return s.Find(".info_edit > span > span:last-of-type").Text()
// 	}

// 	_getIP := func(s *goquery.Selection) string {
// 		return s.Find(".ip").Text()
// 	}

// 	_getContent := func(s *goquery.Selection) string {
// 		html, _ := s.Find(".view_main").Html()
// 		return html
// 	}

// 	_getHit := func(s *goquery.Selection) int {
// 		hit, _ := strconv.Atoi(s.Find("txt_info > .num:first-of-type").Text())
// 		return hit
// 	}

// 	_getComment := func(s *goquery.Selection) int {
// 		comment, _ := strconv.Atoi(s.Find(".comment_dirc").Text())
// 		return comment
// 	}

// 	_getRecommend := func(s *goquery.Selection) int {
// 		recommend, _ := strconv.Atoi(s.Find("#recomm_btn").Text())
// 		return recommend
// 	}

// 	_getNorecommend := func(s *goquery.Selection) int {
// 		norecommend, _ := strconv.Atoi(s.Find("#nonrecomm_btn").Text())
// 		return norecommend
// 	}

// 	_getIsGoJungNick := func(s *goquery.Selection) bool {
// 		return s.Find(".nick_comm").HasClass("nick_comm")
// 	}

// 	_getIsMobile := func(s *goquery.Selection) bool {
// 		return s.Find(".ico_mobile").HasClass("ico_mobile")
// 	}

// 	s := newDocument(url, nil).Find(".gall_content")

// 	article = ArticleReader{
// 		GallogIconURL: _getGallogIconURL(s),
// 		GallogURL:     _getGallogURL(s),
// 		Subject:       _getSubject(s),
// 		Name:          _getName(s),
// 		Date:          _getDate(s),
// 		Ip:            _getIP(s),
// 		Content:       _getContent(s),
// 		Hit:           _getHit(s),
// 		Comment:       _getComment(s),
// 		Recommend:     _getRecommend(s),
// 		Norecommend:   _getNorecommend(s),
// 		IsGoJungNick:  _getIsGoJungNick(s),
// 		IsMobile:      _getIsMobile(s),
// 	}

// 	return
// }

// type Comment struct {
// 	Name          string
// 	Ip            string
// 	Date          string
// 	GallogURL     string
// 	GallogIconURL string
// 	Content       string
// 	Dccon         string
// 	IsGoJungNick  bool
// 	IsDccon       bool
// }

// func CommentParser(url string) (comment []Comment) {
// 	_getGallogIconURL := func(s *goquery.Selection) string {
// 		iconElement := s.Find(".nick_comm")
// 		for key, value := range gallogIconURLsMap {
// 			if iconElement.HasClass(key) {
// 				return value
// 			}
// 		}
// 		return ""
// 	}

// 	_getGallogURL := func(s *goquery.Selection) string {
// 		href, _ := s.Find(".id").Attr("href")
// 		return href
// 	}

// 	_getName := func(s *goquery.Selection) string {
// 		name := s.Find("#id, .id").Text()
// 		return name[1 : len(name)-1]
// 	}

// 	_getDate := func(s *goquery.Selection) string {
// 		return s.Find(".date").Text()
// 	}

// 	_getIP := func(s *goquery.Selection) string {
// 		return s.Find(".ip").Text()
// 	}

// 	_getContent := func(s *goquery.Selection) string {
// 		return s.Find(".txt").Text()
// 	}

// 	_getDccon := func(s *goquery.Selection) string {
// 		src, _ := s.Find(".written_dccon").Attr("src")
// 		if isDcconURL(src) {
// 			return src
// 		}
// 		return ""
// 	}

// 	_getIsGoJungNick := func(s *goquery.Selection) bool {
// 		return s.Find(".nick_comm").HasClass("nick_comm")
// 	}

// 	_getIsDccon := func(s *goquery.Selection) bool {
// 		return s.Find(".written_dccon").HasClass("written_dccon")
// 	}

// 	_eachComment := func(i int, s *goquery.Selection) {
// 		new := Comment{
// 			GallogIconURL: _getGallogIconURL(s),
// 			GallogURL:     _getGallogURL(s),
// 			Name:          _getName(s),
// 			Date:          _getDate(s),
// 			Ip:            _getIP(s),
// 			Content:       _getContent(s),
// 			Dccon:         _getDccon(s),
// 			IsGoJungNick:  _getIsGoJungNick(s),
// 			IsDccon:       _getIsDccon(s),
// 		}

// 		comment = append(comment, new)
// 	}

// 	newDocument(url, nil).Find(".wrap_list > .list_best > li").Each(_eachComment)

// 	return
// }

// func GetGallogMaximumPage(url string, header map[string]string) (ret int) {
// 	text := newDocument(url, header).Find(".navia > .pg_btn1.pg_btn_prev + .pg_num_area1").Text()
// 	ret, _ = strconv.Atoi(strings.Split(text, "/")[1])
// 	return
// }
