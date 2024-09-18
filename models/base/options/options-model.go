package options

import (
	"fmt"

	"github.com/HuBeZa/synth/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Model[T comparable] interface {
	tea.Model
	Value() T
	SetValue(v T) Model[T]
}

type model[T comparable] struct {
	cursor     int
	options    []T
	zonePrefix string
}

func New[T comparable](options []T) Model[T] {
	return model[T]{
		options: options,
		zonePrefix: zone.NewPrefix(),
	}
}

func (m model[T]) Value() T {
	return m.options[m.cursor]
}

func (m model[T]) SetValue(v T) Model[T] {
	for i := range m.options {
		if m.options[i] == v {
			m.cursor = i
		}
	}
	return m
}

func (m model[T]) Init() tea.Cmd {
	return nil
}

func (m model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return m, nil
		}

		for i := range m.options {
			if zone.Get(m.getZoneId(i)).InBounds(msg) {
				m.cursor = i
			}
		}
	}
	return m, nil
}

func (m model[T]) View() string {
	optionsView := make([]string, len(m.options))
	for i := range m.options {
		button := "◇"
		style := lipgloss.NewStyle()
		if i == m.cursor {
			button = "◈"
			style = models.SelectedStyle() 
		}
		id := m.getZoneId(i)
		view := style.Render(fmt.Sprintf("%v %v ", button, m.options[i]))
		optionsView[i] = zone.Mark(id, view)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, optionsView...)
}

func (m model[T]) getZoneId(i int) string {
	return fmt.Sprintf("%v%v", m.zonePrefix, i)
}
