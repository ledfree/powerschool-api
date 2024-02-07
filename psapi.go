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
	strGetToken  string = "https://%s/oauth/access_token/"
	strTimeCheck string = "https://%s/ws/v1/time"
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

type resourceTimeStruc struct {
	Resource struct {
		Time      string `json:"time"`
		TimeStamp string `json:"timestamp"`
	} `json:"resource"`
}

func fetch[T any](httpAction, url string, headerValues map[string]string, strBody string) (status int, v T, resBody []byte, err error) {
	req, er := http.NewRequest(httpAction, url, bytes.NewBuffer(([]byte(strBody))))
	if er != nil {
		err = fmt.Errorf("failed to create new request %s : %s", url, er)
		return
	}
	for k, s := range headerValues {
		req.Header.Add(k, s)
	}

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

	resBody, er = io.ReadAll(res.Body)
	if er != nil {
		err = fmt.Errorf("status = %d : failed reading body : %s", res.StatusCode, er)
		return
	}

	if res.StatusCode == http.StatusOK {
		er := json.Unmarshal(resBody, &v)
		if er != nil {
			err = fmt.Errorf("failed to translate response body : %s", er)
		}
	} else if res.StatusCode == http.StatusNoContent {
		// do nothing, mybe return no content
	} else {
		err = fmt.Errorf("status = %d - body : %s", res.StatusCode, string(resBody))
	}

	return
}

func (a ApiConfig) TimeCheck() (status int, t string, err error) {
	headers := make(map[string]string)
	var emptyBody string

	headers["Accept"] = "application/json"

	u, er := url.Parse(a.PsUrl)
	if er != nil {
		return status, t, fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
	}

	urlTime := fmt.Sprintf(strTimeCheck, u.Hostname())

	status, timeInfo, resBody, er := fetch[resourceTimeStruc](http.MethodGet, urlTime, headers, emptyBody)
	if er != nil {
		err = fmt.Errorf("status = %d : %s - %s", status, resBody, er)
		return
	}

	if status != http.StatusOK {
		err = fmt.Errorf("status = %d : %s", status, resBody)
		return
	}

	if status == http.StatusOK {
		t = timeInfo.Resource.Time
	}

	return
}

func (a ApiConfig) toBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.ClientId, a.ClientSecret)))
}

func (a ApiConfig) GetAccessToken() (status int, token string, err error) {
	headers := make(map[string]string)
	body := `grant_type=client_credentials`

	headers["Authorization"] = fmt.Sprintf("Basic %s", a.toBase64())
	headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8"
	headers["Accept"] = "application/json"

	u, er := url.Parse(a.PsUrl)
	if er != nil {
		err = fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
		return
	}

	urlToken := fmt.Sprintf(strGetToken, u.Hostname())

	status, tokenInfo, resBody, er := fetch[tokenResponse](http.MethodPost, urlToken, headers, body)
	if er != nil {
		var r authErrorDetails
		er := json.Unmarshal(resBody, &r)
		if er != nil {
			err = fmt.Errorf("failed : status = %d : %s", status, r.ErrorDescription)
		} else {
			err = fmt.Errorf("failed : status = %d : %s", status, er)
		}
	} else {
		token = tokenInfo.Access_Token
	}

	return
}
