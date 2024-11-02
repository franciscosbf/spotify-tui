package auth

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/franciscosbf/spotify-tui/internals/uri"
)

func TestBuildAuth(t *testing.T) {
	fakeClientId := "242b42obfevbeber"
	fakeCodeChallenge := "vrwnvr3vnreov3vr3"

	codeAuth := BuildCodeAuth(fakeClientId, fakeCodeChallenge)

	expectedUrl := uri.ACCOUNTS + "/authorize"
	if !strings.HasPrefix(codeAuth.Url, expectedUrl) {
		t.Fatalf("url doesn't start with %s. got=%s", expectedUrl, codeAuth.Url)
	}

	authUrl, err := url.Parse(codeAuth.Url)
	if err != nil {
		t.Fatalf("failed to parse auth url: %s", err)
	}

	query := authUrl.Query()

	expectedParameters := map[string]string{
		"client_id":             fakeClientId,
		"response_type":         "code",
		"redirect_uri":          uri.REDIRECT,
		"state":                 codeAuth.State,
		"scope":                 strings.Join(requiredScopes, " "),
		"code_challenge_method": "S256",
		"code_challenge":        fakeCodeChallenge,
	}

	for key, expectedValue := range expectedParameters {
		value := query.Get(key)

		if value == "" {
			t.Fatalf("missing query parameter %s", key)
		}

		if value != expectedValue {
			t.Fatalf("invalid query parameter %s. got=%s, expected=%s", key, value, expectedValue)
		}
	}
}

func TestWaitForCode(t *testing.T) {
	fakeCode := "br3rb3h5b34b3bnb"
	fakeState := "dvsdab33b44t4btadfasf"
	timeout := time.Second * 4

	type codeResult struct {
		err  error
		code string
	}

	cre := make(chan codeResult, 1)
	resp := make(chan *http.Response, 1)
	stop := make(chan struct{}, 1)

	go func() {
		code, err := WaitForCode(fakeState, timeout)
		cre <- codeResult{err, code}
		stop <- struct{}{}
	}()

	redirectUri, _ := url.Parse(uri.REDIRECT)
	redirectUri = redirectUri.JoinPath("callback")
	query := redirectUri.Query()
	query.Set("state", fakeState)
	query.Set("code", fakeCode)
	redirectUri.RawQuery = query.Encode()

	go func(redirectUri string) {
		for {
			if response, err := http.Get(redirectUri); err == nil {
				resp <- response
				break
			}

			select {
			case <-stop:
			default:
				continue
			}
			break
		}
	}(redirectUri.String())

	result := <-cre

	if result.err != nil {
		t.Fatalf("failed while waiting for code: %s", result.err)
	}

	if result.code != fakeCode {
		t.Fatalf("invalid code. got=%s, expected=%s", result.code, fakeCode)
	}

	response := <-resp

	defer func() {
		response.Body.Close()
	}()

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read body")
	}

	expectedBody := "Go check the app..."
	if body := string(rawBody); body != expectedBody {
		t.Fatalf("invalid body. got=%s, expected=%s", body, expectedBody)
	}
}
