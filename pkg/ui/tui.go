package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	welcomeColorsStyle = []lipgloss.Style{
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#98971a")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#689d6a")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#8ec07c")),
	}
	welcomeStyle = lipgloss.NewStyle().
			Bold(true)
	dotColorsStyle = []lipgloss.Style{
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fabd2f")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fb4934")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fb4934")),
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fb4934")),
	}
	awaitStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#a89984"))
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fb4934"))
	errorMsgStyle = lipgloss.NewStyle().
			Italic(true)
	displayStyle = lipgloss.NewStyle().
			Margin(2, 2).
			Padding(2, 2).
			Width(82).
			Height(15).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#928374"))
)

type defaultKeyMap struct {
	quit key.Binding
}

func (k defaultKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k defaultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.quit},
	}
}

var defaultKm = defaultKeyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
}

type playerKeyMap struct {
	defaultKeyMap
	left  key.Binding
	right key.Binding
	enter key.Binding
}

func (k playerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.left, k.right, k.enter}
}

func (k playerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var playerKm = playerKeyMap{
	defaultKeyMap: defaultKm,
	left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move left"),
	),
	right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move right"),
	),
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "enter"),
	),
}

type (
	welcomeColorMsg     int
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
	help         help.Model
	err          error
	authConf     *config.AuthConf
	token        auth.Token
	view         view
	awaitDots    int
	width        int
	height       int
	welcomeColor int
}

func welcomeMsg() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
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
	return tea.Batch(welcomeMsg(),
		incrementWelcomeColor(m.welcomeColor))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.view = err
		m.err = msg
	case welcomeColorMsg:
		if m.view == initialization {
			m.welcomeColor = int(msg)
			return m, incrementWelcomeColor(m.welcomeColor)
		}
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
		switch {
		case key.Matches(msg, defaultKm.quit):
			return m, tea.Quit
		}

		if m.view != player {
			break
		}

		switch {
		case key.Matches(msg, playerKm.left):
		// TODO:
		case key.Matches(msg, playerKm.right):
			// TODO:
		case key.Matches(msg, playerKm.enter):
			// TODO:
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	display := "\n\n\n\n"
	var keyHelp help.KeyMap = defaultKm

	switch m.view {
	case initialization:
		colored := welcomeColorsStyle[m.welcomeColor].Render("Welcome to Spotify TUI")
		display += welcomeStyle.Render(colored)
	case authConfirmation:
		dots := ""
		for _, style := range dotColorsStyle[:m.awaitDots] {
			dots += style.Render(".")
		}
		display += fmt.Sprintf(
			"\n%s\n%s%s",
			awaitStyle.Render("I sent an authorization request!"),
			awaitStyle.Render("Check your browser, I'm waiting for your"),
			dots)
	case player:
		// TODO:
		keyHelp = playerKm
		display += "TODO"
	case err:
		display += fmt.Sprintf("%s %s.",
			errorStyle.Render("Error:"),
			errorMsgStyle.Render(m.err.Error()))
	}

	newLines := ""
	if m.view == authConfirmation {
		newLines = "\n\n"
	} else {
		newLines = "\n\n\n"
	}
	display += fmt.Sprintf("%s%s", newLines, m.help.View(keyHelp))

	display = displayStyle.Render(display)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, display)
}

type Tui struct {
	m model
}

func New(authConfLocation string) Tui {
	m := model{
		help:     help.New(),
		authConf: config.NewAuthConf(authConfLocation),
		view:     initialization,
	}

	return Tui{m}
}

func (t Tui) Start() error {
	_, err := tea.NewProgram(t.m, tea.WithAltScreen()).Run()

	return err
}
