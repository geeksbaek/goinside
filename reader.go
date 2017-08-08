package goinside

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func fetchSomething(formMap map[string]string, api dcinsideAPI, data interface{}) (err error) {
	resp, err := api.get(formMap)
	if err != nil {
		return
	}
	valid := make(jsonValidation, 1)
	if err = responseUnmarshal(resp, data, &valid); err != nil {
		return
	}
	if err = checkJSONResult(&valid); err != nil {
		return
	}
	return
}

type jsonGallery []struct {
	Category    string `json:"category"`
	ID          string `json:"name"`
	Name        string `json:"ko_name"`
	Number      string `json:"no"`
	Depth       string `json:"depth"`
	CanWrite    bool   `json:"no_write"`
	IsAdultOnly bool   `json:"is_adult"`
}

// FetchAllMajorGallery 함수는 모든 일반 갤러리의 정보를 가져옵니다.
func FetchAllMajorGallery() (mg []*MajorGallery, err error) {
	resp, err := majorGalleryListAPI.get(nil)
	if err != nil {
		return
	}
	jsonResp := jsonGallery{}
	if err = responseUnmarshal(resp, &jsonResp); err != nil {
		return
	}
	mg = make([]*MajorGallery, len(jsonResp))
	for i, v := range jsonResp {
		mg[i] = &MajorGallery{
			ID:       v.ID,
			Name:     v.Name,
			Number:   v.Number,
			CanWrite: !v.CanWrite,
		}
	}
	return
}

type jsonMonirGallery []struct {
	Category    string `json:"category"`
	ID          string `json:"name"`
	Name        string `json:"ko_name"`
	Number      string `json:"no"`
	Depth       string `json:"depth"`
	CanWrite    bool   `json:"no_write"`
	IsAdultOnly bool   `json:"is_adult"`
	Manager     string `json:"manager"`
	SubManagers string `json:"submanager"`
}

// FetchAllMinorGallery 함수는 모든 마이너 갤러리의 정보를 가져옵니다.
func FetchAllMinorGallery() (mg []*MinorGallery, err error) {
	resp, err := minorGalleryListAPI.get(nil)
	if err != nil {
		return
	}
	jsonResp := jsonMonirGallery{}
	if err = responseUnmarshal(resp, &jsonResp); err != nil {
		return
	}
	mg = make([]*MinorGallery, len(jsonResp))
	for i, v := range jsonResp {
		mg[i] = &MinorGallery{
			ID:          v.ID,
			Name:        v.Name,
			Number:      v.Number,
			CanWrite:    !v.CanWrite,
			Manager:     v.Manager,
			SubManagers: strings.Split(v.SubManagers, ","),
		}
	}
	return
}

type jsonList []struct {
	GallInfo []struct {
		CategoryName string `json:"category_name"`
		FileCount    string `json:"file_cnt"`
		FileSize     string `json:"file_size"`
	} `json:"gall_info"`
	GallList []struct {
		Subject      string `json:"subject"`
		Name         string `json:"name"`
		Level        string `json:"level"`
		ImageIcon    string `json:"img_icon"`
		WinnertaIcon string `json:"winnerta_icon"`
		ThumbsUp     string `json:"recommend"`
		ThumbsUpIcon string `json:"recommend_icon"`
		IsBest       string `json:"best_chk"`
		Hit          string `json:"hit"`
		UserID       string `json:"user_id"`
		MemberIcon   string `json:"member_icon"`
		IP           string `json:"ip"`
		TotalComment string `json:"total_comment"`
		TotalVoice   string `json:"total_voice"`
		Number       string `json:"no"`
		Date         string `json:"date_time"`
	} `json:"gall_list"`
}

// FetchList 함수는 해당 갤러리의 해당 페이지에 있는 글의 목록을 가져옵니다.
func FetchList(URL string, page int) (l *List, err error) {
	return fetchList(URL, page, false)
}

// FetchBestList 함수는 해당 갤러리의 해당 페이지에 있는 개념글의 목록을 가져옵니다.
func FetchBestList(URL string, page int) (l *List, err error) {
	return fetchList(URL, page, true)
}

func fetchList(URL string, page int, fetchBestPage bool) (l *List, err error) {
	gallID := gallID(URL)
	gall := &Gall{ID: gallID, URL: URL}
	formMap := map[string]string{
		"app_id": AppID,
		"id":     gallID,
		"page":   fmt.Sprint(page),
	}
	if fetchBestPage {
		formMap["recommend"] = "1"
	}
	respJSON := make(jsonList, 1)
	if err = fetchSomething(formMap, readListAPI, &respJSON); err != nil {
		return
	}
	r := respJSON[0]
	l = &List{
		Info: &ListInfo{
			CategoryName: r.GallInfo[0].CategoryName,
			FileCount:    r.GallInfo[0].FileCount,
			FileSize:     r.GallInfo[0].FileSize,
			Gall:         gall,
		},
		Items: []*ListItem{},
	}
	for _, a := range r.GallList {
		item := &ListItem{
			Gall:               gall,
			URL:                articleURL(gallID, a.Number),
			Subject:            a.Subject,
			Name:               a.Name,
			Level:              Level(a.Level),
			HasImage:           a.ImageIcon == "Y",
			ArticleType:        articleType(a.ImageIcon, a.IsBest),
			ThumbsUp:           mustAtoi(a.ThumbsUp),
			IsBest:             a.IsBest == "Y",
			Hit:                mustAtoi(a.Hit),
			GallogID:           a.UserID,
			GallogURL:          gallogURL(a.UserID),
			IP:                 a.IP,
			CommentLength:      mustAtoi(a.TotalComment),
			VoiceCommentLength: mustAtoi(a.TotalVoice),
			Number:             a.Number,
			Date:               dateFormatter(a.Date),
		}
		l.Items = append(l.Items, item)
	}
	return
}

// Fetch 메소드는 해당 글의 세부 정보(본문, 이미지 주소, 댓글)를 가져옵니다.
func (i *ListItem) Fetch() (*Article, error) {
	return FetchArticle(i.URL)
}

// FetchImageURLs 메소드는 해당 글의 이미지 주소의 슬라이스만을 가져옵니다.
func (i *ListItem) FetchImageURLs() (imageURLs []ImageURLType, err error) {
	formMap := map[string]string{
		"app_id": AppID,
		"id":     i.Gall.ID,
		"no":     fmt.Sprint(i.Number),
	}
	images := make(jsonArticleImages, 1)
	err = fetchSomething(formMap, readArticleImageAPI, &images)
	if err != nil {
		return
	}
	imageURLs = func() (ret []ImageURLType) {
		for _, v := range images {
			ret = append(ret, ImageURLType(v.Image))
		}
		return
	}()
	return
}

// Fetch 메소드는 해당 이미지 주소를 참조하여 이미지의 []byte를 반환합니다.
func (i ImageURLType) Fetch() (data []byte, err error) {
	resp, err := doImage(i)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

type jsonArticleContent []struct {
	Memo       string `json:"memo"`
	ThumbsUp   string `json:"recommend"`
	ThumbsDown string `json:"nonrecommend"`
}

type jsonArticleDetail []struct {
	Subject            string `json:"subject"`
	Number             string `json:"no"`
	Name               string `json:"name"`
	MemberIcon         string `json:"member_icon"`
	IP                 string `json:"ip"`
	TotalComment       string `json:"total_comment"`
	HasImage           string `json:"img_chk"`
	IsBest             string `json:"recommend_chk"`
	IsWinnerta         string `json:"winnerta_chk"`
	Page               string `json:"page"`
	Hit                string `json:"hit"`
	WriteType          string `json:"write_type"`
	UserID             string `json:"user_id"`
	PrevArticleNumber  string `json:"prev_link"`
	PrevArticleSubject string `json:"prev_subject"`
	NextArticleNumber  string `json:"next_link"`
	NextArticleSubject string `json:"next_subject"`
	// _                  string `json:"best_chk"` // ?
	Date string `json:"date_time"`
}

type jsonArticleImages []struct {
	Image string `json:"img"`
	// ImageClone string `json:"img_clone"`
}

// FetchArticle 함수는 해당 글의 정보를 가져옵니다.
func FetchArticle(URL string) (a *Article, err error) {
	gallID := gallID(URL)
	gall := &Gall{ID: gallID, URL: gallURL(gallID)}
	formMap := map[string]string{
		"app_id": AppID,
		"id":     gallID,
		"no":     articleNumber(URL),
	}

	content := make(jsonArticleContent, 1)
	detail := make(jsonArticleDetail, 1)
	images := make(jsonArticleImages, 1)

	ch := func() <-chan error {
		ch := make(chan error)
		go func() {
			ch <- fetchSomething(formMap, readArticleAPI, &content)
		}()
		go func() {
			ch <- fetchSomething(formMap, readArticleDetailAPI, &detail)
		}()
		go func() {
			fetchSomething(formMap, readArticleImageAPI, &images)
			ch <- nil
		}()
		return ch
	}()

	for i := 0; i < 3; i++ {
		if err := <-ch; err != nil {
			return nil, err
		}
	}

	c := content[0]
	d := detail[0]

	article := &Article{
		Gall:          gall,
		URL:           articleURL(gallID, d.Number),
		Subject:       d.Subject,
		Content:       c.Memo,
		ThumbsUp:      mustAtoi(c.ThumbsUp),
		ThumbsDown:    mustAtoi(c.ThumbsDown),
		Name:          d.Name,
		Number:        d.Number,
		Level:         MemberType(mustAtoi(d.MemberIcon)).Level(),
		IP:            d.IP,
		CommentLength: mustAtoi(d.TotalComment),
		HasImage:      d.HasImage == "Y",
		Hit:           mustAtoi(d.Hit),
		ArticleType:   articleType(d.HasImage, d.IsBest),
		GallogID:      d.UserID,
		GallogURL:     gallogURL(d.UserID),
		IsBest:        d.IsBest == "Y",
		ImageURLs: func() (ret []ImageURLType) {
			for _, v := range images {
				ret = append(ret, ImageURLType(v.Image))
			}
			return
		}(),
		Comments: []*Comment{},
		Date:     dateFormatter(d.Date),
	}
	if article.CommentLength > 0 {
		article.Comments, err = fetchComment(URL, article)
		if err != nil {
			return
		}
	}
	return article, nil
}

type jsonComment []struct {
	CommentCount string `json:"total_comment"`
	TotalPage    string `json:"total_page"`
	NowPage      string `json:"re_page"`
	Comments     []struct {
		Name       string `json:"name"`
		UserID     string `json:"user_id"`
		UserNO     string `json:"user_no"`
		Level      string `json:"level"`
		MemberIcon string `json:"member_icon"`
		Content    string `json:"comment_memo"`
		IP         string `json:"ipData"`
		Voice      string `json:"voice"`
		DCcon      string `json:"dccon"`
		Number     string `json:"comment_no"`
		Date       string `json:"date_time"`
	} `json:"comment_list"`
}

func fetchComment(URL string, parents *Article) (cs []*Comment, err error) {
	gallID := gallID(URL)
	gallURL := gallURL(gallID)
	gall := &Gall{ID: gallID, URL: gallURL}
	cs = []*Comment{}
	for commentPage := 1; ; commentPage++ {
		formMap := map[string]string{
			"app_id":  AppID,
			"id":      gallID,
			"no":      parents.Number,
			"re_page": fmt.Sprint(commentPage),
		}
		respJSON := make(jsonComment, 1)
		if err = fetchSomething(formMap, readCommentAPI, &respJSON); err != nil {
			return
		}
		r := respJSON[0]
		for _, c := range r.Comments {
			comment := &Comment{
				Gall:      gall,
				Parents:   parents,
				Type:      commentType(c.DCcon, c.Voice),
				Name:      c.Name,
				GallogID:  c.UserID,
				GallogURL: gallogURL(c.UserID),
				Level:     Level(c.Level),
				IP:        c.IP,
				Number:    c.Number,
				Date:      dateFormatter(c.Date),
			}
			comment.Content, comment.HTML = func() (string, string) {
				switch comment.Type {
				case TextCommentType:
					return c.Content, c.Content
				case DCconCommentType:
					return c.DCcon, toImageElement(c.DCcon)
				case VoiceCommentType:
					return c.Voice, toAudioElement(c.Voice)
				}
				return c.Content, c.Content
			}()
			cs = append(cs, comment)
		}
		if mustAtoi(r.NowPage) >= mustAtoi(r.TotalPage) {
			break
		}
	}
	return
}
