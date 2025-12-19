package ui

import (
	"fmt"

	"github.com/ortolanph/tv-shows-manager/models"
	"github.com/ortolanph/tv-shows-manager/storage"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	listView state = iota
	addView
	editView
)

type model struct {
	storage   *storage.CSVStorage
	shows     []models.Show
	list      list.Model
	form      formModel
	state     state
	err       error
	editIndex int
}

func Run(csvFile string) error {
	s := storage.NewCSVStorage(csvFile)
	shows, err := s.Load()
	if err != nil {
		return err
	}

	items := showsToItems(shows)
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "TV Shows Manager"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	m := model{
		storage: s,
		shows:   shows,
		list:    l,
		state:   listView,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case listView:
		return m.updateList(msg)
	case addView, editView:
		return m.updateForm(msg)
	}
	return m, nil
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.state = addView
			m.form = newFormModel(models.Show{})
			return m, m.form.Init()
		case "e":
			if len(m.shows) > 0 {
				idx := m.list.Index()
				if idx < len(m.shows) {
					m.state = editView
					m.editIndex = idx
					m.form = newFormModel(m.shows[idx])
					return m, m.form.Init()
				}
			}
		case "d":
			if len(m.shows) > 0 {
				idx := m.list.Index()
				if idx < len(m.shows) {
					m.shows = append(m.shows[:idx], m.shows[idx+1:]...)
					m.storage.Save(m.shows)
					m.list.SetItems(showsToItems(m.shows))
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = listView
			return m, nil
		case "ctrl+s":
			show := m.form.getShow()
			if m.state == addView {
				m.shows = append(m.shows, show)
			} else {
				m.shows[m.editIndex] = show
			}
			m.storage.Save(m.shows)
			m.list.SetItems(showsToItems(m.shows))
			m.state = listView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case listView:
		return docStyle.Render(m.list.View())
	case addView, editView:
		return m.form.View()
	}
	return ""
}