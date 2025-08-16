package scanengine

import (
	"path/filepath"
	"strings"
)

const (
	AlgoJaroWinkler = "jarowinkler"
	AlgoNGram       = "ngram"
	AlgoLevenshtein = "levenshtein"
)

type ScanFilteredResult struct {
	Path  string
	Score float64
}

type ScanFilter interface {
	Apply(path string) (ScanFilteredResult, bool)
}

type NoFilter struct {
	Pattern string
}

type ExactFilter struct {
	Pattern string
}

type ContainsFilter struct {
	Pattern string
}

type FuzzyFilter struct {
	Pattern    string
	Threashold float64
	Algo       string
}

func (nf NoFilter) Apply(path string) (ScanFilteredResult, bool) {
	return ScanFilteredResult{Path: path, Score: 1.0}, true
}

func (ef ExactFilter) Apply(path string) (ScanFilteredResult, bool) {
	if path == ef.Pattern {
		return ScanFilteredResult{Path: path, Score: 1.0}, true
	}
	return ScanFilteredResult{}, false
}

func (cf ContainsFilter) Apply(path string) (ScanFilteredResult, bool) {
	if strings.Contains(strings.ToLower(path), strings.ToLower(cf.Pattern)) {
		return ScanFilteredResult{Path: path, Score: 1.0}, true
	}
	return ScanFilteredResult{}, false
}

func (ff FuzzyFilter) Apply(path string) (ScanFilteredResult, bool) {
	if ff.Algo == AlgoNGram {
		pathNgram := createNgram(path, 2)
		patternNgram := createNgram(ff.Pattern, 2)
		overlap := getOverlapCoefficient(pathNgram, patternNgram)
		if overlap > ff.Threashold {
			return ScanFilteredResult{Path: path, Score: overlap}, true
		}
		return ScanFilteredResult{}, false
	} else if ff.Algo == AlgoJaroWinkler {
		slashedPath := filepath.ToSlash(path)
		parts := strings.Split(slashedPath, "/")
		jwSim := jaroWinkler(parts[len(parts)-1], ff.Pattern)
		if jwSim > ff.Threashold {
			return ScanFilteredResult{Path: path, Score: jwSim}, true
		}
		return ScanFilteredResult{}, false
	} else if ff.Algo == AlgoLevenshtein {
		slashedPath := filepath.ToSlash(path)
		parts := strings.Split(slashedPath, "/")
		levSim := levenshteinSimilarity(parts[len(parts)-1], ff.Pattern)
		if levSim > ff.Threashold {
			return ScanFilteredResult{Path: path, Score: levSim}, true
		}
		return ScanFilteredResult{}, false
	} else {
		panic("Unknown Fuzzy matching algorithm")
	}
}
