package config

import (
	"os"
	"testing"
)

func prepareTempFile(pattern string, data string, t *testing.T) string {
	tempDir := os.TempDir()

	tempFile, err := os.CreateTemp(tempDir, pattern)
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}

	defer tempFile.Close()

	if _, err := tempFile.Write([]byte(data)); err != nil {
		t.Fatalf("failed to write temp file: %s", err)
	}

	return tempFile.Name()
}

func TestParseAuth(t *testing.T) {
	expectedClientId := "v4395hb49b4b"
	expectedRefreshToken := "v3rb45jh549h84"

	data := `{"client_id": "` + expectedClientId +
		`", "refresh_token": "` + expectedRefreshToken + `"}`
	filename := prepareTempFile("TestParseAuth", data, t)

	auth, err := ParseAuth(filename)
	if err != nil {
		t.Fatalf("failed to parse auth: %s", err)
	}

	if auth.ClientId == "" {
		t.Fatalf("missing client_id")
	}
	if auth.ClientId != expectedClientId {
		t.Fatalf("invalid client_id. got=%s, expected=%s",
			auth.ClientId, expectedClientId)
	}

	if auth.RefreshToken == "" {
		t.Fatalf("missing refresh_token")
	}
	if auth.RefreshToken != expectedRefreshToken {
		t.Fatalf("invalid refresh_token. got=%s, expected=%s",
			auth.RefreshToken, expectedRefreshToken)
	}
}
