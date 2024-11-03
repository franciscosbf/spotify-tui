package ui

import (
	"fmt"
	"math/rand/v2"
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
	nRepeatStates
)

type button int

const (
	previous button = iota
	resumePause
	next
	shuffle
	repeat
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
			Bold(true)
	buttonStyle = lipgloss.NewStyle().
			Bold(true)
	selectedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#fe8019"))
	clickedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#928374"))
	ackStyle = lipgloss.NewStyle().
			Bold(true)
	warnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fabd2f"))
	warnMsgStyle = lipgloss.NewStyle().
			Bold(true)
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
		key.WithHelp("↵", "press"),
	),
}

type ackKeyMap struct {
	defaultKeyMap
	enter key.Binding
}

var ackKm = ackKeyMap{
	defaultKeyMap: defaultKm,
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "confirm"),
	),
}

func (k ackKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.enter}
}

func (k ackKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
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
	dismissWarnErrMsg   int
)

type newTokenMsg struct {
	token     auth.Token
	refreshed bool
}

type warnErrMsg struct {
	err error
	id  int
}

func newWarnErrMsg(err error) warnErrMsg {
	id := rand.Int()

	return warnErrMsg{err, id}
}

func newNoWarnErrMsg() warnErrMsg {
	return warnErrMsg{nil, -1}
}

func (w warnErrMsg) warn() bool {
	return w.id != -1
}

type model struct {
	help           help.Model
	actions        clientActions
	profile        api.UserProfile
	currentWarnErr warnErrMsg
	err            error
	conf           *config.Config
	token          auth.Token
	view           view
	awaitDots      int
	width          int
	height         int
	welcomeColor   int
	selectedButton int
	repeat         repeatState
	clickedButton  bool
	resume         bool
	shuffle        bool
}

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

func dismissWarnErr(warnErrId int) tea.Cmd {
	return tea.Tick(time.Second*6, func(_ time.Time) tea.Msg {
		return dismissWarnErrMsg(warnErrId)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(welcomeMsg(),
		incrementWelcomeColor(m.welcomeColor))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case warnErrMsg:
		m.currentWarnErr = msg
		return m, dismissWarnErr(msg.id)
	case dismissWarnErrMsg:
		if m.currentWarnErr.id == int(msg) {
			m.currentWarnErr = newNoWarnErrMsg()
		}
	case errMsg:
		m.view = err
		m.err = msg
	case welcomeColorMsg:
		if m.view == initialization {
			m.welcomeColor = int(msg)
			return m, incrementWelcomeColor(m.welcomeColor)
		}
	case initMsg:
		return m, readConfig(m.conf)
	case configReadMsg:
		if m.conf.RefreshToken() != "" {
			return m, regenToken(m.conf)
		}

		return m, requestGenToken()
	case genTokenMsg, failedRegenTokenMsg:
		m.view = authConfirmation
		clientId := m.conf.ClientId()
		return m, tea.Batch(genToken(clientId, m.conf),
			incrementAwaitDots(m.awaitDots))
	case awaitDotsMsg:
		m.awaitDots = int(msg)
		if m.view == authConfirmation {
			return m, incrementAwaitDots(m.awaitDots)
		}
	case expiredTokenMsg:
		return m, regenToken(m.conf)
	case newTokenMsg:
		m.token = msg.token
		m.actions.setToken(msg.token.Access)
		if msg.refreshed {
			m.view = player
			return m, expirationAlert(m.token.ExpiresIn)
		} else {
			m.view = authAck
			return m, tea.Batch(expirationAlert(m.token.ExpiresIn),
				m.actions.getUserProfile())
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

		switch m.view {
		case authAck:
			switch {
			case key.Matches(msg, playerKm.enter):
				m.view = player
			}
		case player:
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

				var cmd tea.Cmd

				switch button(m.selectedButton) {
				case previous:
					cmd = m.actions.skipToPrevious()
				case resumePause:
					if m.resume {
						cmd = m.actions.resume()
					} else {
						cmd = m.actions.pause()
					}
					m.resume = !m.resume
				case next:
					cmd = m.actions.skipToNext()
				case shuffle:
					if m.shuffle {
						cmd = m.actions.enableShuffle()
					} else {
						cmd = m.actions.disableShuffle()
					}
					m.shuffle = !m.shuffle
				case repeat:
					switch m.repeat {
					case track:
						cmd = m.actions.setRepeatTrack()
					case context:
						cmd = m.actions.setRepeatContext()
					case off:
						cmd = m.actions.disableRepeat()
					}
					m.repeat = (m.repeat + 1) % nRepeatStates
				}

				return m, tea.Batch(cmd, removeClick())
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	display := ""

	if m.currentWarnErr.warn() {
		warn := fmt.Sprintf("%s %s",
			warnStyle.Render("Alert:"),
			warnMsgStyle.Render(m.currentWarnErr.err.Error()))
		display = fmt.Sprintf("%s%s", warn, display)
	}

	display += "\n\n\n\n\n"

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
		keyHelp = ackKm
		name := m.profile.Name
		followers := m.profile.Followers.Total
		if name == "" {
			display += fmt.Sprintf("%s\n", ackStyle.Render("Verified with success!"))
		} else {
			fllw := "follower"
			if followers != 1 {
				fllw += "s"
			}
			display += fmt.Sprintf("%s\n%s",
				ackStyle.Render(fmt.Sprintf("Thank you %s with your shitty %d %s!",
					name, followers, fllw)),
				ackStyle.Render("Now I can spoof your miserable account."))
		}
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
	if m.view != authConfirmation && m.view != authAck {
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

func New(confLocation string) Tui {
	client := api.NewClient()

	m := model{
		help:           help.New(),
		actions:        clientActions{client},
		conf:           config.NewConfig(confLocation),
		currentWarnErr: newNoWarnErrMsg(),
		view:           initialization,
		selectedButton: 1,
		shuffle:        true,
		resume:         true,
		repeat:         track,
	}

	return Tui{m}
}

func (t Tui) Start() error {
	_, err := tea.NewProgram(t.m, tea.WithAltScreen()).Run()

	return err
}
