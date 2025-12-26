package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui/components"
	"github.com/elias-gill/yt_player/tui/messages"
)

type Tui struct {
	textInput components.Input
	list      components.VideoList
	player    components.PlayerProgress
	history   components.HistoryList

	currMode  messages.Mode
	WinWidth  int
	WinHeight int
}

func NewModel() tea.Model {
	return Tui{
		textInput: components.NewInput(),
		list:      components.NewVideoList(),
		player:    components.NewPlayer(),
		history:   components.NewHistoryList(),
	}
}

func (t Tui) View() string {
	if t.currMode == messages.HISTORY {
		return t.history.View()
	}

	style := context.Styles.Background.
		Padding(1).
		MaxWidth(t.WinWidth).
		Width(t.WinWidth).
		Height(t.WinHeight).
		MaxHeight(t.WinHeight)

	if t.currMode == messages.HELP {
		helpMessage := `
	Keybinds:
		- "\": Enter history mode
		- Enter: Select entry/search
		- j/k or Arrow Keys: Navigate up and down
		- h/l or Arrow Keys: Go to next page or previous page
		- Tab: Cycle between search and list
		- Ctrl+C: Quit (detach from player if the --detach-on-quit flag is given)
		- q: Similar to "Tab"

	Player Controls:
		- Space: Pause player
		- "+": Skip forward 5 seconds
		- "-": Skip backward 5 seconds

		`
		contMsg := context.Styles.ForegroundRed.Render("... Press any key to continue")
		return style.Render(helpMessage + contMsg)
	}

	return style.Render(
		lipgloss.JoinVertical(
			0,
			t.textInput.View(), // Input
			context.Styles.Background. // List + info
							Height(t.WinHeight-5).
							MaxHeight(t.WinHeight-5).
							Width(t.WinWidth).
							MaxWidth(t.WinWidth).
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

		case "\\":
			t.currMode = messages.HISTORY
			return t, nil

		case "~":
			t.currMode = messages.HELP
			return t, nil

		case tea.KeyTab.String():
			if t.currMode == messages.LIST {
				t.currMode = messages.SEARCH
			} else {
				t.currMode = messages.LIST
			}

			return t, nil

		default:
			switch t.currMode {
			case messages.LIST:
				t.list, _ = t.list.Update(msg)

			case messages.SEARCH:
				t.textInput, _ = t.textInput.Update(msg)

			case messages.HISTORY:
				t.history, _ = t.history.Update(msg)

			case messages.HELP: // Scape help on any key press
				t.currMode = messages.LIST
			}
		}

	case tea.WindowSizeMsg:
		t.WinHeight = msg.Height
		t.WinWidth = msg.Width

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
