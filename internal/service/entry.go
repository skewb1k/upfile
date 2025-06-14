package service

import "github.com/charmbracelet/lipgloss"

type EntryStatus int

const (
	EntryStatusModified EntryStatus = iota
	EntryStatusUpToDate
	EntryStatusDeleted
)

type Entry struct {
	Path   string
	Status EntryStatus
	Err    error
}

func statusAsString(status EntryStatus) string {
	switch status {
	case EntryStatusModified:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Modified")
	case EntryStatusUpToDate:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Up-to-date")
	case EntryStatusDeleted:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Deleted")
	default:
		panic("UNEXPECTED")
	}
}
