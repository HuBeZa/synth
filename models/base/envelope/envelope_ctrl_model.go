package envelope

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/synth/models"
	"github.com/HuBeZa/synth/models/base/options"
	"github.com/HuBeZa/synth/models/base/spinner"
	"github.com/HuBeZa/synth/streamers/composers"
)

const (
	// bubblezone ids:
	attackSpinnerId      = "attackSpinner"
	decaySpinnerId       = "decaySpinner"
	sustainSpinnerId     = "sustainSpinner"
	releaseSpinnerId     = "releaseSpinner"
	attackTypeOptionsId  = "attackTypeOptions"
	decayTypeOptionsId   = "decayTypeOptions"
	releaseTypeOptionsId = "releaseTypeOptions"
)

var (
	durationValues = []time.Duration{0, 10 * time.Millisecond, 25 * time.Millisecond, 50 * time.Millisecond, 100 * time.Millisecond, 150 * time.Millisecond,
		200 * time.Millisecond, 250 * time.Millisecond, 300 * time.Millisecond, 350 * time.Millisecond, 400 * time.Millisecond, 450 * time.Millisecond,
		500 * time.Millisecond, 600 * time.Millisecond, 700 * time.Millisecond, 800 * time.Millisecond, 900 * time.Millisecond, 1 * time.Second}

	gainValues = []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}

	marginLeftStyle = lipgloss.NewStyle().MarginLeft(4)
)

type Model interface {
	tea.Model
	ADSR() (time.Duration, composers.TransitionType, time.Duration, composers.TransitionType, float64, time.Duration, composers.TransitionType)
	Attack() time.Duration
	AttackType() composers.TransitionType
	Decay() time.Duration
	DecayType() composers.TransitionType
	Sustain() float64
	Release() time.Duration
	ReleaseType() composers.TransitionType
}

type model struct {
	attackSpinner      spinner.Model[time.Duration]
	decaySpinner       spinner.Model[time.Duration]
	sustainSpinner     spinner.Model[float64]
	releaseSpinner     spinner.Model[time.Duration]
	attackTypeOptions  options.Model[composers.TransitionType]
	decayTypeOptions   options.Model[composers.TransitionType]
	releaseTypeOptions options.Model[composers.TransitionType]
	zonePrefix         string
	zoneHandlers       models.ZoneHandlers[model]
}

func New() Model {
	m := model{}
	m.attackSpinner = spinner.New(durationValues, true)
	m.decaySpinner = spinner.New(durationValues, true)
	m.sustainSpinner = spinner.New(gainValues, true).SetValue(len(gainValues) - 1)
	m.releaseSpinner = spinner.New(durationValues, true)
	m.attackTypeOptions = options.New(composers.TransitionTypes(), false)
	m.decayTypeOptions = options.New(composers.TransitionTypes(), false)
	m.releaseTypeOptions = options.New(composers.TransitionTypes(), false)
	m.zonePrefix = zone.NewPrefix()
	m.zoneHandlers = models.ZoneHandlers[model]{
		m.zonePrefix + attackSpinnerId:      attackSpinnerHandler,
		m.zonePrefix + decaySpinnerId:       decaySpinnerHandler,
		m.zonePrefix + sustainSpinnerId:     sustainSpinnerHandler,
		m.zonePrefix + releaseSpinnerId:     releaseSpinnerHandler,
		m.zonePrefix + attackTypeOptionsId:  attackTypeOptionsHandler,
		m.zonePrefix + decayTypeOptionsId:   decayTypeOptionsHandler,
		m.zonePrefix + releaseTypeOptionsId: releaseTypeOptionsHandler,
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

func attackSpinnerHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	spinnerModel, cmd := m.attackSpinner.Update(msg)
	m.attackSpinner = spinnerModel.(spinner.Model[time.Duration])
	return m, cmd
}

func decaySpinnerHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	spinnerModel, cmd := m.decaySpinner.Update(msg)
	m.decaySpinner = spinnerModel.(spinner.Model[time.Duration])
	return m, cmd
}

func sustainSpinnerHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	spinnerModel, cmd := m.sustainSpinner.Update(msg)
	m.sustainSpinner = spinnerModel.(spinner.Model[float64])
	return m, cmd
}

func releaseSpinnerHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	spinnerModel, cmd := m.releaseSpinner.Update(msg)
	m.releaseSpinner = spinnerModel.(spinner.Model[time.Duration])
	return m, cmd
}

func attackTypeOptionsHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	optionsModel, cmd := m.attackTypeOptions.Update(msg)
	m.attackTypeOptions = optionsModel.(options.Model[composers.TransitionType])
	return m, cmd
}

func decayTypeOptionsHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	optionsModel, cmd := m.decayTypeOptions.Update(msg)
	m.decayTypeOptions = optionsModel.(options.Model[composers.TransitionType])
	return m, cmd
}

func releaseTypeOptionsHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	optionsModel, cmd := m.releaseTypeOptions.Update(msg)
	m.releaseTypeOptions = optionsModel.(options.Model[composers.TransitionType])
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderLabel(),
		lipgloss.JoinVertical(lipgloss.Left,
			m.renderAttack(),
			m.renderDecay(),
			m.renderSustain(),
			m.renderRelease(),
		),
	)
}

func (m model) renderLabel() string {
	return models.LabelStyle().Render("env")
}

func (m model) renderAttack() string {
	label := "A"
	spinner := zone.Mark(m.zonePrefix+attackSpinnerId, m.attackSpinner.View())
	transitionType := zone.Mark(m.zonePrefix+attackTypeOptionsId, marginLeftStyle.Render(m.attackTypeOptions.View()))
	return fmt.Sprintf("%v %v %v", label, spinner, transitionType)
}

func (m model) renderDecay() string {
	label := "D"
	spinner := zone.Mark(m.zonePrefix+decaySpinnerId, m.decaySpinner.View())
	transitionType := zone.Mark(m.zonePrefix+decayTypeOptionsId, marginLeftStyle.Render(m.decayTypeOptions.View()))
	return fmt.Sprintf("%v %v %v", label, spinner, transitionType)
}

func (m model) renderSustain() string {
	label := "S"
	spinner := zone.Mark(m.zonePrefix+sustainSpinnerId, m.sustainSpinner.View())
	return fmt.Sprintf("%v %v", label, spinner)
}

func (m model) renderRelease() string {
	label := "R"
	spinner := zone.Mark(m.zonePrefix+releaseSpinnerId, m.releaseSpinner.View())
	transitionType := zone.Mark(m.zonePrefix+releaseTypeOptionsId, marginLeftStyle.Render(m.releaseTypeOptions.View()))
	return fmt.Sprintf("%v %v %v", label, spinner, transitionType)
}

func (m model) ADSR() (time.Duration, composers.TransitionType, time.Duration, composers.TransitionType, float64, time.Duration, composers.TransitionType) {
	return m.Attack(), m.AttackType(), m.Decay(), m.DecayType(), m.Sustain(), m.Release(), m.ReleaseType()
}

func (m model) Attack() time.Duration {
	return m.attackSpinner.Value()
}

func (m model) Decay() time.Duration {
	return m.decaySpinner.Value()
}

func (m model) Release() time.Duration {
	return m.releaseSpinner.Value()
}

func (m model) Sustain() float64 {
	return m.sustainSpinner.Value()
}

func (m model) AttackType() composers.TransitionType {
	return m.attackTypeOptions.Value()
}

func (m model) DecayType() composers.TransitionType {
	return m.decayTypeOptions.Value()
}

func (m model) ReleaseType() composers.TransitionType {
	return m.releaseTypeOptions.Value()
}
