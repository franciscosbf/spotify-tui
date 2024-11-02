package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/franciscosbf/spotify-tui/internals/util"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

const codeVerifierLength = 128

func GenCodeVerifier() string {
	r := util.NewRand()
	code := make([]byte, codeVerifierLength)

	for i := range codeVerifierLength {
		code[i] = chars[r.Int63()%int64(len(chars))]
	}

	return string(code)
}

func GenCodeChallenge(codeVerifier string) string {
	bts := sha256.Sum256([]byte(codeVerifier))
	challenge := base64.StdEncoding.EncodeToString(bts[:])

	r := strings.NewReplacer("=", "", "+", "-", "/", "_")
	challenge = r.Replace(challenge)

	return challenge
}
