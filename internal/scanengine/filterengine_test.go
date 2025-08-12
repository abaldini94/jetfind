package scanengine

import (
	"testing"
)

func TestBasicFilterEngine(t *testing.T) {
	testCases := []struct {
		name       string
		pathBuffer []ScanFilteredResult
		scanFilter ScanFilter
		expected   []ScanFilteredResult
	}{
		{
			name:       "Empty buffer",
			pathBuffer: []ScanFilteredResult{},
			scanFilter: ContainsFilter{Pattern: ""},
			expected:   make([]ScanFilteredResult, 0),
		},
		{
			name: "No matches",
			pathBuffer: []ScanFilteredResult{
				{Path: "root/parent/sub1/file1", Score: 1.0},
				{Path: "root/parent/sub2/file2", Score: 1.0},
			},
			scanFilter: ContainsFilter{Pattern: "nomatches"},
			expected:   make([]ScanFilteredResult, 0),
		},
		{
			name: "One match at file level",
			pathBuffer: []ScanFilteredResult{
				{Path: "root/parent/sub1/file1", Score: 1.0},
				{Path: "root/parent/sub2/file2", Score: 1.0},
			},
			scanFilter: ContainsFilter{Pattern: "file1"},
			expected:   []ScanFilteredResult{{Path: "root/parent/sub1/file1", Score: 1.0}},
		},
		{
			name: "One match at folder level",
			pathBuffer: []ScanFilteredResult{
				{Path: "root/parent/sub1/file1", Score: 1.0},
				{Path: "root/parent/sub2/file2", Score: 1.0},
			},
			scanFilter: ContainsFilter{Pattern: "sub2"},
			expected:   []ScanFilteredResult{{Path: "root/parent/sub2/file2", Score: 1.0}},
		},
		{
			name: "Full match",
			pathBuffer: []ScanFilteredResult{
				{Path: "root/parent/sub1/file1", Score: 1.0},
				{Path: "root/parent/sub2/file2", Score: 1.0},
			},
			scanFilter: ContainsFilter{Pattern: "root"},
			expected:   []ScanFilteredResult{{Path: "root/parent/sub1/file1", Score: 1.0}, {Path: "root/parent/sub2/file2", Score: 1.0}},
		},
	}

	for _, tc := range testCases {
		res := FilterEngine(tc.pathBuffer, tc.scanFilter)

		if len(res) != len(tc.expected) {
			t.Errorf("The length of the results is wrong. Expected %v; Obtained %v", len(tc.expected), len(res))
		}

		for i, p := range res {
			if p != tc.expected[i] {
				t.Errorf("The array does not match at sample %d. Expected %v; Obtained %v", i, tc.expected[i], p)
			}
		}
	}
}
