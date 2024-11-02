package browser

import (
	"errors"

	"github.com/pkg/browser"
)

var ErrBrowser = errors.New("failed to open browser with authorization link")

func OpenAuthLink(authUrl string) error {
	if err := browser.OpenURL(authUrl); err != nil {
		return ErrBrowser
	}

	return nil
}
