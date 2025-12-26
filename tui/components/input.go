package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui/messages"
)

type Input struct {
	textInput textinput.Model
	inputHelp help.Model
	prompt    string

	active bool
}

const prompt = " Search> "

func NewInput() Input {
	ti := textinput.New()
	ti.Prompt = context.Styles.BackgroundRed.Render(prompt) + " "
	ti.Focus()

	return Input{
		textInput: ti,
	}
}

func (i Input) View() string {
	i.textInput.Prompt = context.Styles.BackgroundGray.Render(prompt) + " "

	return i.textInput.View()
}

func (i Input) Init() tea.Cmd {
	return nil
}

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.textInput.Width = msg.Width - len(i.textInput.Prompt)
	case messages.ModeChangedMessage:
		if msg.IsModeSearch() {
			i.active = true
		}
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			err := i.Player.Search(i.textInput.Value())
			i.ctx.Error = err
			return i, func() tea.Msg {
				return messages.ModeChangedMessage{
					Mode: messages.LIST,
				}
			}
		}
	}

	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)

	return i, cmd
}
