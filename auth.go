package goinside

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// AppIDResponse 구조체는 appKeyVerificationAPI 요청의 응답을 정의합니다.
type AppIDResponse []struct {
	Result bool   `json:"result"`
	AppID  string `json:"app_id"`
}

// GetAppID 함수는 디시인사이드 서버로부터 App ID를 가져옵니다.
func GetAppID(s session) (appID string) {
	now := time.Now()
	appKey := fmt.Sprintf("dcArdchk_%04d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour())
	hashedAppKey := sha256.Sum256([]byte(appKey))

	res, err := appKeyVerificationAPI.post(s, makeForm(map[string]string{
		"value_token": fmt.Sprintf("%x", hashedAppKey),
		"signature":   "ReOo4u96nnv8Njd7707KpYiIVYQ3FlcKHDJE046Pg6s=",
		// "pkg":         "com.dcinside.app",
		// "vCode":       "?",
		// "vName":       "?",
	}), nonCharsetContentType)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	appIDResponse := AppIDResponse{}
	if err := json.NewDecoder(res.Body).Decode(&appIDResponse); err != nil {
		return ""
	}
	if len(appIDResponse) == 0 || appIDResponse[0].Result == false {
		return ""
	}
	appID = appIDResponse[0].AppID
	return
}
