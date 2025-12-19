package ui

import (
	"tv-shows-manager/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
)

type formModel struct {
	inputs  []textinput.Model
	focused int
}

func newFormModel(show models.Show) formModel {
	inputs := make([]textinput.Model, 7)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Show name"
	inputs[0].SetValue(show.Show)
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Season"
	inputs[1].SetValue(show.Season)

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Year watched"
	inputs[2].SetValue(show.YearWatched)

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Source"
	inputs[3].SetValue(show.Source)

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "TMDB ID"
	inputs[4].SetValue(show.TmdbID)

	inputs[5] = textinput.New()
	inputs[5].Placeholder = "Kind"
	inputs[5].SetValue(show.Kind)

	inputs[6] = textinput.New()
	inputs[6].Placeholder = "Season rating"
	inputs[6].SetValue(show.SeasonRating)

	return formModel{inputs: inputs, focused: 0}
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.inputs[m.focused].Blur()
			m.focused = (m.focused + 1) % len(m.inputs)
			return m, m.inputs[m.focused].Focus()
		case "shift+tab", "up":
			m.inputs[m.focused].Blur()
			m.focused--
			if m.focused < 0 {
				m.focused = len(m.inputs) - 1
			}
			return m, m.inputs[m.focused].Focus()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *formModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m formModel) View() string {
	var s string
	s += focusedStyle.Render("TV Show Form") + "\n\n"

	labels := []string{"Show:", "Season:", "Year:", "Source:", "TMDB ID:", "Kind:", "Rating:"}

	for i := range m.inputs {
		label := labels[i]
		if i == m.focused {
			s += focusedStyle.Render(label) + "\n"
		} else {
			s += blurredStyle.Render(label) + "\n"
		}
		s += m.inputs[i].View() + "\n\n"
	}

	s += blurredStyle.Render("Ctrl+S: Save | Esc: Cancel") + "\n"

	return lipgloss.NewStyle().Margin(1, 2).Render(s)
}

func (m formModel) getShow() models.Show {
	return models.Show{
		Show:         m.inputs[0].Value(),
		Season:       m.inputs[1].Value(),
		YearWatched:  m.inputs[2].Value(),
		Source:       m.inputs[3].Value(),
		TmdbID:       m.inputs[4].Value(),
		Kind:         m.inputs[5].Value(),
		SeasonRating: m.inputs[6].Value(),
	}
}
