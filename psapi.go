package psapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	strGetToken        string = `https://%s/oauth/access_token/`
	strTimeCheck       string = `https://%s/ws/v1/time`
	strAreas           string = `https://%s/ws/schema/area`
	strTables          string = `https://%s/ws/schema/table`
	strApplicationJson string = `application/json`
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

func (t tokenResponse) Authorization() string {
	return fmt.Sprintf("%s %s", t.Token_Type, t.Access_Token)
}

type ApiConfig struct {
	PsUrl        string
	ClientId     string
	ClientSecret string
}

func (a ApiConfig) URL() (u *url.URL, err error) {
	u, err = url.Parse(a.PsUrl)
	return
}

type resourceTimeStruct struct {
	Resource struct {
		Time      string `json:"time"`
		TimeStamp int64  `json:"timestamp"`
	} `json:"resource"`
}

type schemaInfoStruct struct {
	Name        string             `json:"name"`
	UserVisible bool               `json:"userVisible,omitempty"`
	Display     string             `json:"display,omitempty"`
	SubAreas    []schemaInfoStruct `json:"subareas,omitempty"`
	PrimaryKey  string             `json:"primaryKey,omitempty"`
}

func (s schemaInfoStruct) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name : %s\n", s.Name))
	for _, v := range s.SubAreas {
		sb.WriteString(fmt.Sprintf("  %s\n", v.String()))
	}

	return sb.String()
}

type schemaStruct struct {
	Areas  []schemaInfoStruct `json:"areas,omitempty"`
	Tables []schemaInfoStruct `json:"tables,omitempty"`
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

	u, er := a.URL()
	if er != nil {
		return status, t, fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
	}

	urlTime := fmt.Sprintf(strTimeCheck, u.Hostname())

	status, timeInfo, resBody, er := fetch[resourceTimeStruct](http.MethodGet, urlTime, headers, emptyBody)
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

func (a ApiConfig) GetAccessToken() (status int, token tokenResponse, err error) {
	headers := make(map[string]string)
	body := `grant_type=client_credentials`

	headers["Authorization"] = fmt.Sprintf("Basic %s", a.toBase64())
	headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8"
	headers["Accept"] = "application/json"

	u, er := a.URL()
	if er != nil {
		err = fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
		return
	}

	urlToken := fmt.Sprintf(strGetToken, u.Hostname())

	status, token, resBody, er := fetch[tokenResponse](http.MethodPost, urlToken, headers, body)
	if er != nil {
		var r authErrorDetails
		er := json.Unmarshal(resBody, &r)
		if er != nil {
			err = fmt.Errorf("failed : status = %d : %s", status, r.ErrorDescription)
		} else {
			err = fmt.Errorf("failed : status = %d : %s", status, er)
		}
	}

	return
}

func (a ApiConfig) GetAreas() (status int, d schemaStruct, err error) {
	s, t, err := a.GetAccessToken()
	if err != nil {
		return
	}
	if s != http.StatusOK {
		return
	}

	headers := make(map[string]string)
	headers["Authorization"] = t.Authorization()
	headers["Accept"] = strApplicationJson
	headers["Content-Type"] = strApplicationJson

	u, er := a.URL()
	if er != nil {
		err = fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
		return
	}

	urlAreas := fmt.Sprintf(strAreas, u.Hostname())

	status, d, _, err = fetch[schemaStruct](http.MethodGet, urlAreas, headers, "")
	if err != nil {
		err = fmt.Errorf("status = %d - from request url %s - error : %s", status, urlAreas, err)
		return
	}

	return
}

func (a ApiConfig) GetTables() (status int, d schemaStruct, err error) {
	s, t, err := a.GetAccessToken()
	if err != nil {
		return
	}
	if s != http.StatusOK {
		return
	}

	headers := make(map[string]string)
	headers["Authorization"] = t.Authorization()
	headers["Accept"] = strApplicationJson
	headers["Content-Type"] = strApplicationJson

	u, er := a.URL()
	if er != nil {
		err = fmt.Errorf("invalid URL %s : %s", a.PsUrl, er)
		return
	}

	urlTables := fmt.Sprintf(strTables, u.Hostname())

	status, d, _, err = fetch[schemaStruct](http.MethodGet, urlTables, headers, "")
	if err != nil {
		err = fmt.Errorf("status = %d - from request url %s - error : %s", status, urlTables, err)
		return
	}

	return
}
