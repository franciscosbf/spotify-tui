package ui

import "github.com/charmbracelet/bubbles/key"

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
