package goinside

import (
	"fmt"
	"log"
	"testing"
)

func TestGetAllGall(t *testing.T) {
	galls, err := GetAllGall()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range galls {
		fmt.Println(v)
	}
}

func TestGetList(t *testing.T) {
	list, err := GetList("http://m.dcinside.com/list.php?id=baseball_new4", 1)
	if err != nil {
		log.Fatal(err)
	}
	for _, article := range list.Articles {
		fmt.Printf("%#v ", article.AuthorInfo)
		fmt.Printf("%#v ", article.Gall)
		fmt.Printf("%#v\n", article)
	}
}