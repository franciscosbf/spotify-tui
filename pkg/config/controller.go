package config

import (
	"errors"

	"github.com/franciscosbf/spotify-tui/internals/config"
)

var ErrMissingClientId = errors.New("missing client id")

type AuthConf struct {
	location string
	auth     config.Auth
}

func (a *AuthConf) Read() error {
	auth, err := config.ParseAuth(a.location)
	if err != nil {
		return err
	}

	if auth.ClientId == "" {
		return ErrMissingClientId
	}

	a.auth = auth

	return nil
}

func (a *AuthConf) UpdateRefreshToken(refreshToken string) error {
	a.auth.RefreshToken = refreshToken

	return config.WriteAuth(a.location, a.auth)
}

func (a *AuthConf) ClientId() string {
	return a.auth.ClientId
}

func (a *AuthConf) RefreshToken() string {
	return a.auth.RefreshToken
}

func NewAuthConf(location string) *AuthConf {
	return &AuthConf{location: location, auth: config.Auth{}}
}
