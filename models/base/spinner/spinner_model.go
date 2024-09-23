package spinner

import (
	"fmt"

	"github.com/HuBeZa/synth/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	// bubblezone ids:
	prevButtonId = "prevButton"
	nextButtonId = "nextButton"
)

var valStyle = lipgloss.NewStyle().Width(7).AlignHorizontal(lipgloss.Center)

type Model[T any] interface {
	tea.Model
	Value() T
	SetValue(i int) Model[T]
}

type model[T any] struct {
	cursor       int
	options      []T
	allowCyclic  bool
	zonePrefix   string
	zoneHandlers models.ZoneHandlers[model[T]]
}

func New[T any](options []T, allowCyclic bool) Model[T] {
	m := model[T]{
		options:     options,
		allowCyclic: allowCyclic,
		zonePrefix:  zone.NewPrefix(),
	}

	m.zoneHandlers = models.ZoneHandlers[model[T]]{
		m.zonePrefix + prevButtonId: m.prevButtonHandler,
		m.zonePrefix + nextButtonId: m.nextButtonHandler,
	}

	return m
}

func (m model[T]) Value() T {
	return m.options[m.cursor]
}

func (m model[T]) SetValue(i int) Model[T] {
	if i >= 0 && i < len(m.options) {
		m.cursor = i
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

		return m.zoneHandlers.Handle(m, msg)
	}
	return m, nil
}

func (model[T]) prevButtonHandler(m model[T], _ tea.MouseMsg) (tea.Model, tea.Cmd) {
	m.cursor--
	if m.cursor < 0 {
		if m.allowCyclic {
			m.cursor = len(m.options) - 1
		} else {
			m.cursor = 0
		}
	}

	return m, nil
}

func (model[T]) nextButtonHandler(m model[T], _ tea.MouseMsg) (tea.Model, tea.Cmd) {
	m.cursor++
	if m.cursor >= len(m.options) {
		if m.allowCyclic {
			m.cursor = 0
		} else {
			m.cursor = len(m.options) - 1
		}
	}

	return m, nil
}

// ◁   val   ▷
func (m model[T]) View() string {
	prevButton := zone.Mark(m.zonePrefix+prevButtonId, "◁ ")
	val := valStyle.Render(fmt.Sprint(m.Value()))
	nextButton := zone.Mark(m.zonePrefix+nextButtonId, " ▷")

	return lipgloss.JoinHorizontal(lipgloss.Top, prevButton, val, nextButton)
}
