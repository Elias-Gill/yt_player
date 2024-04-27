package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func NewInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Your music"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	return ti
}
