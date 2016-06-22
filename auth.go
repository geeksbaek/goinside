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

func LoginNomember(id, pw string) *Auth {
	return &Auth{id: id, pw: pw, nomember: true}
}

func Login(id, pw string) (*Auth, error) {
	// return &Auth{id: id, pw: pw, nomember: false}
	return nil, nil
}

func (a *Auth) SetTransport(URL *url.URL) {
	a.transport = &http.Transport{Proxy: http.ProxyURL(URL)}
}
