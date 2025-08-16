package scanengine

import (
	"math"
	"reflect"
	"testing"
)

func TestCreateNgrams(t *testing.T) {
	testCases := []struct {
		name     string
		str      string
		ngramLen int
		expected []string
	}{
		{
			name:     "Ngram bigger than string lenght",
			str:      "a",
			ngramLen: 2,
			expected: []string{"a"},
		},
		{
			name:     "Two characters",
			str:      "ab",
			ngramLen: 2,
			expected: []string{"ab"},
		},
		{
			name:     "Three characters",
			str:      "abc",
			ngramLen: 2,
			expected: []string{"ab", "bc"},
		},
		{
			name:     "String with spaces",
			str:      "a b",
			ngramLen: 2,
			expected: []string{"a ", " b"},
		},
		{
			name:     "String with special characters",
			str:      "a-b",
			ngramLen: 2,
			expected: []string{"a-", "-b"},
		},
		{
			name:     "Real case string with 2gram",
			str:      "project",
			ngramLen: 2,
			expected: []string{"pr", "ro", "oj", "je", "ec", "ct"},
		},
		{
			name:     "Real case string with 3gram",
			str:      "project",
			ngramLen: 3,
			expected: []string{"pro", "roj", "oje", "jec", "ect"},
		},
	}
	for _, tc := range testCases {
		res := createNgram(tc.str, tc.ngramLen)
		if !reflect.DeepEqual(res, tc.expected) {
			t.Errorf("Test '%s' failed: got %v, expected %v", tc.name, res, tc.expected)
		}
	}
}

func TestGetOverlapCoefficient(t *testing.T) {
	testCases := []struct {
		name     string
		str      []string
		pattern  []string
		expected float64
	}{
		{
			name:     "Empty pattern",
			str:      []string{"ab", "bc", "cd"},
			pattern:  []string{},
			expected: 0.0,
		},
		{
			name:     "Empty string",
			str:      []string{},
			pattern:  []string{"ab", "bc", "cd"},
			expected: 0.0,
		},
		{
			name:     "No overlap",
			str:      []string{"ab", "bc", "cd"},
			pattern:  []string{"ef", "fg", "hi"},
			expected: 0.0,
		},
		{
			name:     "Single overlap equal len strings",
			str:      []string{"ab", "bc", "cd", "zk"},
			pattern:  []string{"ab", "ef", "hi", "ob"},
			expected: 0.25,
		},
		{
			name:     "Two overlap equal len strings",
			str:      []string{"ab", "bc", "cd", "zk"},
			pattern:  []string{"ab", "bc", "hi", "ob"},
			expected: 0.50,
		},
		{
			name:     "Full overlap",
			str:      []string{"ab", "bc", "cd", "zk"},
			pattern:  []string{"ab", "bc", "cd", "zk"},
			expected: 1.0,
		},
		{
			name:     "Single overlap different len strings",
			str:      []string{"ab", "bc", "cd", "zk"},
			pattern:  []string{"ab", "ef", "hi", "ob", "gf", "ju", "lo", "pq"},
			expected: 0.25,
		},
	}

	for _, tc := range testCases {
		res := getOverlapCoefficient(tc.str, tc.pattern)
		if res != tc.expected {
			t.Errorf("Test %s failed: got %v, expected %v", tc.name, res, tc.expected)
		}
	}

}

func TestJaroWinkler(t *testing.T) {
	const epsilon = 0.001

	testCases := []struct {
		name     string
		str1     string
		str2     string
		expected float64
	}{
		{
			name:     "Perfect match",
			str1:     "apple",
			str2:     "apple",
			expected: 1.0,
		},
		{
			name:     "No match",
			str1:     "zpple",
			str2:     "banana",
			expected: 0.0,
		},
		{
			name:     "One is empty",
			str1:     "",
			str2:     "apple",
			expected: 0.0,
		},
		{
			name:     "transposition",
			str1:     "martha",
			str2:     "marhta",
			expected: 0.961,
		},
		{
			name:     "typo missing char",
			str1:     "project",
			str2:     "projct",
			expected: 0.9714,
		},
		{
			name:     "typo added char",
			str1:     "project",
			str2:     "projectz",
			expected: 0.975,
		},
		{
			name:     "transposed filename",
			str1:     "main.go",
			str2:     "mian.go",
			expected: 0.957,
		},
		{
			name:     "case insensitive",
			str1:     "Test",
			str2:     "test",
			expected: 1.0,
		},
		{
			name:     "spaces",
			str1:     "a b c",
			str2:     "a c b",
			expected: 0.786,
		},
	}

	for _, tc := range testCases {
		res := jaroWinkler(tc.str1, tc.str2)
		if math.Abs(res-tc.expected) > epsilon {
			t.Errorf("Test %s failed: got %v, expected %v", tc.name, res, tc.expected)
		}
	}

}

func TestLevenshtainSimilarity(t *testing.T) {
	const epsilon = 0.001

	testCases := []struct {
		name     string
		str1     string
		str2     string
		expected float64
	}{
		{
			name:     "Both empty; returns 1",
			str1:     "",
			str2:     "",
			expected: 1.0,
		},
		{
			name:     "S1 empty; returns 0",
			str1:     "",
			str2:     "hello",
			expected: 0.0,
		},
		{
			name:     "S2 empty; returns 0",
			str1:     "hello",
			str2:     "",
			expected: 0.0,
		},
		{
			name:     "Yellow vs. Hello",
			str1:     "yellow",
			str2:     "hello",
			expected: 0.667,
		},
		{
			name:     "Table vs. Table",
			str1:     "table",
			str2:     "table",
			expected: 1.0,
		},
		{
			name:     "Bill vs. Gills",
			str1:     "bill",
			str2:     "gills",
			expected: 0.6,
		},
		{
			name:     "Cat vs. Dog",
			str1:     "cat",
			str2:     "dog",
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		res := levenshteinSimilarity(tc.str1, tc.str2)
		if math.Abs(res-tc.expected) > epsilon {
			t.Errorf("Test %s failed: got %v, expected %v", tc.name, res, tc.expected)
		}
	}

}
