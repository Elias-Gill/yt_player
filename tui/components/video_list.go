package components

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
)

const (
	modePlaylists = iota
	modeVideos
)

type VideoList struct {
	context *context.Context
	width   int
	height  int

	details VideoInfo

	mode     int
	currItem int
	currPage int
	pages    int
}

func NewList(ctx *context.Context) VideoList {
	return VideoList{
		context: ctx,
		mode:    modeVideos,
		details: NewVideoInfo(ctx),
	}
}

func (l VideoList) Update(msg tea.Msg) (VideoList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height - 9

	case tea.KeyMsg:
		switch msg.String() {
		case "i", tea.KeyCtrlDown.String():
			if len(l.context.Player.Videos) > 0 {
				l.details = l.details.Update(l.context.Player.Videos[l.currItem].Id)
			}

			return l, nil

		case "j", tea.KeyCtrlDown.String():
			if l.currItem+1 < len(l.context.Player.Videos) && l.currItem+1 < (l.currPage+1)*l.height {
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

		case "q", "/", tea.KeyEsc.String():
			l.context.CurrMode = context.SEARCH
		case "m":
			if l.mode == modeVideos {
				l.mode = modePlaylists
			} else {
				l.mode = modeVideos
			}

		// Player controls
		case tea.KeySpace.String():
			l.context.Player.TogglePause()
		case "+", "=":
			l.context.Player.PlusFiveSecs()
		case "-":
			l.context.Player.LessFiveSecs()
		case tea.KeyEnter.String():
			if l.currItem < len(l.context.Player.Videos) {
				l.context.Player.Play(l.currItem)
			}
		}
	}

	l.pages = len(l.context.Player.Videos) / (l.height)
	return l, nil
}

func (l VideoList) View() string {
	videos := l.context.Player.Videos
	if len(videos) == 0 {
		return l.context.Styles.ForegroundGray.Render("No Available Videos ...")
	}

	var msg strings.Builder
	msg.Grow(len(videos) * 30) // Preallocate memory

	for i := l.currPage * l.height; i < (l.currPage+1)*l.height; i++ {
		if i >= len(l.context.Player.Videos) {
			break
		}
		video := videos[i]
		title := video.Title
		if len(title) > l.context.WinWidth-7 {
			title = title[0 : l.context.WinWidth-8]
		}

		// Prepare the line to write
		var line string
		if i == l.currItem {
			if l.context.CurrMode == context.LIST {
				line = l.context.Styles.ForegroundRed.Render(fmt.Sprintf("%d  %s", i+1, title))
			} else {
				line = l.context.Styles.ForegroundGray.Render(fmt.Sprintf("%d  %s", i+1, title))
			}
		} else {
			line = fmt.Sprintf("%d  %s", i+1, title)
		}

		msg.WriteString(line + "\n")
	}

	msg.WriteString("\n")
	for i := 0; i < l.pages; i++ {
		var line string
		if i == l.currPage {
			line = l.context.Styles.ForegroundRed.Render(fmt.Sprintf("%d", i+1))
		} else {
			line = l.context.Styles.ForegroundGray.Render(fmt.Sprintf("%d", i+1))
		}

		line += l.context.Styles.ForegroundGray.Render("  ")
		msg.WriteString(line)
	}

	listPlusInfo := lipgloss.JoinHorizontal(
		0,
		l.context.Styles.Background.
			MaxWidth(int(math.Round(float64(l.context.WinWidth)*0.65))).
			Width(int(math.Round(float64(l.context.WinWidth)*0.65))).
			Render(msg.String()), // List
		l.context.Styles.Background.
			Width(int(math.Round(float64(l.context.WinWidth)*0.35))).
			Render(l.details.View()), // info
	)

	return listPlusInfo
}

func (l VideoList) Init() tea.Cmd {
	return nil
}
