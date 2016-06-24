/*
goinside는 비공식 dcinside API 입니다.

goinside는 기본적인 글, 댓글 작성, 삭제 및 추천, 비추천 기능을 제공하며,
갤러리나 글의 정보를 가져오는 기능을 제공합니다.

또한 프록시를 설정하여 익명으로 dcinside와 통신하는 기능을 제공합니다.

다음은 programming 갤러리에 비회원 세션으로 이미지가 2개인 글을 작성하고
해당 게시물에 댓글을 하나 작성한 뒤 추천과 비추천을 보내고, 
5초 뒤에 작성한 댓글과 글을 삭제하는 예제입니다.

 ss := goinside.Guest("이름", "비밀번호")

 gall := "programming"
 subject := "글 제목"
 content := "글 내용"
 images := []string{"image1.jpg", "image2.gif"} // 이미지 첨부 파일

 // 글 작성
 article, err := ss.NewArticle(gall, subject, content, images...).Write()
 if err != nil {
    log.Fatal(err)
 }

 // 댓글 작성
 comment, err := ss.NewComment(article, "댓글 내용").Write()
 if err != nil {
    log.Fatal(err)
 }

 ss.ThumbsUp(article)   // 추천
 ss.ThumbsDown(article) // 비추천

 time.Sleep(time.Second * 5)

 ss.DeleteComment(comment) // 댓글 삭제
 ss.DeleteArticle(article) // 글 삭제
*/
package goinside
