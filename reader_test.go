package goinside_test

import "github.com/geeksbaek/goinside"

func ExampleGetAllGall() {
	galls, _ := goinside.GetAllGall()
}

// 프로그래밍 갤러리의 첫 번째 페이지의 글들을 가져온다.
// URL은 데스크톱 버전이든 모바일 버전이든 상관없다.
func ExampleGetList() {
	list, _ := goinside.GetList("http://gall.dcinside.com/board/lists/?id=programming", 1)
    for _, article := range list.Articles {
        // ...
    }
}

// 해당 URL의 글 정보를 가져온다. 마찬가지로 URL은 데스크톱 버전이든
// 모바일 버전이든 상관없다. 해당 글에 달려있는 댓글들까지 모두 가져온다.
func ExampleGetArticle() {
    article, _ := goinside.GetArticle("http://gall.dcinside.com/board/view/?id=programming&no=603814&page=1&exception_mode=recommend")
    for _, comment := range article.Comments {
        // ...
    }
}
