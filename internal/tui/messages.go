package tui

import (
	"jetfind/internal/scanengine"
)

type errMsg error
type newPathMsg scanengine.ScanFilteredResult
type scanDoneMsg struct{}
