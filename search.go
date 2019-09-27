package goinside

import "fmt"

func search(id, keyword, searchType string, page int) (l *List, err error) {
	gall := &Gall{ID: id, URL: gallURL(id)}
	formMap := map[string]string{
		"app_id": RandomGuest().getAppID(),
		"id":     id,
		"s_type": searchType,
		"page":   fmt.Sprint(page),
		"serVal": keyword,
		// "ser_pos": "",
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
			URL:                articleURL(id, a.Number),
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

// Search 함수는 해당 키워드와 제목이나 내용, 혹은 글쓴이가 일치하는 글을 검색합니다.
func Search(id, keyword string) (l *List, err error) {
	return search(id, keyword, "all", 1)
}

// SearchBySubject 함수는 해당 키워드와 제목이 일치하는 글을 검색합니다.
func SearchBySubject(id, keyword string) (l *List, err error) {
	return search(id, keyword, "subject", 1)
}

// SearchByContent 함수는 해당 키워드와 내용이 일치하는 글을 검색합니다.
func SearchByContent(id, keyword string) (l *List, err error) {
	return search(id, keyword, "memo", 1)
}

// SearchBySubjectAndContent 함수는 해당 키워드와 제목이나 내용이 일치하는 글을 검색합니다.
func SearchBySubjectAndContent(id, keyword string) (l *List, err error) {
	return search(id, keyword, "subject_m", 1)
}

// SearchByAuthor 함수는 해당 키워드와 글쓴이가 일치하는 글을 검색합니다.
func SearchByAuthor(id, keyword string) (l *List, err error) {
	return search(id, keyword, "name", 1)
}
