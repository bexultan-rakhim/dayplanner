package model

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Back     key.Binding
	New      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Status   key.Binding
	Graph    key.Binding
	Group    key.Binding
	Undo     key.Binding
	Redo     key.Binding
	Quit     key.Binding
	NextField key.Binding
	PrevField key.Binding
	Confirm  key.Binding
	Regenerate key.Binding
}

var Keys = KeyMap{
	Up:         key.NewBinding(key.WithKeys("up", "k")),
	Down:       key.NewBinding(key.WithKeys("down", "j")),
	Select:     key.NewBinding(key.WithKeys("enter")),
	Back:       key.NewBinding(key.WithKeys("esc")),
	New:        key.NewBinding(key.WithKeys("n")),
	Edit:       key.NewBinding(key.WithKeys("e")),
	Delete:     key.NewBinding(key.WithKeys("d")),
	Status:     key.NewBinding(key.WithKeys(" ")),
	Graph:      key.NewBinding(key.WithKeys("g")),
	Group:      key.NewBinding(key.WithKeys("G")),
	Undo:       key.NewBinding(key.WithKeys("ctrl+z")),
	Redo:       key.NewBinding(key.WithKeys("ctrl+y")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c")),
	NextField:  key.NewBinding(key.WithKeys("tab")),
	PrevField:  key.NewBinding(key.WithKeys("shift+tab")),
	Confirm:    key.NewBinding(key.WithKeys("enter")),
	Regenerate: key.NewBinding(key.WithKeys("ctrl+r")),
}
