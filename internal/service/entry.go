package service

import (
	"errors"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/skewb1k/upfile/pkg/sha256"
)

type Entry struct {
	Path   string
	Status string
	Err    error
}

func getEntry(path string, upstreamHash sha256.SHA256) *Entry {
	entry := &Entry{
		Path:   path,
		Status: lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Up-to-date"),
		Err:    nil,
	}

	existing, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			entry.Status = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Deleted")
		} else {
			entry.Err = errors.Unwrap(err)
		}
	} else if !upstreamHash.EqualBytes(existing) {
		entry.Status = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Modified")
	}

	return entry
}
