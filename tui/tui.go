package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const yt_link = "https://www.youtube.com/watch?v="

type model struct {
	list      list.Model
	textInput textinput.Model
	cmd       *exec.Cmd

	query string

	quitting  bool
	searching bool
	playing   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			m.quitting = true
			return m, tea.Quit

		case "ctrl+c":
			if m.playing {
				m.cmd.Process.Kill()
				m.playing = false
				m.list.Title = "Select youtube music"
				return m, nil
			}

			return m, nil

		case "enter":
			if m.searching {
				m.list = generateVideoList(m.textInput.Value())
				m.searching = false
			} else {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					if m.cmd != nil {
						m.cmd.Process.Kill()
					}

					choice := yt_link + string(i.id)

					cmd := exec.Command("mpv", "--no-video", choice)
					if cmd == nil {
						fmt.Println("Cannot find 'mpv' player on your machine.")
						os.Exit(1)
					}

					m.cmd = cmd

					if err := cmd.Start(); err != nil {
						fmt.Println("Error running mpv: ", err.Error())
						os.Exit(1)
					}

					m.playing = true
					m.list.Title = "▷  Select youtube music"
				}
			}

			return m, nil

		case "/":
			m.searching = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.searching {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		m.cmd.Process.Kill()
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}

	if m.searching {
		return fmt.Sprintf(
			"Search for music on youtube\n\n%s\n",
			m.textInput.View(),
		)
	}

	return "\n" + m.list.View()
}

func InitTUI() {
	m := model{
		searching: true,
		playing:   false,
		textInput: initInput(),
		list:      list.New([]list.Item{}, itemDelegate{}, 80, 25),
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
