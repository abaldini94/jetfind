package findingnore

import (
	"os"
	"regexp"
	"strings"
)

func parseGlob(glob string) (*regexp.Regexp, error) {
	var re strings.Builder

	for _, c := range glob {
		switch c {
		case '*':
			re.WriteString("[^/]*")

		case '.', '\\', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|':
			re.WriteRune('\\')
			re.WriteRune(c)

		default:
			re.WriteRune(c)
		}
	}

	re.WriteString("$")
	return regexp.Compile(re.String())
}

type IgnorePattern struct {
	Pattern   *regexp.Regexp
	IsNegated bool
	IsDir     bool
}

type FindIgnore struct {
	Ignore       []IgnorePattern
	IgnoreHidden bool
}

func New(filename string, ignoreHidden bool) (*FindIgnore, error) {
	lines, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(lines), "\n")
	patterns := make([]IgnorePattern, 0)
	for _, line := range s {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		isNegated := strings.HasPrefix(line, "!")
		if isNegated {
			line = line[1:]
		}

		isDir := strings.HasSuffix(line, "/")
		patternString := line
		if isDir {
			patternString = strings.TrimSuffix(patternString, "/")
		}

		re, err := parseGlob(patternString)
		if err != nil {
			return nil, err
		}

		if isDir {
			originalRe := re.String()
			newReStr := strings.TrimSuffix(originalRe, "$") + `/.*?$`
			re, err = regexp.Compile(newReStr)
			if err != nil {
				return nil, err
			}
		}

		p := IgnorePattern{
			Pattern:   re,
			IsNegated: isNegated,
			IsDir:     isDir,
		}
		patterns = append(patterns, p)
	}

	return &FindIgnore{Ignore: patterns, IgnoreHidden: ignoreHidden}, nil
}

func (f *FindIgnore) ShouldIgnore(path string) bool {
	lastMatchIndex := -1

	if f.IgnoreHidden {
		parts := strings.SplitSeq(path, string(os.PathSeparator))
		for part := range parts {
			if len(part) > 1 && strings.HasPrefix(part, ".") {
				return true
			}
		}
	}

	for i, p := range f.Ignore {
		if p.Pattern.MatchString(path) {
			lastMatchIndex = i
		}
	}

	if lastMatchIndex == -1 {
		return false
	}

	return !f.Ignore[lastMatchIndex].IsNegated

}
