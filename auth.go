package goinside

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
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

func fetchAppID(s session) (valueToken, appID string, err error) {
	valueToken, err = generateValueToken()
	if err != nil {
		return "", "", err
	}
	r, ct := multipartForm(map[string]string{
		"value_token": valueToken,
		"signature":   "ReOo4u96nnv8Njd7707KpYiIVYQ3FlcKHDJE046Pg6s=",
		"pkg":         "com.dcinside.app",
		// "vCode":       "?",
		// "vName":       "?",
	})
	res, err := appKeyVerificationAPI.post(s, r, ct)
	if err != nil {
		return "", "", nil
	}
	defer res.Body.Close()
	appIDResponse := AppIDResponse{}
	if err := json.NewDecoder(res.Body).Decode(&appIDResponse); err != nil {
		return "", "", nil
	}
	if len(appIDResponse) == 0 || appIDResponse[0].Result == false {
		return "", "", nil
	}
	appID = appIDResponse[0].AppID

	time.Sleep(time.Second * 5)
	return
}

func generateValueToken() (string, error) {
	resp, err := appCheckAPI.getWithoutHash()
	if err != nil {
		return "", fmt.Errorf("appCheckAPI.getWithoutHash fail: %v", err)
	}
	body := []map[string]interface{}{}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("appCheckAPI json decode fail: %v", err)
	}
	if len(body) != 1 {
		return "", errors.New("unknown app check response")
	}
	appKey := fmt.Sprintf("dcArdchk_%s", body[0]["date"])
	hashedAppKey := sha256.Sum256([]byte(appKey))
	return fmt.Sprintf("%x", hashedAppKey), nil
}
