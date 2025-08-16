package config

import (
	"testing"
)

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  *Default,
			wantErr: false,
		},
		{
			name: "invalid filter type",
			config: Config{
				Filter: FilterConfig{Type: "invalid"},
			},
			wantErr: true,
		},
		{
			name: "invalid algorithm for fuzzy",
			config: Config{
				Filter: FilterConfig{
					Type: "fuzzy",
					Algo: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "threshold out of range",
			config: Config{
				Filter: FilterConfig{
					Type:       "fuzzy",
					Algo:       "jarowinkler",
					Threashold: 1.5,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid hex color",
			config: Config{
				Filter: Default.Filter,
				Tui: TuiConfig{
					HighlightedFile: HighlightedFileConfig{
						Foreground: "not_a_hex_color",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHexColorValidation(t *testing.T) {
	tests := []struct {
		color string
		valid bool
	}{
		{"#FFFFFF", true},
		{"#000", true},      // Short format
		{"#00000000", true}, // With alpha
		{"", true},          // Empty is OK
		{"#ZZZZZZ", false},  // Invalid chars
		{"FFFFFF", false},   // Missing #
		{"#FF", false},      // Too short
		{"#FFFFFFF", false}, // Wrong length
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			if got := isValidHexColor(tt.color); got != tt.valid {
				t.Errorf("isValidHexColor(%q) = %v, want %v", tt.color, got, tt.valid)
			}
		})
	}
}

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	if !contains(slice, "banana") {
		t.Error("Expected contains to return true for 'banana'")
	}

	if contains(slice, "grape") {
		t.Error("Expected contains to return false for 'grape'")
	}
}
