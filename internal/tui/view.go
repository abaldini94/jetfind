package tui

import (
	"fmt"
	"jetfind/internal/scanengine"
	"strings"
)

func (m *Model) View() string {
	if m.scanErr != nil {
		return fmt.Sprintf("Error: %v\n", m.scanErr)
	}

	if m.userQuery == "" {
		m.filteredPaths = m.scannedPaths
	} else {
		var scanFilter scanengine.ScanFilter
		if m.filterRequested {
			switch m.cfg.Filter.Type {
			case "null":
				scanFilter = scanengine.NoFilter{Pattern: m.userQuery}
			case "contains":
				scanFilter = scanengine.ContainsFilter{Pattern: m.userQuery}
			case "fuzzy":
				scanFilter = scanengine.FuzzyFilter{
					Pattern:    m.userQuery,
					Algo:       m.cfg.Filter.Algo,
					Threashold: m.cfg.Filter.Threashold,
				}

			}
			m.filteredPaths = scanengine.FilterEngine(m.scannedPaths, scanFilter)
			m.filterRequested = false
		}
	}

	var b strings.Builder

	queryText := "Search: " + m.userQuery
	if m.userQuery == "" {
		queryText = "Search: (type to search...)"
	}
	queryBox := QueryBoxStyle.Render(queryText)
	b.WriteString(queryBox + "\n")
	separator := SeparatorStyle.Render("────────────────────────────────────────")
	b.WriteString(separator + "\n")

	lastIdx := m.offset + m.height - 4
	if lastIdx > len(m.filteredPaths) {
		lastIdx = len(m.filteredPaths)
	}

	for i := m.offset; i < lastIdx; i++ {
		line := ""
		path := m.filteredPaths[i]
		if i == m.cursor {
			line = HighlightedFileStyle.Render(fmt.Sprintf("❯ %.1f  %s", path.Score, path.Path))
		} else {
			line = fmt.Sprintf("  %.1f  %s", path.Score, path.Path)
		}
		b.WriteString(line + "\n")
	}

	var status string
	if m.scanDone {
		status = StatusStyle.Render(fmt.Sprintf("--- Scan Completed (%d); Filtered (%d) ---", len(m.scannedPaths), len(m.filteredPaths)))
	} else {
		status = StatusStyle.Render(fmt.Sprintf("--- Scanning (%d); Filtered (%d) ---", len(m.scannedPaths), len(m.filteredPaths)))
	}

	b.WriteString(status)
	return b.String()
}
