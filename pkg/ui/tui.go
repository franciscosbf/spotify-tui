package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Tui struct {
	m model
}

func New(confLocation string) Tui {
	m := newModel(confLocation)

	return Tui{m}
}

func (t Tui) Start() error {
	_, err := tea.NewProgram(t.m, tea.WithAltScreen()).Run()

	return err
}
