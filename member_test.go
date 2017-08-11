package goinside

import (
	"net/url"
	"os"
)

func getTestMemberSession() (ms *MemberSession, err error) {
	id := os.Getenv("GOINSIDE_TEST_ID")
	pw := os.Getenv("GOINSIDE_TEST_PW")
	proxyURL := os.Getenv("GOINSIDE_PROXY_URL")

	ms, err = Login(id, pw)

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return
	}

	ms.Connection().SetTransport(proxy)
	return
}
