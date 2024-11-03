package ui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/franciscosbf/spotify-tui/internals/api"
	"github.com/franciscosbf/spotify-tui/pkg/config"
)

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
