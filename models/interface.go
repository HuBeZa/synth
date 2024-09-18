package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
)

type StreamerModel interface {
	tea.Model
	Equals(other tea.Model) bool
	Streamer() beep.Streamer
}
