package scanengine

import (
	"math"
	"strings"
)

func createNgram(str string, ngramLen int) []string {
	if len(str) <= ngramLen {
		return []string{str}
	}
	ngrams := []string{}
	runes := []rune(str)
	for i := 0; i <= len(runes)-ngramLen; i++ {
		ngrams = append(ngrams, string(runes[i:i+ngramLen]))
	}
	return ngrams
}

func getOverlapCoefficient(str, pattern []string) float64 {
	smallerSetSize := math.Min(float64(len(str)), float64(len(pattern)))
	if smallerSetSize == 0 {
		return 0.0
	}

	set := make(map[string]struct{}, len(pattern))
	for _, v := range pattern {
		set[v] = struct{}{}
	}

	intrs := 0

	for _, ngram := range str {
		if _, ok := set[ngram]; ok {
			intrs++
		}
	}

	return float64(intrs) / smallerSetSize
}

func jaroWinkler(s1, s2 string) float64 {
	s1Len := len(s1)
	s2Len := len(s2)

	if s1Len == 0 || s2Len == 0 {
		return 0.0
	}

	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if s1 == s2 {
		return 1.0
	}

	matchWin := int(math.Floor(math.Max(float64(s1Len), float64(s2Len))/2)) - 1

	matchWin = max(0, matchWin)

	s1Matches := make([]bool, s1Len)
	s2Matches := make([]bool, s2Len)

	matches := 0

	for i := range s1Len {
		start := int(math.Max(0, float64(i-matchWin)))
		end := int(math.Min(float64(s2Len-1), float64(i+matchWin)))

		for j := start; j <= end; j++ {
			if !s2Matches[j] && s1[i] == s2[j] {
				s1Matches[i] = true
				s2Matches[j] = true
				matches++
				break
			}
		}
	}

	if matches == 0 {
		return 0.0
	}

	transpositions := 0
	k := 0
	for i := range s1Len {
		// consider only S1 char with matches
		if s1Matches[i] {
			for k < s2Len && !s2Matches[k] {
				//Advance to the next match in S2
				k++
			}
			if s1[i] != s2[k] {
				transpositions++
			}
			k++
		}
	}

	t := float64(transpositions / 2)
	m := float64(matches)

	sim := (m/float64(s1Len) + m/float64(s2Len) + (m-t)/m) / 3.0

	// Winkler bonus
	bonus := 0
	for i := 0; i < int(math.Min(4, math.Min(float64(s1Len), float64(s2Len)))); i++ {
		if s1[i] == s2[i] {
			bonus++
		} else {
			break
		}
	}

	return sim + (float64(bonus) * 0.1 * (1.0 - sim))

}
