package service

import "github.com/charmbracelet/lipgloss"

const margin = 2

var (
	_headingStyle   = lipgloss.NewStyle().Bold(true)
	_pathStyle      = lipgloss.NewStyle().MarginLeft(margin)
	_errorLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)
