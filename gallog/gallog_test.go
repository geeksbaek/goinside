package gallog

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	id := os.Getenv("GOINSIDE_TEST_ID")
	pw := os.Getenv("GOINSIDE_TEST_PW")

	_, err := Login(id, pw)
	if err != nil {
		t.Fatal(err)
	}
}
