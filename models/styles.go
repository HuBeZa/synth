package models

import "github.com/charmbracelet/lipgloss"

const ColumnWidth = 57

var (
	headerStyle = lipgloss.NewStyle().Underline(true)
	labelStyle = lipgloss.NewStyle().Width(7)
	selectedStyle = ForegroundColor("#87afff")
	playButton  = ForegroundColor("#008000").Render("►")
	stopButton  = ForegroundColor("#800000").Render("■")
	upButton    = ForegroundColor("#005fff").MarginLeft(1).Render("⮝")
	downButton  = ForegroundColor("#005fff").MarginLeft(1).Render("⮟")
	closeButton = ForegroundColor("#DF0000").MarginLeft(1).Render("✘")
)

func HeaderStyle() lipgloss.Style {
	return headerStyle
}

func LabelStyle() lipgloss.Style {
	return labelStyle
}

func SelectedStyle() lipgloss.Style {
	return selectedStyle
}

func PlayButton() string {
	return playButton
}

func StopButton() string {
	return stopButton
}

func UpButton() string {
	return upButton
}

func DownButton() string {
	return downButton
}

func CloseButton() string {
	return closeButton
}

func ForegroundColor(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}
