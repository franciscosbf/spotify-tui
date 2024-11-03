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

	conf, err := Parse(filename)
	if err != nil {
		t.Fatalf("failed to parse auth: %s", err)
	}

	if conf.ClientId == "" {
		t.Fatalf("missing client_id")
	}
	if conf.ClientId != expectedClientId {
		t.Fatalf("invalid client_id. got=%s, expected=%s",
			conf.ClientId, expectedClientId)
	}

	if conf.RefreshToken == "" {
		t.Fatalf("missing refresh_token")
	}
	if conf.RefreshToken != expectedRefreshToken {
		t.Fatalf("invalid refresh_token. got=%s, expected=%s",
			conf.RefreshToken, expectedRefreshToken)
	}
}

func TestPersistAuth(t *testing.T) {
	data := Config{
		ClientId:     "v4395hb49b4b",
		RefreshToken: "v3rb45jh549h84",
	}

	filename := prepareTempFile("TestStoreAuth", "", t)

	if err := Write(filename, data); err != nil {
		t.Fatalf("failed to store file: %s", err)
	}

	conf, err := Parse(filename)
	if err != nil {
		t.Fatalf("failed to parse auth: %s", err)
	}

	if conf.ClientId == "" {
		t.Fatalf("missing client_id")
	}
	if conf.ClientId != data.ClientId {
		t.Fatalf("invalid client_id. got=%s, expected=%s",
			conf.ClientId, data.ClientId)
	}

	if conf.RefreshToken == "" {
		t.Fatalf("missing refresh_token")
	}
	if conf.RefreshToken != data.RefreshToken {
		t.Fatalf("invalid refresh_token. got=%s, expected=%s",
			conf.RefreshToken, conf.RefreshToken)
	}
}
