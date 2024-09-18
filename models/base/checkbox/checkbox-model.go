package checkbox

import (
	"fmt"

	"github.com/HuBeZa/synth/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Model interface {
	tea.Model
	Value() bool
	SetValue(checked bool) Model
}

type model struct {
	label   string
	checked bool
	zoneId  string
}

func New(label string, checked bool) Model {
	return model{
		label:   label,
		checked: checked,
		zoneId:  zone.NewPrefix(),
	}
}

func (m model) Value() bool {
	return m.checked
}

func (m model) SetValue(checked bool) Model {
	m.checked = checked
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return m, nil
		}

		if zone.Get(m.zoneId).InBounds(msg) {
			m.checked = !m.checked
		}
	}
	return m, nil
}

func (m model) View() string {
	button := "◇"
	style := lipgloss.NewStyle()
	if m.checked {
		button = "◈"
		style = models.SelectedStyle() 
	}

	view := button
	if m.label != "" {
		view = fmt.Sprintf("%v %v", button, m.label)
	}

	view = style.Render(view)

	return zone.Mark(m.zoneId, view)
}
