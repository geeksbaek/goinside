package goinside

import (
	"net/http"
	"net/url"
)

// Session 구조체는 사용자의 세션을 위해 사용됩니다.
type Session struct {
	id        string
	pw        string
	ip        string
	cookies   []*http.Cookie
	nomember  bool
	transport *http.Transport
}

// Guest 함수는 전달받은 ID, PASSWORD로 생성한 비회원 객체를 반환합니다.
func Guest(id, pw string) *Session {
	return &Session{id: id, pw: pw, nomember: true}
}

// func Login(id, pw string) (*Session, error) {
// 	// return &Session{id: id, pw: pw, nomember: false}
// 	return nil, nil
// }

// SetTransport 함수는 해당 객체가 주어진 프록시를 통해 디시인사이드와 통신하도록 설정합니다. 프록시 주소는 http://84.192.54.48:8080 같은 형식으로 전달합니다.
func (a *Session) SetTransport(URL string) {
	proxyURL, _ := url.Parse(URL)
	a.transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
}
