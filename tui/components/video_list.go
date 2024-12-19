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
		l.width = msg.Width
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

	if len(videos) == 0 {
		return style.Render(l.context.Styles.GruvboxGray.Render("No Available Videos ..."))
	}

	msg := ""
	for i, video := range videos {
		if i >= l.height-1 {
			break
		}

		line := fmt.Sprintf("%d\t%s", i, video.Title)
		if i == l.context.CurrItem {
			if l.context.CurrMode == context.LIST {
				line = l.context.Styles.GruvboxOrange.Render(line)
			} else {
				line = l.context.Styles.GruvboxGray.Render(line)
			}
		}

		msg += line + "\n"
	}

	return style.Render(msg)
}

func (l List) Init() tea.Cmd {
	return nil
}
