package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/completition"
	"github.com/elias-gill/yt_player/mpv"
	"github.com/elias-gill/yt_player/yt_api"
)

type model struct {
	list      list.Model
	textInput textinput.Model
	window    window
	width     int
	height    int

	playing     bool
	firstSearch bool

	songStatus string
	curSong    string
}

type window int

type TickMsg time.Time

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

func NewModel() tea.Model {
	m := model{
		window:      SEARCH_WINDOW,
		textInput:   NewInput(),
		list:        NewList(),
		firstSearch: true,
		playing:     false,
	}

	return m
}

// to make a timer for updating mpv actual time
func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickEvery()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.songStatus = mpv.GetSongStatus()
		m.list.SetSize(m.width, m.height-2)
		return m, tickEvery()

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
		m.textInput.Update(msg)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // quit app and close the player
			m.window = QUIT
			return m, tea.Quit

		case "ctrl+q": // quit "detaching" mpv
			mpv.DetachPlayer()
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
		completition.PersistHistory()

		return QuitTextStyle.Render(
			"Closing player instance and quitting... ")

	case SEARCH_WINDOW:
		return TitleStyle.Render(searchTitle) + "\n\n" + m.textInput.View()

	default: // MODE_LISTING
		text := ""

		// to be honest, i dont know why the pause and play
		// are inverted but it works
		if m.playing {
			text += pause
		} else {
			text += play
		}

		m.list.Title = text + m.curSong + "\t" + m.songStatus

		return m.list.View()
	}
}

func (m model) updateTrackList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "/": // enter search mode
		m.window = SEARCH_WINDOW
		return m, nil

	case tea.KeySpace.String(): // toggle pause
		mpv.TogglePause()
		m.playing = !m.playing

		return m, nil

	case "+", "=": // +5 secs
		mpv.PlusFiveSecs()

	case "-": // -5 secs
		mpv.LessFiveSecs()

	case "esc", "q": // select song
		m.window = QUIT
		return m, nil

	case "enter": // select song
		curItem, ok := m.list.SelectedItem().(ListItem)
		if !ok {
			return m, nil
		}

		m.curSong = curItem.Title()
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
		value := m.textInput.Value()

		m.list = generateVideoList(value)
		completition.AddHistoryEntry(value)

		m.window = LISTING_WINDOW
		m.firstSearch = false

		return m, nil

	case tea.KeyUp.String():
		actual := m.textInput.Value()
		m.textInput.SetValue(completition.NextEntry(actual))

	case tea.KeyDown.String():
		actual := m.textInput.Value()
		m.textInput.SetValue(completition.PrevEntry(actual))
	}

	m.textInput, _ = m.textInput.Update(msg)
	return m, nil
}
