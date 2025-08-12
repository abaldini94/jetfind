package scanengine

import "testing"

func TestExactFilter(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  string
		path     string
		expected bool
	}{
		{
			name:     "Exact Match",
			pattern:  "README.md",
			path:     "README.md",
			expected: true,
		},
		{
			name:     "No Match",
			pattern:  "main.go",
			path:     "internal/main.go",
			expected: false,
		},
		{
			name:     "Non Valid Partial Match",
			pattern:  "README",
			path:     "README.md",
			expected: false,
		},
		{
			name:     "Empty String",
			pattern:  "",
			path:     "",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := ExactFilter{Pattern: tc.pattern}
			_, result := filter.Apply(tc.path)
			if result != tc.expected {
				t.Errorf("Expected %v, but obtained %v for path '%s' with pattern '%s'", tc.expected, result, tc.path, tc.pattern)
			}
		})
	}
}

func TestContainsFilter(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  string
		path     string
		expected bool
	}{
		{
			name:     "Contains Substring",
			pattern:  "test",
			path:     "/home/user/project/main_test.go",
			expected: true,
		},
		{
			name:     "Does Not Contain The Substring",
			pattern:  "docs",
			path:     "/home/user/project/main.go",
			expected: false,
		},
		{
			name:     "Match Case-Insensitive (lower)",
			pattern:  "readme",
			path:     "/home/user/project/README.md",
			expected: true,
		},
		{
			name:     "Match Case-Insensitive (upper)",
			pattern:  "README",
			path:     "/home/user/project/readme.md",
			expected: true,
		},
		{
			name:     "Empty Pattern",
			pattern:  "",
			path:     "whatever/path",
			expected: true,
		},
		{
			name:     "Empty Path",
			pattern:  "test",
			path:     "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := ContainsFilter{Pattern: tc.pattern}
			_, result := filter.Apply(tc.path)
			if result != tc.expected {
				t.Errorf("Expected %v, but obtained %v for path '%s' with pattern '%s'", tc.expected, result, tc.path, tc.pattern)
			}
		})
	}
}

func TestFuzzyFilter(t *testing.T) {
	testCases := []struct {
		name      string
		pattern   string
		path      string
		threshold float64
		expected  bool
	}{
		{
			name:      "Fuzzy Substring",
			pattern:   "jtfnd",
			threshold: 0.8,
			path:      "/home/user/jetfind",
			expected:  true,
		},
		{
			name:      "No Substring",
			pattern:   "kyusydus",
			threshold: 0.8,
			path:      "/home/user/jetfind/main_test.go",
			expected:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := FuzzyFilter{Pattern: tc.pattern, Threashold: tc.threshold, Algo: AlgoJaroWinkler}
			_, result := filter.Apply(tc.path)
			if result != tc.expected {
				t.Errorf("Expected %v, but obtained %v for path '%s' with pattern '%s'", tc.expected, result, tc.path, tc.pattern)
			}
		})
	}
}
