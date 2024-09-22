package chords

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/synth/models"
	"github.com/HuBeZa/synth/models/base/options"
	"github.com/HuBeZa/synth/models/base/slider"
	"github.com/HuBeZa/synth/streamers/chords"
)

const (
	// bubblezone ids:
	chordsOptionsId = "chordsOptions"
	delaySliderId   = "delaySlider"
)

var (
	delayValues = []time.Duration{0, 50 * time.Millisecond, 100 * time.Millisecond, 150 * time.Millisecond, 200 * time.Millisecond, 250 * time.Millisecond, 
		300 * time.Millisecond, 350 * time.Millisecond, 400 * time.Millisecond, 450 * time.Millisecond, 500 * time.Millisecond}
)

type Model interface {
	tea.Model
	Chord() chords.ChordType
	ArpeggioDelay() time.Duration
}

type model struct {
	chordsOptions options.Model[chords.ChordType]
	delaySlider   slider.Model
	zonePrefix    string
	zoneHandlers  models.ZoneHandlers[model]
}

func New() Model {
	m := model{}
	m.chordsOptions = options.New(chords.ChordTypes(), true)
	m.delaySlider, _ = slider.New(0, len(delayValues)-1, 1, 0, len(delayValues)/2)
	m.zonePrefix = zone.NewPrefix()
	m.zoneHandlers = models.ZoneHandlers[model]{
		m.zonePrefix + chordsOptionsId: chordsOptionsHandler,
		m.zonePrefix + delaySliderId:   delaySliderHandler,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		return m.zoneHandlers.Handle(m, msg)
	}
	return m, nil
}

func chordsOptionsHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	optionsModel, cmd := m.chordsOptions.Update(msg)
	m.chordsOptions = optionsModel.(options.Model[chords.ChordType])
	return m, cmd
}

func delaySliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.delaySlider.Update(msg)
	m.delaySlider = sliderModel.(slider.Model)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderLabel(),
		lipgloss.JoinVertical(lipgloss.Left,
			m.renderChordsOptions(),
			m.renderDelaySlider()))
}

func (m model) renderLabel() string {
	label := models.LabelStyle().Render("chord")
	if m.chordsOptions.Value() != nil {
		label = models.SelectedStyle().Render(label)
	}
	return label
}

func (m model) renderChordsOptions() string {
	id := m.zonePrefix + chordsOptionsId
	return zone.Mark(id, m.chordsOptions.View())
}

func (m model) renderDelaySlider() string {
	label := "arpeggio"
	slider := zone.Mark(m.zonePrefix+delaySliderId, m.delaySlider.View())
	val := m.ArpeggioDelay()
	return fmt.Sprintf("%v %v %v", label, slider, val)
}

func (m model) Chord() chords.ChordType {
	return m.chordsOptions.Value()
}

func (m model) ArpeggioDelay() time.Duration {
	return delayValues[m.delaySlider.Value()]
}
