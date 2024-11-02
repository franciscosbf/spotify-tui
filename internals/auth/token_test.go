package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pkg/browser"
)

func TestFetchTokenAndRefreshToken(t *testing.T) {
	codeVerifier := GenCodeVerifier()
	codeChallenge := GenCodeChallenge(codeVerifier)
	timeout := time.Second * 15

	fmt.Println(codeVerifier)
	fmt.Println(codeChallenge)

	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		t.Fatal("a valid CLIENT_ID must be set")
	}

	auth := BuildAuth(clientId, codeChallenge)

	if err := browser.OpenURL(auth.Url); err != nil {
		t.Fatalf("failed to open browser with auth url: %s", err)
	}

	code, err := WaitForCode(auth.State, timeout)
	if err != nil {
		t.Fatalf("failed to wait for code: %s", err)
	}

	validateToken := func(token Token) {
		if token.Access == "" {
			t.Fatal("access token is empty")
		}

		if token.Refresh == "" {
			t.Fatal("refresh token is empty")
		}

		if token.Expires == 0 {
			t.Fatal("expire is zero")
		}
	}

	token, err := FetchToken(clientId, codeVerifier, code)
	if err != nil {
		t.Fatalf("failed to fetch token: %s", err)
	}
	validateToken(token)

	token, err = RefreshToken(clientId, token.Refresh)
	if err != nil {
		t.Fatalf("failed to refresh token: %s", err)
	}
	validateToken(token)
}
