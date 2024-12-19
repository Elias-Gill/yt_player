package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
)

type List struct {
	context *context.Context
	width   int
	height  int
}

func NewList(ctx *context.Context) List {
	return List{
		context: ctx,
		width:   30,
		height:  30,
	}
}

func (l List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = int(float32(msg.Width) * 0.6)
		l.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if l.context.CurrItem+1 < len(l.context.Player.Videos) {
				l.context.CurrItem++
			}
		case "k":
			if l.context.CurrItem > 0 {
				l.context.CurrItem--
			}

		case "q", "/", tea.KeyEsc.String():
			l.context.CurrMode = context.SEARCH

		case tea.KeyEnter.String():
			l.context.Player.Play(l.context.CurrItem)
		}
	}

	return l, nil
}

func (l List) View() string {
	videos := l.context.Player.Videos
	style := lipgloss.NewStyle().
		MaxWidth(l.width - 1).
		AlignHorizontal(lipgloss.Left).
		PaddingTop(1)

	msg := ""
	for i, video := range videos {
		if i >= l.height {
			break
		}
		line := fmt.Sprintf("%d\t%s", i, video.Title)

		if i == l.context.CurrItem {
			if l.context.CurrMode == context.LIST {
				line = lipgloss.NewStyle().Foreground(lipgloss.Color(context.GruvboxOrange)).Render(line)
			} else {
				line = lipgloss.NewStyle().Foreground(lipgloss.Color(context.GruvboxGray)).Render(line)
			}
		}

		msg += line + "\n"
	}

	return style.Render(msg)
}

func (l List) Init() tea.Cmd {
	return nil
}
