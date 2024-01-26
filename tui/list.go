package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/globals"
	"github.com/elias-gill/yt_player/yt_api"
)

type item struct {
	title string
	id    string
}

// filter is deactivated here
func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func generateVideoList(input string) list.Model {
	items := []list.Item{}

	videos := yt_api.RetrieveVideos(
		input,
		globals.GetMaxResults(),
		globals.GetApiKey(),
	)

	for key, value := range videos {
		items = append(items, item{id: value, title: key})
	}

	l := list.New(items, itemDelegate{}, 80, 25)
	l.Title = "Select youtube music"

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}
