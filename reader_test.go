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
	article, err := GetArticle("http://m.dcinside.com/view.php?id=baseball_new4&no=8109733&page=1&recommend=1")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%#v\n", article.AuthorInfo)
	// fmt.Printf("%#v\n", article.Gall)
	// fmt.Printf("%#v\n", article.Number)
	for _, comment := range article.Comments {
		fmt.Println(
			// comment.AuthorInfo.Name,
			// comment.AuthorInfo.IP,
			// comment.AuthorInfo.IsGuest,
			// comment.AuthorInfo.GallogID,
			// comment.AuthorInfo.GallogURL,
			// comment.AuthorInfo.GallogIcon,
			// comment.Gall.URL,
			// comment.Gall.ID,
			// comment.Gall.Name,
			comment.Number,
			comment.Content,
			comment.Date,
		)
	}
}
