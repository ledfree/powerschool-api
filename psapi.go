package psapi

import "net/url"

func GetAccessToken(psurl, clientId, clientSecret string) (status int, token string, err error) {
	status = 400 // bad request

	u, err := url.Parse(psurl)
	if err != nil {
		return
	}

	token = u.Hostname()

	return
}
