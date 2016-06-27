package goinside

import (
	"net/http"
	"net/url"
	"regexp"
)

var (
	ipRe = regexp.MustCompile(`((?:\d{1,3}\.){3}\d{1,3})`)
)

// Guest 함수는 전달받은 ID, PASSWORD로 생성한 비회원 객체를 반환합니다.
func Guest(id, pw string) *Session {
	return &Session{id: id, pw: pw, nomember: true}
}

// func Login(id, pw string) (*Session, error) {
// 	// return &Session{id: id, pw: pw, nomember: false}
// 	return nil, nil
// }

// SetTransport 함수는 해당 세션이 주어진 프록시를 통해 통신하도록 설정합니다.
// 프록시 주소는 http://84.192.54.48:8080 와 같은 형식으로 전달합니다.
func (s *Session) SetTransport(proxy *url.URL) {
	s.transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	ip := ipRe.FindStringSubmatch(proxy.String())
	if len(ip) == 2 {
		s.ip = ip[1]
	}
}
