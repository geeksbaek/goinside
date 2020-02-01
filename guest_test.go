package goinside

import (
	"net/url"
	"os"
	"time"
)

func getTestGuestSession() (gs *GuestSession, err error) {
	proxyURL := os.Getenv("GOINSIDE_PROXY_URL")

	gs, err = Guest("ㅇㅇ", "123")
	if err != nil {
		return
	}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		gs.Connection().SetTransport(proxy)
	}

	gs.Connection().SetTimeout(time.Second * 5)
	return
}
