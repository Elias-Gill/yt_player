package components

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/player"
)

const (
	modePlaylists = iota
	modeVideos
)

type VideoList struct {
	context *context.Context
	width   int
	height  int

	details string

	mode     int
	currItem int
	currPage int
	pages    int
}

func NewVideoList(ctx *context.Context) VideoList {
	return VideoList{
		context: ctx,
		mode:    modeVideos,
		details: "To retrive video information press 'i'",
	}
}

func (l VideoList) Update(msg tea.Msg) (VideoList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height - 10

	case tea.KeyMsg:
		switch msg.String() {
		case "i", tea.KeyCtrlDown.String():
			if len(l.context.Player.Videos) > 0 {
				l.updateVideoDetails()
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

	listWidth := int(math.Round(float64(l.context.WinWidth) * 0.62))
	infoWidth := int(math.Round(float64(l.context.WinWidth) * 0.35))

	var list strings.Builder
	list.Grow(len(videos) * 30) // Preallocate memory

	for i := l.currPage * l.height; i < (l.currPage+1)*l.height; i++ {
		if i >= len(l.context.Player.Videos) {
			break
		}

		video := videos[i]
		title := video.Title

		// Prepare the line to write
		var line string
		if i == l.currItem {
			if l.context.CurrMode == context.LIST {
				line = l.context.Styles.ForegroundRed.MaxWidth(listWidth - 1).Render(fmt.Sprintf("%d  %s", i+1, title))
			} else {
				line = l.context.Styles.ForegroundGray.MaxWidth(listWidth - 1).Render(fmt.Sprintf("%d  %s", i+1, title))
			}
		} else {
			line = l.context.Styles.Background.MaxWidth(listWidth - 1).Render(fmt.Sprintf("%d  %s", i+1, title))
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

	// Display context errors alongside pagination
	var err = ""
	if l.context.Error != nil {
		err = l.context.Error.Error()
	}
	pagination.WriteString("  " + l.context.Styles.Background.Render(err))

	listPlusInfo := lipgloss.JoinHorizontal(
		0,
		l.context.Styles.Background. // list
						MaxWidth(listWidth).
						Width(listWidth).
						Render(list.String()),
		lipgloss.NewStyle(). // info
					Width(infoWidth).
					PaddingLeft(1).
					MaxWidth(infoWidth).
					Height(l.height).
					MaxHeight(l.height).
					Render(l.details),
	)

	return lipgloss.JoinVertical(0,
		l.context.Styles.Background.Render(listPlusInfo),
		l.context.Styles.Background.Width(l.width).Render(pagination.String()))
}

func (l VideoList) Init() tea.Cmd {
	return nil
}

func (l *VideoList) updateVideoDetails() {
	video := l.context.Player.Videos[l.currItem].Id
	details, err := player.GetVideoDetails(video)
	if err != nil {
		l.details = "Cannot retrieve video details\n" + err.Error()
	}

	l.details = fmt.Sprintf("Author: %s\nDuration: %s\nTitle: %s\nDescription: %s",
		details.Author,
		details.Duration,
		details.Title,
		details.Description)
}
