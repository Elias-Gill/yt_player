package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui/components"
)

type Tui struct {
	context   *context.Context
	textInput components.Input
	list      components.List
}

func NewModel(ctx *context.Context) tea.Model {
	return Tui{
		context:   ctx,
		textInput: components.NewInput(ctx),
		list:      components.NewList(ctx),
	}
}

func (t Tui) View() string {
	return lipgloss.JoinHorizontal(
		0,
		lipgloss.JoinVertical(
			0, t.textInput.View(),
			t.list.View(),
		),
		drawDivision(t.context.WinHeight))
}

func (t Tui) Init() tea.Cmd {
	return nil
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if tea.KeyCtrlC == msg.Type {
			return t, tea.Quit
		}
	case tea.WindowSizeMsg:
		t.context.WinHeight = msg.Height
		t.context.WinWidth = msg.Width

		t.list, _ = t.list.Update(msg)
		t.textInput, _ = t.textInput.Update(msg)

		return t, nil
	}

	switch t.context.CurrMode {
	case context.LIST:
		t.list, cmd = t.list.Update(msg)
	case context.SEARCH:
		t.textInput, cmd = t.textInput.Update(msg)
	}

	return t, cmd
}

func drawDivision(h int) string {
	var div string
	for i := 0; i < h-1; i++ {
		div += "│\n"
	}

	return div
}
