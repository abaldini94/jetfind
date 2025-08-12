package executor

import (
	"fmt"
	"jetfind/internal/cli"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	config *cli.Config
}

func New(config *cli.Config) *Executor {
	return &Executor{
		config: config,
	}
}

func (e *Executor) Execute(selectedFile string) error {
	if !e.config.HasPostCommand() {
		fmt.Println(selectedFile)
		return nil
	}
	
	return e.executeCommand(selectedFile)
}

func (e *Executor) executeCommand(selectedFile string) error {
	cmdParts := strings.Fields(e.config.PostCmd)
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
		return fmt.Errorf("failed to execute command '%s': %w", e.config.PostCmd, err)
	}
	
	return nil
}