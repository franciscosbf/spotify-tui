package api

import (
	"time"

	"github.com/franciscosbf/spotify-tui/internals/auth"
	"github.com/franciscosbf/spotify-tui/internals/browser"
)

const verificationTimeout = time.Second * 15

func GenerateToken(clientId string) (auth.Token, error) {
	codeVerifier := auth.GenCodeVerifier()
	codeChallenge := auth.GenCodeChallenge(codeVerifier)

	codeAuth := auth.BuildCodeAuth(clientId, codeChallenge)

	if err := browser.OpenAuthLink(codeAuth.Url); err != nil {
		return auth.Token{}, err
	}

	code, err := auth.WaitForCode(codeAuth.State, verificationTimeout)
	if err != nil {
		return auth.Token{}, err
	}

	token, err := auth.FetchToken(clientId, codeVerifier, code)
	if err != nil {
		return auth.Token{}, err
	}

	return token, nil
}

func RegenerateToken(clientId, refreshToken string) (auth.Token, error) {
	return auth.RefreshToken(clientId, refreshToken)
}
