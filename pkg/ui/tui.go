package ui

import (
	"fmt"
	"strings"
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
	authAck
	player
	err
)

type repeatState int

const (
	track repeatState = iota
	context
	off
)

var buttonSymbols = []string{
	"<",
	"▶/‖",
	">",
	"∞",
	"⟳",
}

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
	buttonStyle = lipgloss.NewStyle().
			Bold(true)
	selectedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#fe8019"))
	clickedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#928374"))
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fb4934"))
	errorMsgStyle = lipgloss.NewStyle().
			Italic(true)
	displayStyle = lipgloss.NewStyle().
			Margin(2, 2).
			Padding(2, 2).
			Width(68).
			Height(13).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
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
		key.WithHelp("←", "left"),
	),
	right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "right"),
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
	genTokenMsg         struct{}
	failedRegenTokenMsg struct{}
	removeClickMsg      struct{}
	errMsg              error
	ackedAuthMsg        auth.Token
	userInfoMsg         api.UserProfile
)

type newTokenMsg struct {
	token     auth.Token
	refreshed bool
}

type model struct {
	help           help.Model
	client         *api.Client
	profile        api.UserProfile
	err            error
	authConf       *config.AuthConf
	token          auth.Token
	view           view
	awaitDots      int
	width          int
	height         int
	welcomeColor   int
	selectedButton int
	repeat         repeatState
	clickedButton  bool
	play           bool
	shuffle        bool
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

		return newTokenMsg{token: token}
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

func getUserProfile(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		profile, _ := client.GetUserProfile()

		return userInfoMsg(profile)
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

		return m, requestGenToken()
	case genTokenMsg, failedRegenTokenMsg:
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
		m.token = msg.token
		m.client.SetToken(msg.token.Access)
		if msg.refreshed {
			m.view = player
			return m, expirationAlert(m.token.ExpiresIn)
		} else {
			m.view = authAck
			return m, tea.Batch(expirationAlert(m.token.ExpiresIn),
				getUserProfile(m.client))
		}
	case userInfoMsg:
		m.profile = api.UserProfile(msg)
	case ackedAuthMsg:
		m.view = player
	case removeClickMsg:
		m.clickedButton = false
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
			if m.selectedButton > 0 {
				m.clickedButton = false
				m.selectedButton--
			}
		case key.Matches(msg, playerKm.right):
			if m.selectedButton < len(buttonSymbols)-1 {
				m.clickedButton = false
				m.selectedButton++
			}
		case key.Matches(msg, playerKm.enter):
			m.clickedButton = true
			return m, removeClick()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	display := "\n\n\n\n\n"
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
			"%s\n%s%s",
			awaitStyle.Render("I sent an authorization request!"),
			awaitStyle.Render("Check your browser, I'm waiting for your"),
			dots)
	case authAck:
		if m.profile.Name == "" {
			break
		}
		display += fmt.Sprintf("%s | %d", m.profile.Name, m.profile.Followers.Total)
	case player:
		keyHelp = playerKm
		buttons := []string{}
		for _, symbol := range buttonSymbols {
			buttons = append(buttons, buttonStyle.Render(symbol))
		}
		if m.clickedButton {
			buttons[m.selectedButton] = clickedButtonStyle.Render(buttons[m.selectedButton])
		} else {
			buttons[m.selectedButton] = selectedButtonStyle.Render(buttons[m.selectedButton])
		}
		display += strings.Join(buttons, "   ")
	case err:
		display += fmt.Sprintf("%s %s.",
			errorStyle.Render("Error:"),
			errorMsgStyle.Render(m.err.Error()))
	}

	newLines := "\n\n\n\n"
	if m.view != authConfirmation {
		newLines += "\n"
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
		help:           help.New(),
		client:         api.NewClient(),
		authConf:       config.NewAuthConf(authConfLocation),
		view:           initialization,
		selectedButton: 1,
	}

	return Tui{m}
}

func (t Tui) Start() error {
	_, err := tea.NewProgram(t.m, tea.WithAltScreen()).Run()

	return err
}
