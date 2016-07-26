package goinside

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	ipRe = regexp.MustCompile(`((?:\d{1,3}\.){3}\d{1,3})`)
)

// Guest 함수는 전달받은 ID, PASSWORD로 생성한 비회원 세션을 반환합니다.
func Guest(id, pw string) (*Session, error) {
	if id == "" || pw == "" {
		return nil, errors.New("Invaild ID or PW")
	}
	return &Session{id: id, pw: pw, isGuest: true}, nil
}

// Login 함수는 전달받은 ID, PASSWORD로 로그인한 뒤 해당 세션을 반환합니다.
func Login(id, pw string) (*Session, error) {
	s := &Session{}
	loginPageConKey, err := fnLoginGetConKeyFromLoginPage(loginPageURL)
	if err != nil {
		return nil, err
	}
	formForAccessToken := form(map[string]string{
		"token_verify": "login",
		"con_key":      loginPageConKey,
	})
	resp, err := s.post(accessTokenURL, nil, formForAccessToken, nonCharsetContentType)
	if err != nil {
		return nil, err
	}
	cookies := resp.Cookies()
	conKey, err := parseAuthKey(resp)
	if err != nil {
		return nil, err
	}
	formForLogin := form(map[string]string{
		"user_id": id,
		"user_pw": pw,
		"id_chk":  "on",
		"con_key": conKey,
	})
	resp, err = s.post(mobileLoginURL, cookies, formForLogin, nonCharsetContentType)
	if err != nil {
		return nil, err
	}
	if len(resp.Cookies()) == 0 {
		return nil, errors.New("No Response. Login Fail")
	}
	s.id, s.pw, s.cookies, s.isGuest = id, pw, resp.Cookies(), false
	return s, nil
}

// Logout 함수는 해당 세션을 종료합니다.
func (s *Session) Logout() error {
	if s.isGuest {
		return nil
	}
	_, err := s.get(logoutURL)
	s = &Session{}
	return err
}

// SetTransport 함수는 해당 세션이 주어진 프록시를 통해 통신하도록 설정합니다.
// 프록시 주소는 http://84.192.54.48:8080 와 같은 형식으로 전달합니다.
func (s *Session) SetTransport(proxy *url.URL) {
	s.proxy = http.ProxyURL(proxy)
	ip := ipRe.FindStringSubmatch(proxy.String())
	if len(ip) == 2 {
		s.ip = ip[1]
	}
}

// SetTimeout 함수는 해당 세션의 통신에 timeout 값을 설정합니다.
func (s *Session) SetTimeout(time time.Duration) {
	s.timeout = time
}
