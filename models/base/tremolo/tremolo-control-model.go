package tremolo

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/synth/models/base/checkbox"
	"github.com/HuBeZa/synth/models/base/slider"
)

const (
	gainSliderRatio = 10

	isOnCheckboxId    = "isOnCheckbox"
	speedSliderId     = "speedSlider"
	gainSliderId      = "gainSlider"
	pulsingCheckboxId = "pulsingCheckbox"
	reverseCheckboxId = "reverseCheckbox"
)

var (
	speedValues = []float64{1 / 1000.0, 1 / 500.0, 1 / 100.0, 1 / 75.0, 1 / 50.0, 1 / 25.0, 1 / 20.0, 1 / 10.0, 1 / 5.0, 1 / 3.0, 1 / 2.0, 1, 2, 3, 5}
	labelStyle  = lipgloss.NewStyle().Width(4)
	marginRight = lipgloss.NewStyle().MarginRight(2)
)

type Model interface {
	tea.Model
	IsOn() bool
	Duration() time.Duration
	StartGain() float64
	EndGain() float64
	Pulsing() bool
	Reverse() bool
}

type model struct {
	isOnCheckbox    checkbox.Model
	SpeedSlider     slider.Model
	gainSlider      slider.Model
	pulsingCheckbox checkbox.Model
	reverseCheckbox checkbox.Model
	zonePrefix      string
	zoneHandlers    map[string]func(model, tea.MouseMsg) (tea.Model, tea.Cmd)
}

func New() Model {
	m := model{}
	m.isOnCheckbox = checkbox.New("tremolo", false)
	m.SpeedSlider, _ = slider.New(0, len(speedValues)-1, 1, len(speedValues)/2, len(speedValues)/2)
	m.gainSlider, _ = slider.New(-gainSliderRatio, 2*gainSliderRatio, 1, gainSliderRatio/2, 0, gainSliderRatio)
	m.pulsingCheckbox = checkbox.New("pulsing", false)
	m.reverseCheckbox = checkbox.New("reverse", false)

	m.zonePrefix = zone.NewPrefix()
	m.zoneHandlers = map[string]func(model, tea.MouseMsg) (tea.Model, tea.Cmd){
		m.zonePrefix + isOnCheckboxId:    isOnCheckboxHandler,
		m.zonePrefix + speedSliderId:     speedSliderHandler,
		m.zonePrefix + gainSliderId:      gainSliderHandler,
		m.zonePrefix + pulsingCheckboxId: pulsingCheckboxHandler,
		m.zonePrefix + reverseCheckboxId: reverseCheckboxHandler,
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

func isOnCheckboxHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	checkboxModel, cmd := m.isOnCheckbox.Update(msg)
	m.isOnCheckbox = checkboxModel.(checkbox.Model)
	return m, cmd
}

func speedSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.SpeedSlider.Update(msg)
	m.SpeedSlider = sliderModel.(slider.Model)
	return m, cmd
}

func gainSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.gainSlider.Update(msg)
	m.gainSlider = sliderModel.(slider.Model)
	return m, cmd
}

func pulsingCheckboxHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	checkboxModel, cmd := m.pulsingCheckbox.Update(msg)
	m.pulsingCheckbox = checkboxModel.(checkbox.Model)
	return m, cmd
}

func reverseCheckboxHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	checkboxModel, cmd := m.reverseCheckbox.Update(msg)
	m.reverseCheckbox = checkboxModel.(checkbox.Model)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		marginRight.Render(m.renderIsOn()),
		lipgloss.JoinVertical(lipgloss.Left,
			m.renderGain(),
			m.renderSpeed(),
			lipgloss.JoinHorizontal(lipgloss.Top,
				marginRight.Render(m.renderPulsing()),
				m.renderReverse(),
			),
		),
	)
}

func (m model) renderIsOn() string {
	return zone.Mark(m.zonePrefix+isOnCheckboxId, m.isOnCheckbox.View())
}

func (m model) renderGain() string {
	label := labelStyle.Render("gain")
	slider := zone.Mark(m.zonePrefix+gainSliderId, m.gainSlider.View())
	val := m.Gain()
	return fmt.Sprintf("%v %v %v", label, slider, val)
}

func (m model) renderSpeed() string {
	label := labelStyle.Render("spd")
	slider := zone.Mark(m.zonePrefix+speedSliderId, m.SpeedSlider.View())
	val := fmt.Sprintf("%v/sec", m.timesPerSecond())
	return fmt.Sprintf("%v %v %v", label, slider, val)
}

func (m model) renderPulsing() string {
	return zone.Mark(m.zonePrefix+pulsingCheckboxId, m.pulsingCheckbox.View())
}

func (m model) renderReverse() string {
	return zone.Mark(m.zonePrefix+reverseCheckboxId, m.reverseCheckbox.View())
}

func (m model) IsOn() bool {
	return m.isOnCheckbox.Value()
}

func (m model) Duration() time.Duration {
	return time.Duration(m.durationSeconds() * float64(time.Second))
}

func (m model) durationSeconds() float64 {
	return speedValues[m.SpeedSlider.Value()]
}

func (m model) timesPerSecond() string {
	timesPerSec := 1.0 / m.durationSeconds()
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", timesPerSec), "0"), ".")
}

func (m model) Gain() float64 {
	return float64(m.gainSlider.Value()) / float64(gainSliderRatio)
}

func (m model) StartGain() float64 {
	if m.Reverse() {
		return 1
	}
	return m.Gain()
}

func (m model) EndGain() float64 {
	if m.Reverse() {
		return m.Gain()
	}
	return 1
}

func (m model) Pulsing() bool {
	return m.pulsingCheckbox.Value()
}

func (m model) Reverse() bool {
	return m.reverseCheckbox.Value()
}

// func main() {
// 	zone.NewGlobal()
// 	fmt.Println(New().View())
// }
