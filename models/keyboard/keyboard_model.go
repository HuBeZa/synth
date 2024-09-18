package keyboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/synth/frequencies"
	"github.com/HuBeZa/synth/models"
	"github.com/HuBeZa/synth/models/base/options"
	"github.com/HuBeZa/synth/models/base/slider"
	"github.com/HuBeZa/synth/models/base/tremolo"
	"github.com/HuBeZa/synth/streamers"
)

const (
	// basic ascii keyboard - vertical bars & underscores only
	asciiKeyboard = "" +
		"_________________________________________\n" +
		"|  | | | |  |  | | | | | |  |  | | | |  |\n" +
		"|  | | | |  |  | | | | | |  |  | | | |  |\n" +
		"|  |w| |e|  |  |t| |y| |u|  |  |o| |p|  |\n" +
		"|  |_| |_|  |  |_| |_| |_|  |  |_| |_|  |\n" +
		"|   |   |   |   |   |   |   |   |   |   |\n" +
		"| a | s | d | f | g | h | j | k | l | ; |\n" +
		"|___|___|___|___|___|___|___|___|___|___|\n"

	// Unicode Box-drawing characters keyboard (see https://en.wikipedia.org/wiki/Box-drawing_characters)
	keyboardTall = "" +
		"╒══╤═╤═╤═╤══╤══╤═╤═╤═╤═╤═╤══╤══╤═╤═╤═╤══╕\n" +
		"│  │ │ │ │  │  │ │ │ │ │ │  │  │ │ │ │  │\n" +
		"│  │ │ │ │  │  │ │ │ │ │ │  │  │ │ │ │  │\n" +
		"│  │w│ │e│  │  │t│ │y│ │u│  │  │o│ │p│  │\n" +
		"│  └┬┘ └┬┘  │  └┬┘ └┬┘ └┬┘  │  └┬┘ └┬┘  │\n" +
		"│   │   │   │   │   │   │   │   │   │   │\n" +
		"│ a │ s │ d │ f │ g │ h │ j │ k │ l │ ; │\n" +
		"└───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘\n"

	keyboard = "" +
		"╒══╤═╤═╤═╤══╤══╤═╤═╤═╤═╤═╤══╤══╤═╤═╤═╤══╕\n" +
		"│  │w│ │e│  │  │t│ │y│ │u│  │  │o│ │p│  │\n" +
		"│  └┬┘ └┬┘  │  └┬┘ └┬┘ └┬┘  │  └┬┘ └┬┘  │\n" +
		"│ a │ s │ d │ f │ g │ h │ j │ k │ l │ ; │\n" +
		"└───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘\n"

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
	tremoloCtrlId     = "tremoloCtrl"
)

var (
	octaveToKeys = initOctaveToKeys()
	currKeyStyle = lipgloss.NewStyle().Reverse(true)
)

func initOctaveToKeys() map[int]map[string]frequencies.Frequency {
	octavesMap := make(map[int]map[string]frequencies.Frequency, 11)
	keys := []string{"a", "w", "s", "e", "d", "f", "t", "g", "y", "h", "u", "j", "k", "o", "l", "p", ";"}

	baseLow := frequencies.C0
	baseHigh := baseLow.ShiftTone(len(keys) - 1)

	for octaveId := -1; octaveId <= 9; octaveId++ {
		octavesMap[octaveId] = make(map[string]frequencies.Frequency, len(keys))
		octaveFrequencies := frequencies.GetRange(baseLow.ShiftOctave(octaveId), baseHigh.ShiftOctave(octaveId))
		for i, key := range keys {
			octavesMap[octaveId][key] = octaveFrequencies[i]
		}
	}
	return octavesMap
}

type model struct {
	waveformOptions options.Model[streamers.Waveform]
	octaveSlider    slider.Model
	panSlider       slider.Model
	gainSlider      slider.Model
	tremoloCtrl     tremolo.Model

	isSilenced    bool
	keyPressTimer timer.Model
	currKey       string
	currFreq      frequencies.Frequency
	streamer      streamers.DynamicStreamer
	zonePrefix    string
	zoneHandlers  map[string]func(model, tea.MouseMsg) (tea.Model, tea.Cmd)
}

func New(sr beep.SampleRate) models.StreamerModel {
	m := model{}
	m.waveformOptions = options.New(streamers.AllWaveforms())
	m.octaveSlider, _ = slider.New(-1, 9, 1, 3, 4)
	m.panSlider, _ = slider.New(-panSliderRatio, panSliderRatio, 1, 0, 0)
	m.gainSlider, _ = slider.New(0, gainSliderRatio*4, 1, 0, gainSliderRatio, gainSliderRatio*2, gainSliderRatio*3)
	m.tremoloCtrl = tremolo.New()
	m.zonePrefix = zone.NewPrefix()
	m.zoneHandlers = map[string]func(model, tea.MouseMsg) (tea.Model, tea.Cmd){
		m.zonePrefix + upButtonId:        upButtonHandler,
		m.zonePrefix + downButtonId:      downButtonHandler,
		m.zonePrefix + closeButtonId:     closeButtonHandler,
		m.zonePrefix + playStopButtonId:  playStopButtonHandler,
		m.zonePrefix + waveformOptionsId: waveformOptionsHandler,
		m.zonePrefix + octaveSliderId:    octaveSliderHandler,
		m.zonePrefix + panSliderId:       panSliderHandler,
		m.zonePrefix + gainSliderId:      gainSliderHandler,
		m.zonePrefix + tremoloCtrlId:     tremoloCtrlHandler,
	}

	m.streamer, _ = streamers.NewWaveformDynamicStreamer(sr, 0, 0, 1, streamers.Sine)
	m.streamer.Silence()

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
		for id, handler := range m.zoneHandlers {
			if zone.Get(id).InBounds(msg) {
				return handler(m, msg)
			}
		}
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "a", "w", "s", "e", "d", "f", "t", "g", "y", "h", "u", "j", "k", "o", "l", "p", ";":
			keyPressTimeout := 40
			if key != m.currKey {
				freq := octaveToKeys[m.octaveSlider.Value()][key]
				if !m.isSilenced {
					m.streamer.SetFrequency(freq.Frequency())
					m.streamer.Unsilence()
				}
				m.currKey = key
				m.currFreq = freq
				keyPressTimeout = 300
			}

			m.keyPressTimer = timer.NewWithInterval(time.Duration(keyPressTimeout)*time.Millisecond, 10*time.Millisecond)
			return m, m.keyPressTimer.Init()
		}
	case timer.TickMsg:
		if msg.Timeout {
			m.currKey = ""
			m.streamer.Silence()
		}

		var cmd tea.Cmd
		m.keyPressTimer, cmd = m.keyPressTimer.Update(msg)
		return m, cmd
		// case timer.TimeoutMsg:	// already handled on TickMsg
		// case timer.StartStopMsg:	// required only if Start/Stop/Toggle is called
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
		m.isSilenced = !m.isSilenced
	}
	return m, nil
}

func octaveSliderHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	sliderModel, cmd := m.octaveSlider.Update(msg)
	m.octaveSlider = sliderModel.(slider.Model)
	return m, cmd
}

func waveformOptionsHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	optionsModel, cmd := m.waveformOptions.Update(msg)
	m.waveformOptions = optionsModel.(options.Model[streamers.Waveform])
	m.streamer.SetWaveform(m.currentWaveform())
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

func tremoloCtrlHandler(m model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	tremoloModel, cmd := m.tremoloCtrl.Update(msg)
	m.tremoloCtrl = tremoloModel.(tremolo.Model)

	if m.tremoloCtrl.IsOn() {
		m.streamer.SetTremolo(m.tremoloCtrl.Duration(), m.tremoloCtrl.StartGain(), m.tremoloCtrl.EndGain(), m.tremoloCtrl.Pulsing())
	} else {
		m.streamer.SetTremoloOff()
	}
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderHeader(models.ColumnWidth),
		m.renderKeyboard(),
		m.renderWaveformOptions(),
		m.renderOctaveSlider(),
		m.renderPanSlider(),
		m.renderGainSlider(),
		m.renderTremoloCtrl())
}

func (m model) renderHeader(width int) string {
	widthLeft := width * 9 / 10
	widthRight := width - widthLeft

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderHeaderText(widthLeft),
		m.renderHeaderButtons(widthRight))
}

func (m model) renderHeaderText(width int) string {
	var header string
	if m.currFreq == nil {
		header = models.HeaderStyle().Render(m.currentWaveform().String())
	} else {
		header = models.HeaderStyle().Render(fmt.Sprintf("%v %v (%vHz)", m.currentWaveform(), m.currFreq.Name(), m.currFreq.Frequency()))
	}

	playStopButton := models.PlayButton()
	if m.isSilenced {
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

func (m model) renderKeyboard() string {
	if m.currKey != "" {
		return strings.Replace(keyboard, m.currKey, currKeyStyle.Render(m.currKey), 1)
	}
	return keyboard
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

func (m model) renderTremoloCtrl() string {
	id := m.zonePrefix + tremoloCtrlId
	return zone.Mark(id, m.tremoloCtrl.View())
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
