package goinside

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// App 구조체는 AppKey와 AppId를 구성합니다.
type App struct {
	Token string
	ID    string
}

// AppIDResponse 구조체는 appKeyVerificationAPI 요청의 응답을 정의합니다.
type AppIDResponse []struct {
	Result bool   `json:"result"`
	AppID  string `json:"app_id"`
}

func fetchAppID(s session) (valueToken, appID string) {
	valueToken = generateValueToken()
	r, ct := multipartForm(map[string]string{
		"value_token": valueToken,
		"signature":   "ReOo4u96nnv8Njd7707KpYiIVYQ3FlcKHDJE046Pg6s=",
		"pkg":         "com.dcinside.app",
		// "vCode":       "?",
		// "vName":       "?",
	})
	res, err := appKeyVerificationAPI.post(s, r, ct)
	if err != nil {
		return "", ""
	}
	defer res.Body.Close()
	appIDResponse := AppIDResponse{}
	if err := json.NewDecoder(res.Body).Decode(&appIDResponse); err != nil {
		return "", ""
	}
	if len(appIDResponse) == 0 || appIDResponse[0].Result == false {
		return "", ""
	}
	appID = appIDResponse[0].AppID

	time.Sleep(time.Second * 5)
	return
}

func generateValueToken() string {
	now := time.Now()
	appKey := fmt.Sprintf("dcArdchk_%04d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour())
	hashedAppKey := sha256.Sum256([]byte(appKey))
	return fmt.Sprintf("%x", hashedAppKey)
}
