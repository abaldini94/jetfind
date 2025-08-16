package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	validYAML := `filter:
  type: "fuzzy"
  algorithm: "jarowinkler"
  threashold: 0.8
findignore:
  enable: false
  hidden_ignore: true
tui:
  highlighted_file:
    foreground: "#FFFFFF"
    background: "#4B5563"
  query_box:
    text_foreground: "#F9FAFB"
    text_background: "#374151"
    border_foreground: "#6B7280"`

	tmpFile := createTempConfig(t, validYAML)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.Filter.Type != "fuzzy" {
		t.Errorf("Expected filter type 'fuzzy', got '%s'", cfg.Filter.Type)
	}
	if cfg.Filter.Threashold != 0.8 {
		t.Errorf("Expected threshold 0.8, got %f", cfg.Filter.Threashold)
	}
	if cfg.Tui.HighlightedFile.Foreground != "#FFFFFF" {
		t.Errorf("Expected foreground '#FFFFFF', got '%s'", cfg.Tui.HighlightedFile.Foreground)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{
			name: "invalid filter type",
			yaml: `filter:
  type: "invalid"
  algorithm: "jarowinkler"
  threashold: 0.8`,
		},
		{
			name: "invalid threshold",
			yaml: `filter:
  type: "fuzzy"
  algorithm: "jarowinkler"
  threashold: 1.5`,
		},
		{
			name: "invalid color",
			yaml: `tui:
  highlighted_file:
    foreground: "invalid_color"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := createTempConfig(t, tt.yaml)
			defer os.Remove(tmpFile)

			_, err := Load(tmpFile)
			if err == nil {
				t.Error("Expected validation error, got none")
			}
		})
	}
}

func TestLoadNonexistentConfig(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got none")
	}
}

func createTempConfig(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	return tmpFile
}
