package psapi

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	ps := ApiConfig{TestUrl, TestClientId, TestClientSecret}

	s, token, err := ps.GetAccessToken()
	if err != nil {
		t.Fatalf("fail : GetAccessToken - %s", err)
	} else {
		t.Logf("Status       : %d", s)
		t.Logf("Bearer Token : %s", token)
	}
}
