package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	cliFlags CliFlags
}

func NewExecutor(cliFlags *CliFlags) *Executor {
	return &Executor{
		cliFlags: *cliFlags,
	}
}

func (e *Executor) Execute(selectedFile string) error {
	if !e.cliFlags.HasPostCommand() {
		fmt.Println(selectedFile)
		return nil
	}

	return e.executeCommand(selectedFile)
}

func (e *Executor) executeCommand(selectedFile string) error {
	cmdParts := strings.Fields(e.cliFlags.PostCmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmdName := cmdParts[0]
	cmdArgs := append(cmdParts[1:], selectedFile)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command '%s': %w", e.cliFlags.PostCmd, err)
	}

	return nil
}

