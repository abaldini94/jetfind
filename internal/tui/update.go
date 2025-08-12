package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"jetfind/internal/scanengine"
)

func popFromScanChanCmd(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		if p, ok := <-ch; ok {
			return newPathMsg(scanengine.ScanFilteredResult{Path: p, Score: 1.0})
		}
		return scanDoneMsg{}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case newPathMsg:
		m.scannedPaths = append(m.scannedPaths, scanengine.ScanFilteredResult(msg))
		return m, popFromScanChanCmd(m.scanChan)
	case errMsg:
		m.scanErr = msg
		return m, nil
	case scanDoneMsg:
		m.scanDone = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.filteredPaths) > 0 && m.cursor < len(m.filteredPaths) {
				m.SelectedFile = m.filteredPaths[m.cursor].Path
				return m, tea.Quit
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.filteredPaths)-1 {
				m.cursor++
			}
		case "backspace":
			if m.userQuery != "" {
				m.userQuery = m.userQuery[:len(m.userQuery)-1]
			}
			if m.userQuery != "" {
				m.filterRequested = true
				m.cursor = 0
			}
		}

		if msg.Type == tea.KeyRunes {
			m.userQuery += string(msg.Runes)
			m.filterRequested = true
			m.cursor = 0
		}

		visibleLines := m.height - 4

		if m.cursor < m.offset {
			m.offset = m.cursor
		}
		if m.cursor >= m.offset+visibleLines {
			m.offset = m.cursor - visibleLines + 1
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.height = msg.Height - 1
		return m, nil
	}
	return m, nil
}
