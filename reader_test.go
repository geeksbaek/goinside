package goinside

import (
	"fmt"
	"testing"
)

func TestReader(t *testing.T) {
	a, _ := FetchArticle("http://gall.dcinside.com/board/view/?id=stock_new1&no=3516851&page=1")

	fmt.Println(a.Detail.Content)
	fmt.Println(a.Detail.ImageURLs)
}
