package goinside_test

import (
	"log"
	"net/url"

	"github.com/geeksbaek/goinside"
)

func ExampleGuest() {
	s, err := goinside.Guest("닉네임", "비밀번호")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSession_SetTransport() {
	proxy, err := url.Parse("http://1.2.3.4:80")
	if err != nil {
		log.Fatal(err)
	}

	s, _ := goinside.Guest("닉네임", "비밀번호")
	s.SetTransport(proxy)

	// ...
}
