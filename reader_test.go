package goinside

import (
	"fmt"
	"testing"
)

func TestGetList(t *testing.T) {
	list, _ := GetList("http://gall.dcinside.com/board/lists/?id=baseball_new4", 16592)
	for _, v := range list.Articles {
		fmt.Println(v.URL, v.Hit)
	}
}

func TestGetArticle(t *testing.T) {
	article, _ := GetArticle("http://gall.dcinside.com/board/view/?id=programming&no=615863&page=1")
	fmt.Println(article.Content)
	fmt.Println(article.Images)
}
