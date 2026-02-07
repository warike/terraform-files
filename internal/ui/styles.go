package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	CheckboxStyle  = lipgloss.NewStyle().PaddingLeft(2)
	CheckedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	UncheckedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	HelpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(2)
	SpinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	SuccessStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
)
