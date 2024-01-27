package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/mpv"
)

type mode int

const (
	MODE_SEARCHING = iota
	MODE_LISTING
	MODE_QUITTING
)

var (
	defaultTitle = "Select a song to reproduce"
	playingTitle = "▷  Playing: "
	searchTitle  = "Search for music on youtube"
)

type model struct {
	list      list.Model
	textInput textinput.Model
	cmd       *exec.Cmd

	mode    mode
	playing bool

	firstSearch bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		m.textInput.Update(msg)
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case MODE_SEARCHING:
			return m.updateSearchInput(msg)

		case MODE_LISTING:
			return m.updateTrackList(msg)

		default:
			m.mode = MODE_QUITTING
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.mode {
	case MODE_QUITTING:
		mpv.KillInstance(m.cmd)
		return quitTextStyle.Render(
			"Closing all player instances and quitting... ")

	case MODE_SEARCHING:
		return titleStyle.Render(searchTitle) + "\n\n" + m.textInput.View()

	case MODE_LISTING:
		return "\n" + m.list.View()
	}

	// unrecheable code
	return ""
}

func InitTUI() {
	l := list.New([]list.Item{}, itemDelegate{}, 80, 25)
	l.Styles.Title = titleStyle

	m := model{
		mode:        MODE_SEARCHING,
		playing:     false,
		textInput:   newInput(),
		firstSearch: true,
		list:        l,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) updateTrackList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "/": // enter searching mode
		m.mode = MODE_SEARCHING
		return m, nil

	case "ctrl+q": // quit app without closing the player
		return m, tea.Quit

	case "esc", "q": // quit app and close the player
		m.mode = MODE_QUITTING
		return m, tea.Quit

	case "ctrl+c": // kill player
		mpv.KillInstance(m.cmd)

		m.cmd = nil
		m.playing = false
		m.list.Title = defaultTitle

		return m, nil

	case "enter": // select song
		curItem, ok := m.list.SelectedItem().(item)
		if !ok {
			return m, nil
		}

		// Kill the previous mpv process if exists
		if m.playing {
			mpv.KillInstance(m.cmd)
			m.cmd = nil
		}

		m.cmd = mpv.NewPlayer(string(curItem.Id()))

		m.playing = true
		m.list.Title = playingTitle + string(curItem.Title())

		return m, nil
	}

	m.list, _ = m.list.Update(msg)
	return m, nil
}

func (m model) updateSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+q", "ctrl+c": // quit app
		m.mode = MODE_QUITTING
		return m, tea.Quit

	case "esc": // cancel search
		if !m.firstSearch {
			m.mode = MODE_LISTING
			return m, nil
		}

		m.mode = MODE_QUITTING
		return m, tea.Quit

	case "enter":
		m.list = generateVideoList(m.textInput.Value())

		m.mode = MODE_LISTING
		m.firstSearch = false

		return m, nil

	default:
		m.textInput, _ = m.textInput.Update(msg)
		return m, nil
	}
}
