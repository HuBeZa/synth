package oscillator

import (
	"fmt"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/synth/models"
	"github.com/HuBeZa/synth/models/base/options"
	"github.com/HuBeZa/synth/models/base/slider"
	"github.com/HuBeZa/synth/streamers"
	"github.com/HuBeZa/synth/streamers/frequencies"
)

const (
	panSliderRatio  = 10
	gainSliderRatio = 5

	// bubblezone ids:
	upButtonId        = "upButton"
	downButtonId      = "downButton"
	closeButtonId     = "closeButton"
	playStopButtonId  = "playStopButton"
	waveformOptionsId = "waveformOptions"
	octaveSliderId    = "octaveSlider"
	panSliderId       = "panSlider"
	gainSliderId      = "gainSlider"
	freqSliderId      = "freqSlider"
)

var (
	octaveToFrequencies = initOctaveToFrequencies()
	cFreqIndexes        = initCFreqIndexes()
)

func initOctaveToFrequencies() map[int][]frequencies.Frequency {
	baseLow := frequencies.C0()
	baseHigh := baseLow.ShiftOctave(4)

	octaveToFrequencies := make(map[int][]frequencies.Frequency, 7)
	for octaveId := -1; octaveId <= 6; octaveId++ {
		octaveToFrequencies[octaveId] = frequencies.GetRange(baseLow.ShiftOctave(octaveId), baseHigh.ShiftOctave(octaveId))
	}
	return octaveToFrequencies
}

func initCFreqIndexes() []int {
	isC := regexp.MustCompile("^C[0-9-]+$")
	freq := octaveToFrequencies[0]
	res := make([]int, 0)
	for i := 1; i < len(freq)-1; i++ {
		if isC.MatchString(freq[i].Name()) {
			res = append(res, i)
		}
	}
	return res
}

type model struct {
	waveformOptions options.Model[streamers.Waveform]
	octaveSlider    slider.Model
	panSlider       slider.Model
	gainSlider      slider.Model
	freqSlider      slider.Model
	streamer        streamers.DynamicStreamer
	zonePrefix      string
	zoneHandlers    models.ZoneHandlers[model]
}

func New(sr beep.SampleRate) models.StreamerModel {
	m := model{}
	m.waveformOptions = options.New(streamers.AllWaveforms(), false)
	m.octaveSlider, _ = slider.New(-1, 6, 1, 3, 3)
	m.panSlider, _ = slider.New(-panSliderRatio, panSliderRatio, 1, 0, 0)
	m.gainSlider, _ = slider.New(0, gainSliderRatio*4, 1, gainSliderRatio, gainSliderRatio, gainSliderRatio*2, gainSliderRatio*3)
	m.freqSlider, _ = slider.New(0, len(m.currentOctave())-1, 1, 0, cFreqIndexes...)
	m.zonePrefix = zone.NewPrefix()
	m.zoneHandlers = models.ZoneHandlers[model]{
		m.zonePrefix + upButtonId:        upButtonHandler,
		m.zonePrefix + downButtonId:      downButtonHandler,
		m.zonePrefix + closeButtonId:     closeButtonHandler,
		m.zonePrefix + playStopButtonId:  playStopButtonHandler,
		m.zonePrefix + waveformOptionsId: waveformOptionsHandler,
		m.zonePrefix + octaveSliderId:    octaveSliderHandler,
		m.zonePrefix + panSliderId:       panSliderHandler,
		m.zonePrefix + gainSliderId:      gainSliderHandler,
		m.zonePrefix + freqSliderId:      freqSliderHandler,
	}

	var err error
	m.streamer, err = streamers.NewWaveformDynamicStreamer(sr, m.currentFrequency(), m.currentPan(), m.currentGain(), m.currentWaveform())
	if err != nil {
		panic(err)
	}

	return m
}

func (m model) Equals(other tea.Model) bool {
	if other, ok := other.(model); ok {
		return m.zonePrefix == other.zonePrefix
	}
	return false
}

func (m model) Streamer() beep.Streamer {
	return m.streamer
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

func upButtonHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		return m, models.StreamerUpFunc(m)
	}
	return m, nil
}

func downButtonHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		return m, models.StreamerDownFunc(m)
	}
	return m, nil
}

func closeButtonHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		return m, models.RemoveStreamerFunc(m)
	}
	return m, nil
}

func playStopButtonHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		m.streamer.ToggleSilence()
	}
	return m, nil
}

func waveformOptionsHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	optionsModel, cmd := m.waveformOptions.Update(msg)
	m.waveformOptions = optionsModel.(options.Model[streamers.Waveform])
	m.streamer.SetWaveform(m.currentWaveform())
	return m, cmd
}

func octaveSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.octaveSlider.Update(msg)
	m.octaveSlider = sliderModel.(slider.Model)
	m.streamer.SetFrequency(m.currentFrequency())
	return m, cmd
}

func panSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.panSlider.Update(msg)
	m.panSlider = sliderModel.(slider.Model)
	m.streamer.SetPan(m.currentPan())
	return m, cmd
}

func gainSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.gainSlider.Update(msg)
	m.gainSlider = sliderModel.(slider.Model)
	m.streamer.SetGain(m.currentGain())
	return m, cmd
}

func freqSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.freqSlider.Update(msg)
	m.freqSlider = sliderModel.(slider.Model)
	m.streamer.SetFrequency(m.currentFrequency())
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderHeader(models.ColumnWidth),
		m.renderOctaveSlider(),
		m.renderFreqSlider(),
		m.renderWaveformOptions(),
		m.renderPanSlider(),
		m.renderGainSlider())
}

func (m model) renderHeader(width int) string {
	widthLeft := width * 9 / 10
	widthRight := width - widthLeft

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderHeaderText(widthLeft),
		m.renderHeaderButtons(widthRight))
}

func (m model) renderHeaderText(width int) string {
	header := models.HeaderStyle().Render(fmt.Sprintf("%v %v (%vHz)", m.currentWaveform(), m.currentFrequency().Name(), m.currentFrequency().Frequency()))

	playStopButton := models.PlayButton()
	if m.streamer.IsSilenced() {
		playStopButton = models.StopButton()
	}

	id := m.zonePrefix + playStopButtonId
	view := zone.Mark(id, fmt.Sprintf("%v %v", playStopButton, header))
	return lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Left).Render(view)
}

func (m model) renderHeaderButtons(width int) string {
	upButton := zone.Mark(m.zonePrefix+upButtonId, models.UpButton())
	downButton := zone.Mark(m.zonePrefix+downButtonId, models.DownButton())
	closeButton := zone.Mark(m.zonePrefix+closeButtonId, models.CloseButton())
	view := lipgloss.JoinHorizontal(lipgloss.Top, upButton, downButton, closeButton)

	return lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Right).MarginRight(1).Render(view)
}

func (m model) renderWaveformOptions() string {
	id := m.zonePrefix + waveformOptionsId
	return zone.Mark(id, m.waveformOptions.View())
}

func (m model) renderOctaveSlider() string {
	id := m.zonePrefix + octaveSliderId
	return models.LabelStyle().Render("octave") + zone.Mark(id, m.octaveSlider.View()) + fmt.Sprintf(" %v", m.octaveSlider.Value())
}

func (m model) renderPanSlider() string {
	id := m.zonePrefix + panSliderId
	return models.LabelStyle().Render("pan") + zone.Mark(id, m.panSlider.View()) + fmt.Sprintf(" %v", m.streamer.Pan())
}

func (m model) renderGainSlider() string {
	id := m.zonePrefix + gainSliderId
	return models.LabelStyle().Render("gain") + zone.Mark(id, m.gainSlider.View()) + fmt.Sprintf(" %v", m.streamer.Gain())
}

func (m model) renderFreqSlider() string {
	id := m.zonePrefix + freqSliderId
	return models.LabelStyle().Render("freq") + zone.Mark(id, m.freqSlider.View())
}

func (m model) currentWaveform() streamers.Waveform {
	return m.waveformOptions.Value()
}

func (m model) currentPan() float64 {
	return float64(m.panSlider.Value()) / float64(panSliderRatio)
}

func (m model) currentGain() float64 {
	return float64(m.gainSlider.Value()) / float64(gainSliderRatio)
}

func (m model) currentOctave() []frequencies.Frequency {
	return octaveToFrequencies[m.octaveSlider.Value()]
}

func (m model) currentFrequency() frequencies.Frequency {
	return m.currentOctave()[m.freqSlider.Value()]
}
