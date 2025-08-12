package tui

import (
	"jetfind/internal/config"
	"jetfind/internal/findignore"
	"jetfind/internal/scanengine"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cfg             *config.Config
	userQuery       string
	scanChan        <-chan string
	scannedPaths    []scanengine.ScanFilteredResult
	filteredPaths   []scanengine.ScanFilteredResult
	filterRequested bool
	scanDone        bool
	scanErr         error
	cursor          int
	offset          int
	height          int
	SelectedFile    string
}

func NewModel(cfg *config.Config) *Model {
	return &Model{
		cfg:          cfg,
		scannedPaths: []scanengine.ScanFilteredResult{},
		cursor:       0,
		offset:       0,
	}
}

func (m *Model) Init() tea.Cmd {
	var fi *findingnore.FindIgnore
	var err error
	cfgDir := config.GetConfigDir()
	if m.cfg.Findignore.Enable {
		fi, err = findingnore.New(filepath.Join(cfgDir, ".findignore"), true)
		if err != nil {
			return func() tea.Msg {
				return errMsg(err)
			}
		}
	} else {
		fi = nil
	}

	scanCfg := scanengine.Config{
		Root:       "./",
		FindIgnore: fi,
	}

	scanner := scanengine.New(scanCfg)
	m.scanChan = scanner.Run()
	return popFromScanChanCmd(m.scanChan)
}
