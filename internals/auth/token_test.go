package auth

import (
	"testing"
	"time"

	"github.com/franciscosbf/spotify-tui/internals/browser"
	"github.com/franciscosbf/spotify-tui/internals/internalstest"
)

func TestFetchTokenAndRefreshToken(t *testing.T) {
	codeVerifier := GenCodeVerifier()
	codeChallenge := GenCodeChallenge(codeVerifier)
	timeout := time.Second * 15

	clientId := internalstest.GetClientId(t)

	codeAuth := BuildCodeAuth(clientId, codeChallenge)

	if err := browser.OpenAuthLink(codeAuth.Url); err != nil {
		t.Fatalf("failed to open browser with auth url: %s", err)
	}

	code, err := WaitForCode(codeAuth.State, timeout)
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

		expectedExpiresIn := time.Second * 3600
		if token.ExpiresIn == expectedExpiresIn {
			t.Fatalf("invalid expire time. got=%s, expected=%s",
				token.ExpiresIn, expectedExpiresIn)
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
