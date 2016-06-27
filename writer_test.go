package goinside

// func TestWriter(t *testing.T) {
// 	ss := Guest("이름", "비밀번호")
// 	proxy, _ := url.Parse("http://104.200.129.174:80")
// 	ss.SetTransport(proxy)

// 	gall := "china"
// 	subject := "글 제목"
// 	content := "글 내용"
// 	// images := []string{"image1.jpg", "image2.gif"} // 이미지 첨부 파일

// 	// 글 작성
// 	article, err := ss.WriteArticle(gall, subject, content)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	browser.OpenURL(article.URL)

// 	// 댓글 작성
// 	comment, err := ss.WriteComment(article, "댓글 내용")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

//     ss.PrefetchDetail(article)

// 	fmt.Println(ss.ThumbsUp(article))
// 	fmt.Println(ss.ThumbsDown(article))

// 	time.Sleep(time.Second * 10)

// 	ss.Delete(comment) // 댓글 삭제
// 	ss.Delete(article) // 글 삭제
// }
