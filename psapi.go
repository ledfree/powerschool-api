package psapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

type ApiConfig struct {
	PsUrl        string
	ClientId     string
	ClientSecret string
}

func validateConfig() error {

	return nil
}

func GetAccessToken(psurl, clientId, clientSecret string) (status int, token string, err error) {
	status = 400 // bad request

	str64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret)))
	u, err := url.Parse(psurl)
	if err != nil {
		return
	}

	resp, err := http.Post(u.Hostname(), "applicaion/json", nil)
	if err != nil {
		return
	}
	status = resp.StatusCode
	token = resp.Body

	token = u.Hostname()

	return
}
