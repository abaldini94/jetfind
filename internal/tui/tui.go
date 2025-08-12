package tui

import (
	"fmt"
	"jetfind/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() (string, error) {
	cfg := config.LoadOrDefault()
	ConfiguredStyles(
		cfg.Tui.HighlightedFile.Foreground,
		cfg.Tui.QueryBox.TextForeground,
		cfg.Tui.QueryBox.TextBackground,
		cfg.Tui.QueryBox.BorderForeground,
	)

	model := NewModel(cfg)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("TUI error: %w", err)
	}

	if m, ok := finalModel.(*Model); ok {
		return m.SelectedFile, nil
	}

	return "", nil
}
