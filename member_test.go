package goinside

import (
	"os"
	"time"
)

func getTestMemberSession() (ms *MemberSession, err error) {
	id := os.Getenv("GOINSIDE_TEST_ID")
	pw := os.Getenv("GOINSIDE_TEST_PW")
	// proxyURL := os.Getenv("GOINSIDE_PROXY_URL")

	ms, err = Login(id, pw)
	if err != nil {
		return
	}

	// proxy, err := url.Parse(proxyURL)
	// if err != nil {
	// 	return
	// }

	// ms.Connection().SetTransport(proxy)
	ms.Connection().SetTimeout(time.Second * 5)
	return
}
