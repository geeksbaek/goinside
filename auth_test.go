package goinside

import (
	"testing"
)

func TestGetAppID(t *testing.T) {
	gs, err := Guest("ㅇㅇ", "123")
	if err != nil {
		t.Fatal(err)
	}
	if appID := gs.getAppID(); appID == "" {
		t.Fatal("could not get app id")
	} else {
		t.Log(appID)
	}
}
