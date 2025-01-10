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
	list      components.VideoList
	player    components.PlayerProgress
	history   components.HistoryList
}

func NewModel(ctx *context.Context) tea.Model {
	return Tui{
		context:   ctx,
		textInput: components.NewInput(ctx),
		list:      components.NewVideoList(ctx),
		player:    components.NewPlayer(ctx),
		history:   components.NewHistoryList(ctx),
	}
}

func (t Tui) View() string {
	if t.context.CurrMode == context.HISTORY {
		return t.history.View()
	}

	return t.context.Styles.Background.
		Padding(1).
		MaxWidth(t.context.WinWidth).
		Width(t.context.WinWidth).
		Height(t.context.WinHeight).
		MaxHeight(t.context.WinHeight).
		Render(
			lipgloss.JoinVertical(
				0,
				t.textInput.View(), // Input
				t.context.Styles.Background. // List + info
								Height(t.context.WinHeight-5).
								MaxHeight(t.context.WinHeight-5).
								Width(t.context.WinWidth).
								MaxWidth(t.context.WinWidth).
								PaddingTop(1).
								Render(t.list.View()),
				t.player.View(), // Player
			))
}

func (t Tui) Init() tea.Cmd {
	return t.player.Init()
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return t, tea.Quit

		case "~":
			t.context.CurrMode = context.HISTORY
			return t, nil

		case tea.KeyTab.String():
			if t.context.CurrMode == context.LIST {
				t.context.CurrMode = context.SEARCH
			} else {
				t.context.CurrMode = context.LIST
			}

			return t, nil

		default:
			switch t.context.CurrMode {
			case context.LIST:
				t.list, _ = t.list.Update(msg)

			case context.SEARCH:
				t.textInput, _ = t.textInput.Update(msg)

			case context.HISTORY:
				t.history, _ = t.history.Update(msg)
			}
		}

	case tea.WindowSizeMsg:
		t.context.WinHeight = msg.Height
		t.context.WinWidth = msg.Width

		t.list, _ = t.list.Update(msg)
		t.textInput, _ = t.textInput.Update(msg)
		t.player, _ = t.player.Update(msg)
		t.history, _ = t.history.Update(msg)

		return t, nil

	default:
		t.player, cmd = t.player.Update(msg)
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
