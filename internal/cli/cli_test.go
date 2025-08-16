package cli

import (
	"testing"
)

func TestCliFlagsHasPostCommand(t *testing.T) {
	tests := []struct {
		name     string
		postCmd  string
		expected bool
	}{
		{
			name:     "empty post command",
			postCmd:  "",
			expected: false,
		},
		{
			name:     "has post command",
			postCmd:  "vim",
			expected: true,
		},
		{
			name:     "has post command with args",
			postCmd:  "code --wait",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CliFlags{
				PostCmd: tt.postCmd,
			}

			result := c.HasPostCommand()
			if result != tt.expected {
				t.Errorf("HasPostCommand() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCliFlagsDefaultValues(t *testing.T) {
	c := &CliFlags{}

	if c.PostCmd != "" {
		t.Errorf("Default PostCmd should be empty, got %s", c.PostCmd)
	}

	if c.Help != false {
		t.Errorf("Default Help should be false, got %v", c.Help)
	}

	if c.Version != false {
		t.Errorf("Default Version should be false, got %v", c.Version)
	}
}

