package models

import tea "github.com/charmbracelet/bubbletea"

type StreamerMsg struct{ Model StreamerModel }

type RemoveStreamerMsg StreamerMsg

func RemoveStreamerFunc(model StreamerModel) func() tea.Msg {
	return func() tea.Msg {
		return RemoveStreamerMsg{model}
	}
}

type StreamerUpMsg StreamerMsg

func StreamerUpFunc(model StreamerModel) func() tea.Msg {
	return func() tea.Msg {
		return StreamerUpMsg{model}
	}
}

type StreamerDownMsg StreamerMsg

func StreamerDownFunc(model StreamerModel) func() tea.Msg {
	return func() tea.Msg {
		return StreamerDownMsg{model}
	}
}
