package goinside

import (
	"fmt"
	"testing"
)

func TestLoginAndWrite(t *testing.T) {
	s, err := Login("goinside", "qweqwe")
	if err != nil {
		panic(err)
	}
	fmt.Println(s.cookies)
}
