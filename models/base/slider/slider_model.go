package slider

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

const (
	sliderEdge  = "¦"
	sliderMid   = "|"
	sliderStep  = "-"
	sliderValue = "■"
)

type Model interface {
	tea.Model
	Value() int
	SetValue(val int) (Model, error)
}

type model struct {
	min, max, step int
	value          int
	mids           []int
	zonePrefix     string
}

func New(min, max, step, val int, mids ...int) (Model, error) {
	if min >= max {
		return nil, fmt.Errorf("min should be smaller than max")
	}
	if (max-min)%step != 0 {
		return nil, fmt.Errorf("cannot divided by step")
	}

	m := model{min, max, step, val, mids, zone.NewPrefix()}
	if !m.validate(min) {
		return nil, fmt.Errorf("mid is not a legal value")
	}
	if !m.validate(val) {
		return nil, fmt.Errorf("val is not a legal value")
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Button != tea.MouseButtonLeft {
			return m, nil
		}

		for i := m.min; i <= m.max; i += m.step {
			if zone.Get(fmt.Sprintf("%v%v", m.zonePrefix, i)).InBounds(msg) {
				m.value = i
				return m, nil
			}
		}

	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder
	for i := m.min; i <= m.max; i += m.step {
		var val string
		if i == m.value {
			val = sliderValue
		} else if i == m.min || i == m.max {
			val = sliderEdge
		} else if slices.Contains(m.mids, i) {
			val = sliderMid
		} else {
			val = sliderStep
		}

		id := fmt.Sprintf("%v%v", m.zonePrefix, i)
		sb.WriteString(zone.Mark(id, val))
	}

	return sb.String()
}

func (m model) Value() int {
	return m.value
}

func (m model) SetValue(val int) (Model, error) {
	if !m.validate(val) {
		return m, fmt.Errorf("val is not a legal value")
	}
	m.value = val
	return m, nil
}

func (m model) validate(val int) bool {
	return val >= m.min && val <= m.max && (m.max-val)%m.step == 0
}
