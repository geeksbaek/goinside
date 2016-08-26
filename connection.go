package goinside

import (
	"net/http"
	"net/url"
	"time"
)

// Connection 구조체는 HTTP 통신에 필요한 정보를 나타냅니다.
type Connection struct {
	proxy   func(*http.Request) (*url.URL, error)
	timeout time.Duration
}

// SetTransport 함수는 해당 세션이 주어진 프록시를 통해 통신하도록 설정합니다.
// 프록시 주소는 http://84.192.54.48:8080 와 같은 형식으로 전달합니다.
func (c *Connection) SetTransport(proxy *url.URL) {
	c.proxy = http.ProxyURL(proxy)
}

// SetTimeout 함수는 해당 세션의 통신에 timeout 값을 설정합니다.
func (c *Connection) SetTimeout(time time.Duration) {
	c.timeout = time
}
