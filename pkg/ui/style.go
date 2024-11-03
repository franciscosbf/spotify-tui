package ui

import "github.com/charmbracelet/lipgloss"

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
