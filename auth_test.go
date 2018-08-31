package goinside

import (
	"fmt"
	"testing"
)

func TestGetAppID(t *testing.T) {
	gs, err := Guest("ㅇㅇ", "123")
	if err != nil {
		t.Fatal(err)
	}
	appID := GetAppID(gs)
	if appID == "" {
		t.Fatal("could not get app id")
	}
	fmt.Println(appID)
}
