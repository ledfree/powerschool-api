package psapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type tokenResponse struct {
	Access_Token string `json:"access_token"`
	Token_Type   string `json:"token_type"`
	Expires_In   string `json:"expires_in"`
}

type ApiConfig struct {
	PsUrl        string
	ClientId     string
	ClientSecret string
}

func (a ApiConfig) toBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.ClientId, a.ClientSecret)))
}

func (a ApiConfig) GetAccessToken() (status int, token string, err error) {
	status = 400 // bad request

	u, err := url.Parse(a.PsUrl)
	if err != nil {
		err = fmt.Errorf("invalid URL %s : %s", a.PsUrl, err)
		return
	}

	urlToken := fmt.Sprintf("https://%s/oauth/access_token/", u.Hostname())
	body := []byte(`{"body": "grant_type=client_credentials"}`)
	req, err := http.NewRequest(http.MethodPost, urlToken, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", a.toBase64()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	status = res.StatusCode

	if res.StatusCode != http.StatusOK {
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var t tokenResponse
	err = json.Unmarshal(resBody, &t)
	if err != nil {
		return
	}

	token = t.Access_Token

	return
}
