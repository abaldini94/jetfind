package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isValidHexColor(color string) bool {
	if color == "" {
		return true
	}
	validLengths := []int{4, 7, 9}
	validLength := false
	for _, length := range validLengths {
		if len(color) == length {
			validLength = true
			break
		}
	}
	if !validLength || color[0] != '#' {
		return false
	}
	for i := 1; i < len(color); i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func (c *Config) Validate() error {

	validFilterTypes := []string{"fuzzy", "contains"}
	if !contains(validFilterTypes, c.Filter.Type) {
		return fmt.Errorf("invalid filter type: %s. Must be one of: %v", c.Filter.Type, validFilterTypes)
	}

	if c.Filter.Type == "fuzzy" {
		validAlgos := []string{"jarowinkler", "ngram", "levenshtein"}
		if !contains(validAlgos, c.Filter.Algo) {
			return fmt.Errorf("invalid filter algorithm: %s. Must be one of: %v", c.Filter.Algo, validAlgos)
		}

		if c.Filter.Threashold < 0 || c.Filter.Threashold > 1 {
			return fmt.Errorf("invalid filter threashold: %.2f. Must be in the [0, 1] interval", c.Filter.Threashold)
		}
	}

	if c.Findignore.Enable {
		findignorePath := filepath.Join(GetConfigDir(), ".findignore")
		if _, err := os.Stat(findignorePath); os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist: %s", findignorePath)
		}
	}

	colorFields := []string{
		c.Tui.HighlightedFile.Foreground,
		c.Tui.HighlightedFile.Background,
		c.Tui.QueryBox.TextForeground,
		c.Tui.QueryBox.TextBackground,
		c.Tui.QueryBox.BorderForeground,
	}

	for _, color := range colorFields {
		if color != "" && !isValidHexColor(color) {
			return fmt.Errorf("invalid hex color format: %s", color)
		}
	}

	return nil
}
