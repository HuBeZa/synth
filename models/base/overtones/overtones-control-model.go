package overtones

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/synth/models"
	"github.com/HuBeZa/synth/models/base/slider"
)

const (
	gainSliderRatio = 10

	countSliderId = "countSlider"
	gainSliderId  = "gainSlider"
)

var (
	marginRight = lipgloss.NewStyle().MarginRight(2)
)

type Model interface {
	tea.Model
	Count() int
	Gain() float64
}

type model struct {
	countSlider  slider.Model
	gainSlider   slider.Model
	zonePrefix   string
	zoneHandlers map[string]func(model, tea.MouseMsg) (tea.Model, tea.Cmd)
}

func New() Model {
	m := model{}
	m.countSlider, _ = slider.New(0, 4, 1, 0)
	m.gainSlider, _ = slider.New(0, 2*gainSliderRatio, 1, gainSliderRatio, gainSliderRatio)

	m.zonePrefix = zone.NewPrefix()
	m.zoneHandlers = map[string]func(model, tea.MouseMsg) (tea.Model, tea.Cmd){
		m.zonePrefix + countSliderId: countSliderHandler,
		m.zonePrefix + gainSliderId:  gainSliderHandler,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		for id, handler := range m.zoneHandlers {
			if zone.Get(id).InBounds(msg) {
				return handler(m, msg)
			}
		}
	}
	return m, nil
}

func countSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.countSlider.Update(msg)
	m.countSlider = sliderModel.(slider.Model)
	return m, cmd
}

func gainSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.gainSlider.Update(msg)
	m.gainSlider = sliderModel.(slider.Model)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		marginRight.Render(m.renderLabel()),
		marginRight.Render(m.renderCount()),
		m.renderGain(),
	)
}

func (m model) renderLabel() string {
	label := "overtones"
	if m.Count() > 0 && m.Gain() > 0 {
		label = models.SelectedStyle().Render(label)
	}
	return label
}

func (m model) renderCount() string {
	slider := zone.Mark(m.zonePrefix+countSliderId, m.countSlider.View())
	val := m.Count()
	return fmt.Sprintf("%v %v", slider, val)
}

func (m model) renderGain() string {
	label := "gain"
	slider := zone.Mark(m.zonePrefix+gainSliderId, m.gainSlider.View())
	val := m.Gain()
	return fmt.Sprintf("%v %v %v", label, slider, val)
}

func (m model) Count() int {
	return m.countSlider.Value()
}

func (m model) Gain() float64 {
	return float64(m.gainSlider.Value()) / float64(gainSliderRatio)
}
