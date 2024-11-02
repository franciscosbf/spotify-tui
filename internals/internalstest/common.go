package internalstest

import (
	"os"
	"testing"
)

func GetClientId(t *testing.T) string {
	clientId := os.Getenv("CLIENT_ID")

	if clientId == "" {
		t.Fatal("a valid CLIENT_ID must be set")
	}

	return clientId
}
