package goinside

import (
	"fmt"
	"log"
	"testing"
)

// func TestGetAllGall(t *testing.T) {
// 	galls, err := GetAllGall()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, v := range galls {
// 		fmt.Println(v)
// 	}
// }

// func TestGetList(t *testing.T) {
// 	list, err := GetList("http://m.dcinside.com/list.php?id=baseball_new4", 1)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, article := range list.Articles {
// 		fmt.Printf("%#v ", article.AuthorInfo)
// 		fmt.Printf("%#v ", article.Gall)
// 		fmt.Printf("%#v\n", article)
// 	}
// }

func TestGetArticle(t *testing.T) {
	a, err := GetArticle("http://m.dcinside.com/view.php?id=baseball_new4&no=8129734&page=1&exception_mode=recommend")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}
