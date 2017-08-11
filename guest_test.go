package goinside

import (
	"net/url"
	"os"
)

func getTestGuestSession() (gs *GuestSession, err error) {
	proxyURL := os.Getenv("GOINSIDE_PROXY_URL")

	gs, err = Guest("ㅇㅇ", "123")

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return
	}

	gs.Connection().SetTransport(proxy)

	return
}
