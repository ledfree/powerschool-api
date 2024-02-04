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

const (
	strGetToken string = "https://%s/oauth/access_token/"
)

type authErrorDetails struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

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

	u, er := url.Parse(a.PsUrl)
	if er != nil {
		err = fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
		return
	}

	urlToken := fmt.Sprintf(strGetToken, u.Hostname())
	body := []byte(`grant_type=client_credentials`)
	req, er := http.NewRequest(http.MethodPost, urlToken, bytes.NewBuffer(body))
	if er != nil {
		err = fmt.Errorf("failed to create new request %s : %s", urlToken, er)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", a.toBase64()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Accept", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, er := client.Do(req)
	if er != nil {
		err = fmt.Errorf("failed to Client.Do : %s", er)
		return
	}
	defer res.Body.Close()

	status = res.StatusCode

	resBody, er := io.ReadAll(res.Body)
	if er != nil {
		err = fmt.Errorf("status = %d : failed reading body : %s", res.StatusCode, er)
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status = %d : %s", res.StatusCode, resBody)
		return
	} else {

	}

	if res.StatusCode == http.StatusOK {
		var t tokenResponse
		er := json.Unmarshal(resBody, &t)
		if er != nil {
			err = fmt.Errorf("failed to translate response body to token : %s", er)
			return
		}

		token = t.Access_Token
	} else {
		var r authErrorDetails
		er := json.Unmarshal(resBody, &r)
		if er != nil {
			err = fmt.Errorf("failed : status = %d : %s", res.StatusCode, r.ErrorDescription)
		} else {
			err = fmt.Errorf("failed : status = %d : %s", res.StatusCode, er)
		}
		return
	}

	return
}
