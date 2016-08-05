package goinside

import (
	"fmt"
	"testing"
)

func TestGetList(t *testing.T) {
	list, _ := GetList("http://gall.dcinside.com/board/lists/?id=baseball_new4", 1)
	for _, v := range list.Articles {
		article, _ := GetArticle(v.URL)
		fmt.Println(article.Content)
		fmt.Println(article.Images)
		fmt.Println("-------------------------------------------------")
	}
}

func TestGetArticle(t *testing.T) {
	article, _ := GetArticle("http://gall.dcinside.com/board/view/?id=game_classic&no=10501388&page=1")
	fmt.Println(article.Content)
	fmt.Println(article.HasImage)
	fmt.Println(article.Images)
}
