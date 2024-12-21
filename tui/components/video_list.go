package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/context"
)

type VideoList struct {
	context *context.Context
	width   int
	height  int
}

func NewList(ctx *context.Context) VideoList {
	return VideoList{
		context: ctx,
		width:   30,
		height:  30,
	}
}

func (l VideoList) Update(msg tea.Msg) (VideoList, tea.Cmd) {
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
			if l.context.CurrItem < len(l.context.Player.Videos) {
				l.context.Player.Play(l.context.CurrItem)
			}
		}
	}

	return l, nil
}

func (l VideoList) View() string {
	videos := l.context.Player.Videos
	if len(videos) == 0 {
		return l.context.Styles.ForegroundGray.Render("No Available Videos ...")
	}

	var msg strings.Builder
	msg.Grow(len(videos) * 30) // Preallocate memory

	for i, video := range videos {
		if i >= l.height-1 {
			break
		}

		// Prepare the line to write
		var line string
		if i == l.context.CurrItem {
			if l.context.CurrMode == context.LIST {
				line = l.context.Styles.ForegroundRed.Render(fmt.Sprintf("%d\t%s", i+1, video.Title))
			} else {
				line = l.context.Styles.ForegroundGray.Render(fmt.Sprintf("%d\t%s", i+1, video.Title))
			}
		} else {
			line = fmt.Sprintf("%d\t%s", i+1, video.Title)
		}

		msg.WriteString(line + "\n")
	}

	return msg.String()
}

func (l VideoList) Init() tea.Cmd {
	return nil
}
