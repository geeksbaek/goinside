package goinside

import (
	"fmt"
	"testing"
)

func TestReport(t *testing.T) {
	URL := "http://gall.dcinside.com/board/view/?id=baseball_new4&no=8946613&page=1"
	s, _ := Guest("ㅇㅇ", "123")
	err := s.Report(URL, "불량 게시물")
	fmt.Println(err)
}
