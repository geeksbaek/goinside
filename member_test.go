package goinside

import (
	"encoding/json"
	"io/ioutil"
)

func getTestMemberSession() (ms *MemberSession, err error) {
	authFile, err := ioutil.ReadFile("auth.json")
	if err != nil {
		return
	}
	var auth struct {
		ID string
		PW string
	}
	err = json.Unmarshal(authFile, &auth)
	if err != nil {
		return
	}
	ms, err = Login(auth.ID, auth.PW)
	return
}
