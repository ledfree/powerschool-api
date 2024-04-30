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

func TestTimeCheck(t *testing.T) {
	ps := ApiConfig{TestUrl, TestClientId, TestClientSecret}

	_, v, err := ps.TimeCheck()
	if err != nil {
		t.Fatalf("fail : TimeCheck - %s", err)
	} else if v != "don't know" {
		t.Fatalf("fail : TimeCheck value = %s", v)
	}
}

/*
func TestGetAreas(t *testing.T) {
	ps := ApiConfig{TestUrl, TestClientId, TestClientSecret}

	s, err := ps.GetTest()
	if err != nil {
		t.Fatalf("fail : status = %d - %s", s, err)
	} else if s != http.StatusOK {
		t.Fatalf("fail : status = %d - %s", s, err)
	}
}
*/
