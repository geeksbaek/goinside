package goinside_test

import (
	"log"
	"time"

	"github.com/geeksbaek/goinside"
)

// ExampleGuest 함수는 비회원 세션을 사용하는 예제입니다.
// 먼저 programming이라는 ID를 가진 갤러리에 글을 하나 작성합니다.
// 그리고 바로 해당 글에 댓글을 하나 작성하고요. 추천과 비추천도 누릅니다.
// 그리고 5초 뒤에 댓글과 글을 삭제합니다.
func ExampleGuest() {
	s := goinside.Guest("닉네임", "비밀번호")

	gall := "programming"
	subject := "글 제목"
	content := "글 내용"
	images := []string{"첨부파일1.jpg", "첨부파일2.gif"}

	// 글 작성
	article, err := s.WriteArticle(gall, subject, content, images...)
	if err != nil {
		log.Panic(err)
	}

	// 댓글 작성
	comment, err := s.WriteComment(article, "댓글 내용")
	if err != nil {
		log.Panic(err)
	}

	s.ThumbsUp(article)   // 추천
	s.ThumbsDown(article) // 비추천

	time.Sleep(time.Second * 5) // 5초 뒤

	s.Delete(comment) // 댓글 삭제
	s.Delete(article) // 글 삭제
}
