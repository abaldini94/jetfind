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

	m.applyFiltering()

	var b strings.Builder

	m.renderQueryBox(&b)
	m.renderSeparator(&b)
	m.renderPathList(&b)
	m.renderStatus(&b)

	return b.String()
}

func (m *Model) applyFiltering() {
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
}

func (m *Model) renderQueryBox(b *strings.Builder) {
	queryText := "Search: " + m.userQuery
	if m.userQuery == "" {
		queryText = "Search: (type to search...)"
	}
	queryBox := QueryBoxStyle.Render(queryText)
	b.WriteString(queryBox + "\n")
}

func (m *Model) renderSeparator(b *strings.Builder) {
	separator := SeparatorStyle.Render("────────────────────────────────────────")
	b.WriteString(separator + "\n")
}

func (m *Model) renderPathList(b *strings.Builder) {
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
}

func (m *Model) renderStatus(b *strings.Builder) {
	var status string
	if m.scanDone {
		status = StatusStyle.Render(fmt.Sprintf("--- Scan Completed (%d); Filtered (%d) ---", len(m.scannedPaths), len(m.filteredPaths)))
	} else {
		status = StatusStyle.Render(fmt.Sprintf("--- Scanning (%d); Filtered (%d) ---", len(m.scannedPaths), len(m.filteredPaths)))
	}

	b.WriteString(status)
}
