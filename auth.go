package goinside

import (
	"net/http"
	"net/url"
)

type Auth struct {
	id        string
	pw        string
	ip        string
	cookies   []*http.Cookie
	nomember  bool
	transport *http.Transport
}

// GetNomemberAuth 함수는 전달받은 ID, PASSWORD로 생성한 비회원 객체를 반환합니다.
func GetNomemberAuth(id, pw string) *Auth {
	return &Auth{id: id, pw: pw, nomember: true}
}

// func Login(id, pw string) (*Auth, error) {
// 	// return &Auth{id: id, pw: pw, nomember: false}
// 	return nil, nil
// }

// SetTransport 함수는 해당 객체가 주어진 프록시를 통해 디시인사이드와 통신하도록 설정합니다.
func (a *Auth) SetTransport(URL *url.URL) {
	a.transport = &http.Transport{Proxy: http.ProxyURL(URL)}
}
