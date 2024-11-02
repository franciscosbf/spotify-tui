package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/franciscosbf/spotify-tui/internals/api"
	"github.com/franciscosbf/spotify-tui/internals/auth"
	"github.com/franciscosbf/spotify-tui/pkg/config"
)

type view int

const (
	initialization view = iota
	authConfirmation
	player
	err
)

var (
	welcomeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#98971a"))
	awaitAuthStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fe8019"))
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fb4934"))
	displayStyle = lipgloss.NewStyle().
			Margin(2, 2).
			Padding(2, 2).
			Width(54).
			Height(9).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#928374"))
)

type (
	initMsg             struct{}
	configReadMsg       struct{}
	awaitDotsMsg        int
	expiredTokenMsg     struct{}
	regenTokenMsg       struct{}
	failedRegenTokenMsg struct{}
	errMsg              error
	newTokenMsg         auth.Token
)

type model struct {
	err       error
	authConf  *config.AuthConf
	token     auth.Token
	view      view
	awaitDots int
}

func welcomeMsg() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
		return initMsg(struct{}{})
	})
}

func incrementAwaitDots(dots int) tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		if dots += 1; dots < 4 {
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

func readConfig(authConf *config.AuthConf) tea.Cmd {
	return func() tea.Msg {
		if err := authConf.Read(); err != nil {
			return errMsg(err)
		}

		return configReadMsg(struct{}{})
	}
}

func genToken(clientId string, authConf *config.AuthConf) tea.Cmd {
	return func() tea.Msg {
		token, err := api.GenerateToken(clientId)
		if err != nil {
			return errMsg(err)
		}

		if err = authConf.UpdateRefreshToken(token.Refresh); err != nil {
			return errMsg(err)
		}

		return newTokenMsg(token)
	}
}

func regenToken(authConf *config.AuthConf) tea.Cmd {
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

		return newTokenMsg(token)
	}
}

func requestRegenToken() tea.Cmd {
	return func() tea.Msg {
		return regenTokenMsg(struct{}{})
	}
}

func (m model) Init() tea.Cmd {
	return welcomeMsg()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.view = err
		m.err = msg
	case initMsg:
		return m, readConfig(m.authConf)
	case configReadMsg:
		if m.authConf.RefreshToken() != "" {
			return m, regenToken(m.authConf)
		}

		return m, requestRegenToken()
	case regenTokenMsg, failedRegenTokenMsg:
		m.view = authConfirmation
		clientId := m.authConf.ClientId()
		return m, tea.Batch(genToken(clientId, m.authConf),
			incrementAwaitDots(m.awaitDots))
	case awaitDotsMsg:
		m.awaitDots = int(msg)
		if m.view == authConfirmation {
			return m, incrementAwaitDots(m.awaitDots)
		}
	case expiredTokenMsg:
		return m, regenToken(m.authConf)
	case newTokenMsg:
		m.view = player
		m.token = auth.Token(msg)
		return m, expirationAlert(m.token.ExpiresIn)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyLeft:
		// TODO:
		case tea.KeyRight:
			// TODO:
		}
	}

	return m, nil
}

func (m model) View() string {
	display := ""

	switch m.view {
	case initialization:
		display += welcomeStyle.Render("Welcome to Spotify TUI!!!")
	case authConfirmation:
		display += awaitAuthStyle.Render("I'm waiting for your") +
			strings.Repeat(".", m.awaitDots)
	case player:
		// TODO:
		display += "TODO"
	case err:
		display += errorStyle.Render("Error:") + " " + m.err.Error()
	}

	return displayStyle.Render(display)
}

type Tui struct {
	m model
}

func New(authConfLocation string) Tui {
	m := model{
		authConf: config.NewAuthConf(authConfLocation),
		view:     initialization,
	}

	return Tui{m}
}

func (t Tui) Start() error {
	_, err := tea.NewProgram(t.m).Run()

	return err
}
