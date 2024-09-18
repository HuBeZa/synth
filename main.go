package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	zone "github.com/lrstanley/bubblezone"
	"golang.org/x/term"

	"github.com/HuBeZa/synth/models"
	"github.com/HuBeZa/synth/models/keyboard"
	"github.com/HuBeZa/synth/models/oscillator"
)

const defaultSampleRate = beep.SampleRate(48000)

var (
	streamerModelStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	helpStyle          = models.ForegroundColor("#626262").MarginTop(1).MarginLeft(1)
)

type mainModel struct {
	streamers []tea.Model
}

func newModel() tea.Model {
	m := mainModel{}.addNewKeyboard()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			speaker.Close()
			return m, tea.Quit
		case "ctrl+k":
			return m.addNewKeyboard(), nil
		case "ctrl+o":
			return m.addNewOscillator(), nil
		default:
			return m.updateStreamers(msg)
		}
	case timer.TickMsg:
		return m.updateStreamers(msg)
	case models.StreamerUpMsg:
		return m.moveStreamer(msg.Model, -1)
	case models.StreamerDownMsg:
		return m.moveStreamer(msg.Model, +1)
	case models.RemoveStreamerMsg:
		return m.removeStreamer(msg.Model)
	case tea.MouseMsg:
		for i := range m.streamers {
			if zone.Get(getStreamerZoneId(i)).InBounds(msg) {
				var cmd tea.Cmd
				m.streamers[i], cmd = m.streamers[i].Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m mainModel) updateStreamers(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	for i := range m.streamers {
		var cmd tea.Cmd
		m.streamers[i], cmd = m.streamers[i].Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	return zone.Scan(
		lipgloss.JoinVertical(lipgloss.Left, m.renderStreamers(), m.renderHelp()),
	)
}

func (m mainModel) renderStreamers() string {
	if len(m.streamers) == 0 {
		return ""
	}

	// the -=2 is to account for help text
	screenWidth, screenHeight := getTerminalSize()
	screenHeight -= 2
	maxColumns := max(screenWidth/models.ColumnWidth, 1)

	columns := make([][]string, 1, maxColumns)
	currCol := 0
	colHeight := 0
	for i, streamer := range m.streamers {
		id := getStreamerZoneId(i)
		view := streamerModelStyle.Render(streamer.View())

		if screenHeight > 0 && currCol < maxColumns-1 {
			lines := strings.Count(view, "\n") + 1
			colHeight += lines
			if colHeight > screenHeight {
				currCol++
				colHeight = lines
				columns = append(columns, make([]string, 0, 1))
			}
		}

		columns[currCol] = append(columns[currCol], zone.Mark(id, view))
	}

	views := make([]string, len(columns))
	for i, col := range columns {
		views[i] = lipgloss.JoinVertical(lipgloss.Left, col...)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}

func (m mainModel) renderHelp() string {
	return helpStyle.Render("ctrl+k: add keyboard • ctrl+o: add oscillator • ctrl-q: exit")
}

func (m mainModel) addNewKeyboard() mainModel {
	model := keyboard.New(defaultSampleRate)
	speaker.Play(model.Streamer())
	m.streamers = append(m.streamers, model)
	return m
}

func (m mainModel) addNewOscillator() mainModel {
	model := oscillator.New(defaultSampleRate)
	speaker.Play(model.Streamer())
	m.streamers = append(m.streamers, model)
	return m
}

func (m mainModel) moveStreamer(model models.StreamerModel, diff int) (tea.Model, tea.Cmd) {
	i := m.indexOf(model)
	if i != -1 && i+diff >= 0 && i+diff < len(m.streamers) {
		m.streamers[i], m.streamers[i+diff] = m.streamers[i+diff], m.streamers[i]
	}
	return m, nil
}

func (m mainModel) removeStreamer(model models.StreamerModel) (tea.Model, tea.Cmd) {
	if i := m.indexOf(model); i != -1 {
		m.streamers = slices.Delete(m.streamers, i, i+1)
		m.restartSpeaker()
	}
	return m, nil
}

func (m mainModel) indexOf(model models.StreamerModel) int {
	for i, streamerModel := range m.streamers {
		if model.Equals(streamerModel) {
			return i
		}
	}
	return -1
}

func (m mainModel) restartSpeaker() {
	streamers := make([]beep.Streamer, len(m.streamers))
	for i, streamer := range m.streamers {
		streamers[i] = streamer.(models.StreamerModel).Streamer()
	}
	speaker.Clear()
	speaker.Play(streamers...)
}

func getStreamerZoneId(i int) string {
	return fmt.Sprintf("streamer_%v", i)
}

func getTerminalSize() (width, height int) {
	width, height, _ = term.GetSize(int(os.Stdout.Fd()))
	return
}

func main() {
	zone.NewGlobal()
	speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/10))

	m := newModel()
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
