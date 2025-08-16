package cli

import (
	"testing"
)

func TestNewExecutor(t *testing.T) {
	flags := &CliFlags{
		PostCmd: "vim",
		Help:    false,
		Version: false,
	}

	executor := NewExecutor(flags)

	if executor == nil {
		t.Error("NewExecutor() returned nil")
	}

	if executor.cliFlags.PostCmd != "vim" {
		t.Errorf("Expected PostCmd 'vim', got '%s'", executor.cliFlags.PostCmd)
	}
}

func TestExecutorExecuteWithoutPostCommand(t *testing.T) {
	flags := &CliFlags{
		PostCmd: "",
	}

	executor := NewExecutor(flags)

	err := executor.Execute("test.txt")
	if err != nil {
		t.Errorf("Execute() without post command should not return error, got: %v", err)
	}
}

