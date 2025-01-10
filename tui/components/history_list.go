package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
)

type HistoryList struct {
	context *context.Context
	width   int
	height  int

	currItem int
	currPage int
	pages    int
}

func NewHistoryList(ctx *context.Context) HistoryList {
	return HistoryList{
		context: ctx,
	}
}

func (l HistoryList) Update(msg tea.Msg) (HistoryList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width - 2 // minus padding
		l.height = msg.Height - 10

	case tea.KeyMsg:
		history := l.context.Player.GetHistory()

		switch msg.String() {
		case "j", tea.KeyCtrlDown.String():
			if l.currItem+1 < len(history) && l.currItem+1 < (l.currPage+1)*l.height {
				l.currItem++
			}
		case "k", tea.KeyCtrlPgUp.String():
			if l.currItem > 0 && l.currItem > l.currPage*l.height {
				l.currItem--
			}
		case "l", tea.KeyRight.String():
			if l.currPage+1 < l.pages {
				l.currPage++
				l.currItem = l.currPage * l.height
			}
		case "h", tea.KeyLeft.String():
			if l.currPage > 0 {
				l.currPage--
				l.currItem = l.currPage * l.height
			}

		case "q", tea.KeyEsc.String():
			l.context.CurrMode = context.LIST

		case tea.KeyEnter.String():
			if l.currItem < len(history) {
				item := history[l.currItem]
				l.context.Player.Playlists = item.Playlists
				l.context.Player.Videos = item.Videos
			}

			l.context.CurrMode = context.LIST
		}
	}

	l.pages = len(l.context.Player.GetHistory()) / (l.height)
	return l, nil
}

func (l HistoryList) View() string {
	entries := l.context.Player.GetHistory()
	if entries == nil || len(entries) == 0 {
		return l.context.Styles.ForegroundGray.Render("No Available History ...")
	}

	var list strings.Builder
	list.Grow(len(entries) * 30) // Preallocate memory

	for i := l.currPage * l.height; i < (l.currPage+1)*l.height; i++ {
		if i >= len(entries) {
			break
		}

		entry := entries[i]
		title := entry.Input

		// Prepare the line to write
		var line string
		if i == l.currItem {
			line = l.context.Styles.ForegroundRed.MaxWidth(l.width).Render(fmt.Sprintf("%d  %s", i+1, title))
		} else {
			line = l.context.Styles.Background.MaxWidth(l.width).Render(fmt.Sprintf("%d  %s", i+1, title))
		}

		list.WriteString(line + "\n")
	}

	var pagination strings.Builder
	list.Grow(30)
	for i := 0; i < l.pages; i++ {
		var line string
		if i == l.currPage {
			line = l.context.Styles.ForegroundRed.Render(fmt.Sprintf("%d", i+1))
		} else {
			line = l.context.Styles.ForegroundGray.Render(fmt.Sprintf("%d", i+1))
		}

		line += l.context.Styles.ForegroundGray.Render("  ")
		pagination.WriteString(line)
	}

	return l.context.Styles.Background.
		Height(l.context.WinHeight).
		Padding(1).
		Render(
			lipgloss.JoinVertical(0,
				l.context.Styles.BackgroundGray.Render(" Select History Entry "),
				l.context.Styles.Background.
					Padding(1).
					MaxWidth(l.width).
					Width(l.width).
					Render(list.String()),
				l.context.Styles.Background.
					Width(l.width).
					Render(pagination.String()),
			))
}

func (l HistoryList) Init() tea.Cmd {
	return nil
}
