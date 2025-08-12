package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	HighlightedFileStyle lipgloss.Style
	QueryBoxStyle        lipgloss.Style
	SeparatorStyle       lipgloss.Style
	StatusStyle          lipgloss.Style
)

func DefaultStyles() {

	HighlightedFileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	StatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Italic(true)

	QueryBoxStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F9FAFB")).
		Background(lipgloss.Color("#374151")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6B7280"))

	SeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563")).
		Bold(true)
}

func ConfiguredStyles(hfForeground, qbTextForeground, qbTextBackground, qbBorderForeground string) {
	HighlightedFileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(hfForeground)).
		Bold(true).
		Padding(0, 1)

	StatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Italic(true)

	QueryBoxStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(qbTextForeground)).
		Background(lipgloss.Color(qbTextBackground)).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(qbBorderForeground))

	SeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563")).
		Bold(true)
}
