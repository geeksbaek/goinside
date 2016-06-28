package goinside_test

import (
	"log"

	"github.com/geeksbaek/goinside"
)

func ExampleSession_WriteArticle() {
	s, _ := goinside.Guest("닉네임", "비밀번호")

	gall := "programming"
	subject := "글 제목"
	content := "글 내용"
	images := []string{"첨부파일1.jpg", "첨부파일2.gif"}

	// 글 작성
	article, err := s.WriteArticle(gall, subject, content, images...)
	if err != nil {
		log.Panic(err)
	}
	s.Delete(article) // 글 삭제
}

func ExampleSession_WriteComment() {
	s, _ := goinside.Guest("닉네임", "비밀번호")
	comment, err := s.WriteComment(article, "댓글 내용")
	if err != nil {
		log.Panic(err)
	}
	s.Delete(comment) // 댓글 삭제
}

// Delete는 삭제 가능한 인자(글, 댓글)들을 가변 인자로 받아
// 동시적으로 삭제하는 함수이다.
// func ExampleSession_Delete() {
// 	s.Delete(article)
//     s.Delete(articleSlice...)
//     s.Delete(article1, comment1, article2, comment2)
// }