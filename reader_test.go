package goinside

import (
	"fmt"
	"testing"
)

func TestReader(t *testing.T) {
	a, _ := FetchArticle("http://gall.dcinside.com/board/view/?id=stock_new1&no=3516942&page=1")

	fmt.Println(a.Detail.Content)
	fmt.Println(a.Detail.ImageURLs)
}

func TestFetchGallerys(t *testing.T) {

}

func TestComments(t *testing.T) {
	a, _ := FetchArticle("http://gall.dcinside.com/board/view/?id=baseball_new4&no=9509049&page=1&exception_mode=recommend")
	fmt.Println(a.CommentCount)
	fmt.Println(len(a.Detail.Comments))
	for _, v := range a.Detail.Comments {
		fmt.Printf("%#v %#v\n", v.Author, v.Author.Detail)
	}
}
