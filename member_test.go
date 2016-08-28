package goinside

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestLogin(t *testing.T) {
	authFile, err := ioutil.ReadFile("auth.json")
	if err != nil {
		t.Error(err)
	}
	var auth struct {
		ID string
		PW string
	}
	err = json.Unmarshal(authFile, &auth)
	if err != nil {
		t.Error(err)
	}
	_, err = Login(auth.ID, auth.PW)
	if err != nil {
		t.Error(err)
	}
}
