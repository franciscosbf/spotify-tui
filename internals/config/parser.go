package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	ClientId     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

var (
	ErrFailedToReadConfig  = errors.New("failed to read config file")
	ErrFailedToWriteConfig = errors.New("failed to  writeconfig file")
	ErrInvalidConfig       = errors.New("config file is invalid")
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

func write(path string, config any) error {
	raw, _ := json.MarshalIndent(config, "", "  ")

	if err := os.WriteFile(path, raw, 0644); err != nil {
		return ErrFailedToWriteConfig
	}

	return nil
}

func Parse(path string) (Config, error) {
	var auth Config

	if err := parse(path, &auth); err != nil {
		return Config{}, err
	}

	return auth, nil
}

func Write(path string, auth Config) error {
	return write(path, auth)
}
