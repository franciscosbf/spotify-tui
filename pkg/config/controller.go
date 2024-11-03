package config

import (
	"errors"

	"github.com/franciscosbf/spotify-tui/internals/config"
)

var ErrMissingClientId = errors.New("missing client_id field in config")

type Config struct {
	location string
	conf     config.Config
}

func (a *Config) Read() error {
	auth, err := config.Parse(a.location)
	if err != nil {
		return err
	}

	if auth.ClientId == "" {
		return ErrMissingClientId
	}

	a.conf = auth

	return nil
}

func (a *Config) UpdateRefreshToken(refreshToken string) error {
	a.conf.RefreshToken = refreshToken

	return config.Write(a.location, a.conf)
}

func (a *Config) ClientId() string {
	return a.conf.ClientId
}

func (a *Config) RefreshToken() string {
	return a.conf.RefreshToken
}

func NewConfig(location string) *Config {
	return &Config{location: location, conf: config.Config{}}
}
