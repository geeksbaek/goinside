package goinside_test

import (
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/geeksbaek/goinside"
)

func Example() {
	// 세션 생성하기
	s, err := goinside.Guest("id", "pw")
	if err != nil {
		log.Fatal(err)
	}

	var (
		gallID  = "gallID" // ex) baseball_new4
		subject = "글 제목"
		content = "글 내용"
		images  = []string{ // 이미지 첨부파일
			"example.jpg",
			"example.gif",
		}
	)

	// 글 작성하기
	article, err := s.WriteArticle(gallID, subject, content, images...)
	if err != nil {
		log.Fatal(err)
	}

	// 댓글 작성하기
	comment, err := s.WriteComment(article, "댓글 내용")
	if err != nil {
		log.Fatal(err)
	}

	// 글 추천하기
	err = s.ThumbsUp(article)
	if err != nil {
		log.Fatal(err)
	}

	// 글 비추천하기
	err = s.ThumbsDown(article)
	if err != nil {
		log.Fatal(err)
	}

	// 댓글 삭제하기
	err = s.Delete(comment)
	if err != nil {
		log.Fatal(err)
	}

	// 글 삭제하기
	err = s.Delete(article)
	if err != nil {
		log.Fatal(err)
	}

	// Option: 프록시 설정하기
	proxy, err := url.Parse("http://1.2.3.4:80")
	if err != nil {
		log.Fatal(err)
	}
	s.SetTransport(proxy)

	// Option: timeout 설정하기
	s.SetTimeout(time.Second * 10)

	// 모든 갤러리 목록 가져오기
	galls, err := goinside.GetAllGall()
	if err != nil {
		log.Fatal(err)
	}
	for _, gall := range galls {
		_ = gall // unused
	}

	// 특정 갤러리 글 목록 가져오기
	page := 1
	list, err := goinside.GetList(galls[rand.Intn(len(galls))].URL, page)
	if err != nil {
		log.Fatal(err)
	}
	for _, article := range list.Articles {
		_ = article // unused
	}

	// 특정 글 정보 가져오기
	article, err = goinside.GetArticle(list.Articles[rand.Intn(len(list.Articles))].URL)
	if err != nil {
		log.Fatal(err)
	}
	for _, comment := range article.Comments {
		_ = comment // unused
	}
}
