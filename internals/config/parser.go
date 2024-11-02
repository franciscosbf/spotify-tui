package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Auth struct {
	ClientId     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

var (
	ErrFailedToReadConfig = errors.New("failed to read config file")
	ErrInvalidConfig      = errors.New("config file is invalid")
)

func parse(path string, config any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ErrFailedToReadConfig
	}

	if err := json.Unmarshal(raw, config); err != nil {
		return ErrInvalidConfig
	}

	return nil
}

func ParseAuth(path string) (Auth, error) {
	var auth Auth

	if err := parse(path, &auth); err != nil {
		return Auth{}, err
	}

	return auth, nil
}
