package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/franciscosbf/spotify-tui/internals/uri"
	"github.com/franciscosbf/spotify-tui/internals/util"
)

var requiredScopes = []string{
	"user-read-playback-state",
	"user-modify-playback-state",
	"user-read-currently-playing",
	"playlist-read-private",
}

var (
	ErrInvalidAuth = errors.New("invalid authentication")
	ErrAuthTimeout = errors.New("authentication timed out")
)

type CodeAuth struct {
	Url   string
	State string
}

type callbackResponse struct {
	code  string
	state string
	error string
}

type callbackServer struct {
	*http.Server
	response <-chan callbackResponse
}

func genState() string {
	r := util.NewRand()

	return strconv.FormatInt(r.Int63(), 10)
}

func startCallbackServer() callbackServer {
	url, _ := url.Parse(uri.REDIRECT)
	port := url.Port()

	result := make(chan callbackResponse, 1)

	handler := http.NewServeMux()
	handler.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		result <- callbackResponse{
			code:  query.Get("code"),
			state: query.Get("state"),
			error: query.Get("error"),
		}

		fmt.Fprint(w, "Go check the app...")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	go func() {
		server.ListenAndServe()
	}()

	return callbackServer{server, result}
}

func BuildCodeAuth(clientId, codeChallenge string) CodeAuth {
	tokenUrl, _ := url.Parse(uri.ACCOUNTS)

	tokenUrl = tokenUrl.JoinPath("authorize")

	state := genState()
	scope := strings.Join(requiredScopes, " ")

	query := url.Values{}
	query.Set("client_id", clientId)
	query.Set("response_type", "code")
	query.Set("redirect_uri", uri.REDIRECT)
	query.Set("state", state)
	query.Set("scope", scope)
	query.Set("code_challenge_method", "S256")
	query.Set("code_challenge", codeChallenge)
	tokenUrl.RawQuery = query.Encode()

	return CodeAuth{Url: tokenUrl.String(), State: state}
}

func WaitForCode(state string, timeout time.Duration) (string, error) {
	callback := startCallbackServer()

	defer func() {
		callback.Close()
	}()

	select {
	case response := <-callback.response:
		if response.error != "" || response.state != state {
			return "", ErrInvalidAuth
		}

		return response.code, nil
	case <-time.After(timeout):
		return "", ErrAuthTimeout
	}
}
