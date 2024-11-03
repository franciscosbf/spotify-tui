package ui

import (
	"fmt"
	"strings"

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
		display += fmt.Sprintf("%s %s",
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
