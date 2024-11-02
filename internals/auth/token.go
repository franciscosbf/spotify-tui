package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/franciscosbf/spotify-tui/internals/uri"
)

var (
	ErrTokenRequestFailed   = errors.New("failed to request token")
	ErrTokenInvalidResponse = errors.New("invalid token response")
)

type Token struct {
	Access  string
	Refresh string
	Expires time.Duration
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func requestToken(clientId string, parameters *url.Values) (Token, error) {
	tokenUrl, _ := url.Parse(uri.ACCOUNTS)

	tokenUrl = tokenUrl.JoinPath("api").JoinPath("token")

	parameters.Set("client_id", clientId)
	body := bytes.NewBuffer([]byte(parameters.Encode()))

	request, _ := http.NewRequest("POST", tokenUrl.String(), body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return Token{}, ErrTokenRequestFailed
	}

	var tokenMeta tokenResponse
	json.NewDecoder(response.Body).Decode(&tokenMeta)
	token := Token{
		Access:  tokenMeta.AccessToken,
		Refresh: tokenMeta.RefreshToken,
		Expires: time.Duration(tokenMeta.ExpiresIn),
	}

	return token, nil
}

func FetchToken(clientId, codeVerifier, code string) (Token, error) {
	parameters := &url.Values{}
	parameters.Set("grant_type", "authorization_code")
	parameters.Set("code", code)
	parameters.Set("redirect_uri", uri.REDIRECT)
	parameters.Set("code_verifier", codeVerifier)

	return requestToken(clientId, parameters)
}

func RefreshToken(clientId, refreshToken string) (Token, error) {
	parameters := &url.Values{}
	parameters.Set("grant_type", "refresh_token")
	parameters.Set("refresh_token", refreshToken)

	return requestToken(clientId, parameters)
}
