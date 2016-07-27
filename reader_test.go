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
