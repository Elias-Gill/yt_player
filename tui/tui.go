package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui/components"
)

type Tui struct {
	context    *context.Context
	textInput  components.Input
	list       components.VideoList
	playerInfo components.PlayerProgress
}

func NewModel(ctx *context.Context) tea.Model {
	return Tui{
		context:    ctx,
		textInput:  components.NewInput(ctx),
		list:       components.NewList(ctx),
		playerInfo: components.NewPlayerInfo(ctx),
	}
}

func (t Tui) View() string {
	return t.context.Styles.GruvboxBg.
		Width(t.context.WinWidth).
		Height(t.context.WinHeight).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				t.textInput.View(),
				t.list.View(),
				lipgloss.NewStyle().Align(lipgloss.Bottom).Render(t.playerInfo.View()),
			))
}

func (t Tui) Init() tea.Cmd {
	return t.playerInfo.Init()
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return t, tea.Quit
		case tea.KeyTab.String():
			t.context.NextMode()
			return t, nil
		}

	case tea.WindowSizeMsg:
		t.context.WinHeight = msg.Height
		t.context.WinWidth = msg.Width

		t.list, _ = t.list.Update(msg)
		t.textInput, _ = t.textInput.Update(msg)
		t.playerInfo, _ = t.playerInfo.Update(msg)

		return t, nil
	}

	t.playerInfo, cmd = t.playerInfo.Update(msg)

	switch t.context.CurrMode {
	case context.LIST:
		t.list, _ = t.list.Update(msg)
	case context.SEARCH:
		t.textInput, _ = t.textInput.Update(msg)
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
