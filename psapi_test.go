package psapi

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	ps := ApiConfig{TestUrl, TestClientId, TestClientSecret}

	s, _, err := ps.GetAccessToken()
	if err != nil {
		t.Fatalf("fail : GetAccessToken - status code = %d; %s", s, err)
	}
	if s != 200 {
		t.Fatalf("fail : GetAccessToken - status code = %d; %s", s, err)
	}
}
