package scanengine

import (
	findignore "jetfind/internal/findignore"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// Structure:
// root/
// ├── file1.txt
// ├── sub/
// │   ├── file2.go
// │   └── nested/
// │       └── file3.md
// ├── ignored_dir/
// │   └── ignored_file.txt
// └── .git/
//	    └── config

func createTestDir(t *testing.T) (string, func()) {
	t.Helper()
	root, err := os.MkdirTemp("", "scanner_test_*")
	if err != nil {
		t.Fatalf("Impossible to load directory: %v", err)
	}

	mustMkdir(t, filepath.Join(root, "sub", "nested"))
	mustMkdir(t, filepath.Join(root, "ignored_dir"))
	mustMkdir(t, filepath.Join(root, ".git"))

	mustWriteFile(t, filepath.Join(root, "file1.txt"), "test")
	mustWriteFile(t, filepath.Join(root, "sub", "file2.go"), "package main")
	mustWriteFile(t, filepath.Join(root, "sub", "nested", "file3.md"), "# Title")
	mustWriteFile(t, filepath.Join(root, "ignored_dir", "ignored_file.txt"), "ignore")
	mustWriteFile(t, filepath.Join(root, ".git", "config"), "config")

	cleanup := func() {
		os.RemoveAll(root)
	}

	return root, cleanup
}

func mustMkdir(t *testing.T, path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Impossible to create directory %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Impossible to write the file %s: %v", path, err)
	}
}

func collectResults(resultsChan <-chan string) []string {
	var results []string

	for res := range resultsChan {
		results = append(results, res)
	}
	sort.Strings(results)
	return results
}

func TestScanEmptyDirectory(t *testing.T) {
	root, err := os.MkdirTemp("", "empty_test")
	if err != nil {
		t.Fatalf("Impossibile creare directory temporanea: %v", err)
	}
	defer os.RemoveAll(root)

	config := Config{
		Root:       root,
		NumWorkers: 2,
		FindIgnore: nil,
	}
	scanner := New(config)
	resultsChan := scanner.Run()
	results := collectResults(resultsChan)

	if len(results) != 0 {
		t.Errorf("Empty directory scan must return 0 results, but %d paths were obtained", len(results))
	}
}

func TestBasicScan(t *testing.T) {
	root, cleanup := createTestDir(t)
	defer cleanup()

	config := Config{
		Root:       root,
		NumWorkers: 2,
		FindIgnore: nil,
	}
	scanner := New(config)
	resultsChan := scanner.Run()
	results := collectResults(resultsChan)

	expected := []string{
		filepath.Join(root, ".git", "config"),
		filepath.Join(root, "file1.txt"),
		filepath.Join(root, "ignored_dir", "ignored_file.txt"),
		filepath.Join(root, "sub", "file2.go"),
		filepath.Join(root, "sub", "nested", "file3.md"),
	}
	sort.Strings(expected)

	if len(results) != len(expected) {
		t.Fatalf("Wrong number of results. Expected: %d, Got: %d\nExpected: %v\nGot: %v", len(expected), len(results), expected, results)
	}

	for i := range results {
		if results[i] != expected[i] {
			t.Errorf("Unexpected results at %d. Expected: %s, Got: %s", i, expected[i], results[i])
		}
	}
}

func TestScanWithFindIgnore(t *testing.T) {
	root, cleanup := createTestDir(t)
	defer cleanup()

	ignoreContent := "sub/\n"
	ignoreFilePath := filepath.Join(root, ".findignore")
	mustWriteFile(t, ignoreFilePath, ignoreContent)

	fi, err := findignore.New(ignoreFilePath, true)
	if err != nil {
		t.Fatalf("Impossible to create temporary findignore file: %v", err)
	}
	config := Config{
		Root:       root,
		NumWorkers: 2,
		FindIgnore: fi,
	}
	scanner := New(config)
	resultsChan := scanner.Run()
	results := collectResults(resultsChan)

	expected := []string{
		filepath.Join(root, "file1.txt"),
		filepath.Join(root, "ignored_dir", "ignored_file.txt"),
	}
	sort.Strings(expected)

	if len(results) != len(expected) {
		t.Fatalf("Number of results is wrong. Expected: %d, Obtained: %d\nExpected: %v\nObtained: %v", len(expected), len(results), expected, results)
	}

	for i := range results {
		if results[i] != expected[i] {
			t.Errorf("Unexpected result at position %d. Expected: %s, Obtained: %s", i, expected[i], results[i])
		}
	}
}
