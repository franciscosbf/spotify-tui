package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscosbf/spotify-tui/internals/api"
	"github.com/franciscosbf/spotify-tui/internals/auth"
	"github.com/franciscosbf/spotify-tui/pkg/config"
)

func welcomeMsg() tea.Cmd {
	return tea.Tick(time.Millisecond*2000, func(_ time.Time) tea.Msg {
		return initMsg(struct{}{})
	})
}

func incrementWelcomeColor(color int) tea.Cmd {
	return tea.Tick(time.Millisecond*400, func(_ time.Time) tea.Msg {
		if color += 1; color < len(welcomeColorsStyle) {
			return welcomeColorMsg(color)
		} else {
			return welcomeColorMsg(0)
		}
	})
}

func incrementAwaitDots(dots int) tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		if dots += 1; dots <= len(dotColorsStyle) {
			return awaitDotsMsg(dots)
		} else {
			return awaitDotsMsg(0)
		}
	})
}

func expirationAlert(expiresIn time.Duration) tea.Cmd {
	d := expiresIn - (expiresIn / 6)

	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return expiredTokenMsg(struct{}{})
	})
}

func readConfig(authConf *config.Config) tea.Cmd {
	return func() tea.Msg {
		if err := authConf.Read(); err != nil {
			return errMsg(err)
		}

		return configReadMsg(struct{}{})
	}
}

func genToken(clientId string, authConf *config.Config) tea.Cmd {
	return func() tea.Msg {
		token, err := api.GenerateToken(clientId)
		if err != nil {
			return errMsg(err)
		}

		if err = authConf.UpdateRefreshToken(token.Refresh); err != nil {
			return errMsg(err)
		}

		return newTokenMsg{token: token}
	}
}

func regenToken(authConf *config.Config) tea.Cmd {
	return func() tea.Msg {
		var (
			token auth.Token
			err   error
		)

		clientId := authConf.ClientId()
		refreshToken := authConf.RefreshToken()

		if token, err = api.RegenerateToken(clientId, refreshToken); err != nil {
			return failedRegenTokenMsg{}
		}

		if err = authConf.UpdateRefreshToken(token.Refresh); err != nil {
			return errMsg(err)
		}

		return newTokenMsg{token: token, refreshed: true}
	}
}

func requestGenToken() tea.Cmd {
	return func() tea.Msg {
		return genTokenMsg(struct{}{})
	}
}

func removeClick() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(_ time.Time) tea.Msg {
		return removeClickMsg(struct{}{})
	})
}

func dismissWarnErr(warnErrId int) tea.Cmd {
	return tea.Tick(time.Second*6, func(_ time.Time) tea.Msg {
		return dismissWarnErrMsg(warnErrId)
	})
}

type clientActions struct {
	client *api.Client
}

func (c clientActions) directOperation(op func() error) tea.Cmd {
	return func() tea.Msg {
		err := op()
		if err, ok := err.(api.ErrResponse); ok && err.Status == 403 {
			return nil
		}
		if err != nil {
			return newWarnErrMsg(err)
		}

		return nil
	}
}

func (c clientActions) setToken(token string) {
	c.client.SetToken(token)
}

func (c clientActions) getUserProfile() tea.Cmd {
	return func() tea.Msg {
		profile, err := c.client.GetUserProfile()
		if err != nil {
			return newWarnErrMsg(err)
		}

		return userInfoMsg(profile)
	}
}

func (c clientActions) resume() tea.Cmd {
	return c.directOperation(c.client.Resume)
}

func (c clientActions) pause() tea.Cmd {
	return c.directOperation(c.client.Pause)
}

func (c clientActions) skipToPrevious() tea.Cmd {
	return c.directOperation(c.client.SkipToPrevious)
}

func (c clientActions) skipToNext() tea.Cmd {
	return c.directOperation(c.client.SkipToNext)
}

func (c clientActions) enableShuffle() tea.Cmd {
	return c.directOperation(c.client.EnableShuffle)
}

func (c clientActions) disableShuffle() tea.Cmd {
	return c.directOperation(c.client.DisableShuffle)
}

func (c clientActions) setRepeatTrack() tea.Cmd {
	return c.directOperation(c.client.SetRepeatTrack)
}

func (c clientActions) setRepeatContext() tea.Cmd {
	return c.directOperation(c.client.SetRepeatContext)
}

func (c clientActions) disableRepeat() tea.Cmd {
	return c.directOperation(c.client.DisableRepeat)
}
