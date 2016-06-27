package goinside_test

import (
	"log"
	"net/url"
	"time"

	"github.com/geeksbaek/goinside"
)

// ExampleGuest 함수는 비회원 세션을 사용하는 예제입니다.
// 먼저 programming이라는 ID를 가진 갤러리에 2개의 이미지를 첨부한 글을 하나 작성합니다.
// WriteArticle의 마지막 매개 변수인 images는 가변 인자이기 때문에
// 슬라이스를 전달하기 위해서 ...을 사용해야 합니다.
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

// ExampleSetTransport 함수는 SetTransport()를 이용하여 세션을 프록시로
// 통신하게 하는 방법을 보여줍니다. 이제 해당 세션은 모든 통신에서
// http://1.2.3.4:80 라는 프록시를 통해 통신합니다.
func ExampleSetTransport() {
	proxy, err := url.Parse("http://1.2.3.4:80")
	if err != nil {
		log.Fatal(err)
	}

	s := goinside.Guest("닉네임", "비밀번호")
	s.SetTransport(proxy)

	// ...
}
