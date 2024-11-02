package api

import (
	"testing"

	"github.com/franciscosbf/spotify-tui/internals/internalstest"
)

func TestGenAndRegenToken(t *testing.T) {
	clientId := internalstest.GetClientId(t)

	token, err := GenerateToken(clientId)
	if err != nil {
		t.Fatalf("failed to generate token: %s", err)
	}

	_, err = RegenerateToken(clientId, token.Refresh)
	if err != nil {
		t.Fatalf("failed to regenerate token: %s", err)
	}
}
