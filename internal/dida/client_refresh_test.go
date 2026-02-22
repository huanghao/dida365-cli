package dida

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExchangeRefreshToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); !strings.HasPrefix(got, "Basic ") {
			t.Fatalf("expected basic auth header")
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.Form.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected grant_type: %s", r.Form.Get("grant_type"))
		}
		_ = json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "new-access",
			RefreshToken: "new-refresh",
			TokenType:    "bearer",
			ExpiresIn:    3600,
		})
	}))
	defer ts.Close()

	old := oauthTokenURL
	oauthTokenURL = ts.URL
	defer func() { oauthTokenURL = old }()

	resp, err := ExchangeRefreshToken("cid", "secret", "refresh-1", "tasks:read tasks:write")
	if err != nil {
		t.Fatalf("ExchangeRefreshToken returned error: %v", err)
	}
	if resp.AccessToken != "new-access" {
		t.Fatalf("unexpected token response: %+v", resp)
	}
}

func TestExchangeRefreshToken_RequireRefreshToken(t *testing.T) {
	_, err := ExchangeRefreshToken("cid", "secret", "", "")
	if err == nil {
		t.Fatalf("expected error")
	}
}
