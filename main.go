package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/elias-gill/yt_player/globals"
	"github.com/elias-gill/yt_player/tui"
	"github.com/elias-gill/yt_player/yt_api"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/mpv"
)

type window int

const (
	SEARCH_WINDOW = iota
	LISTING_WINDOW
	QUIT
)

var (
	play        = "▷  Playing: "
	pause       = "⏸︎  Pause: "
	searchTitle = "Search for music on youtube"
)

type model struct {
	list      list.Model
	textInput textinput.Model
	window    window
	curSong   string

	playing     bool
	firstSearch bool
}

func main() {
	globals.ParseFlags()

	mpv.StartPlayer()

	m := model{
		window:      SEARCH_WINDOW,
		textInput:   tui.NewInput(),
		list:        tui.NewList(),
		firstSearch: true,
		playing:     false,
	}

	// NOTE: tea.Every() for updating every x seconds

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
		m.textInput.Update(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // quit app and close the player
			m.window = QUIT
			return m, tea.Quit

		default:
			switch m.window {
			case SEARCH_WINDOW:
				return m.updateSearchInput(msg)

			case LISTING_WINDOW:
				return m.updateTrackList(msg)

			default:
				m.window = QUIT
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.window {
	case QUIT:
		mpv.StopPlayer()

		// TODO: delete
		go func() {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}()

		return tui.QuitTextStyle.Render(
			"Closing player instance and quitting... ")

	case SEARCH_WINDOW:
		return tui.TitleStyle.Render(searchTitle) + "\n\n" + m.textInput.View()

	default: // MODE_LISTING
		text := "\n\n"

		if m.playing {
			text += play
		} else {
			text += pause
		}

		m.list.Title = text + m.curSong

		return m.list.View()
	}
}

func (m model) updateTrackList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "/": // enter searching mode
		m.window = SEARCH_WINDOW
		return m, nil

	case tea.KeySpace.String(): // toggle pause
		mpv.TogglePause()
		m.playing = !m.playing

		return m, nil

	case "+", "=": // select song
		mpv.PlusFiveSecs()

	case "esc", "q": // select song
		m.window = QUIT
		return m, nil

	case "-": // select song
		mpv.LessFiveSecs()

	case "enter": // select song
		curItem, ok := m.list.SelectedItem().(tui.ListItem)
		if !ok {
			return m, nil
		}

		m.curSong = string(curItem.Title())
		mpv.ChangeSong(yt_api.Yt_url + curItem.Id())

		return m, nil
	}

	m.list, _ = m.list.Update(msg)
	return m, nil
}

func (m model) updateSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "esc": // cancel search
		if m.firstSearch {
			m.window = QUIT
			return m, nil
		}

		m.window = LISTING_WINDOW
		return m, nil

	case "enter":
		m.list = tui.GenerateVideoList(m.textInput.Value())

		m.window = LISTING_WINDOW
		m.firstSearch = false

		return m, nil

	default:
		m.textInput, _ = m.textInput.Update(msg)
		return m, nil
	}
}
