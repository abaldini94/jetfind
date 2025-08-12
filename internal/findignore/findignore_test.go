package findingnore

import (
	"os"
	"testing"
)

func createTempFindIgnoreFile(t *testing.T, ignoreContent string) (string, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", ".findignore")
	if err != nil {
		t.Fatalf("Impossible to create temporary findignore file: %v", err)
	}

	if _, err := tmpFile.WriteString(ignoreContent); err != nil {
		t.Fatalf("Impossible to write in the temporary findignore file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Impossible to close the temporary findignore file: %v", err)
	}

	cleanup := func() {
		os.Remove(tmpFile.Name())
	}
	return tmpFile.Name(), cleanup
}

func TestParseGlob(t *testing.T) {
	testCases := []struct {
		glob     string
		path     string
		expected bool
	}{
		{"*.log", "file.log", true},
		{"*.log", "file.log.txt", false},
		{"src/*.js", "src/app.js", true},
		{"src/*.js", "src/lib/app.js", false},
		{"data", "data", true},
		{"data/", "data/", true},
	}

	for _, tc := range testCases {
		re, err := parseGlob(tc.glob)
		if err != nil {
			t.Fatalf("parseGlob failed for '%s': %v", tc.glob, err)
		}

		if re.MatchString(tc.path) != tc.expected {
			t.Errorf("Glob '%s' for path '%s': expected %v, obtained %v", tc.glob, tc.path, tc.expected, !tc.expected)
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	ignoreContent := `
# comments
*.log
build/
!important.log
node_modules/
`
	ignoreFile, cleanup := createTempFindIgnoreFile(t, ignoreContent)
	defer cleanup()

	fi, err := New(ignoreFile, true)
	if err != nil {
		t.Fatalf("New() raised an unexpected error: %v", err)
	}

	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "simple log file", path: "app.log", expected: true},
		{name: "log file in subdirectory", path: "src/app.log", expected: true},
		{name: "Not ignored file", path: "main.go", expected: false},
		{name: "Directory build", path: "build/", expected: true},
		{name: "File within ignored directory", path: "build/output.bin", expected: true},
		{name: "File named as an ignored directory", path: "mybuild", expected: false},
		{name: "Negated log file", path: "important.log", expected: false},
		{name: "Directory node_modules", path: "node_modules/react", expected: true},
		{name: "Hidden file", path: ".env", expected: true},
		{name: "File within hidden directory", path: "src/.cache/data", expected: true},
	}

	for _, tc := range testCases {
		result := fi.ShouldIgnore(tc.path)
		if result != tc.expected {
			t.Errorf("Path: %s - Expected %v Obtained %v", tc.path, tc.expected, result)
		}
	}
}
