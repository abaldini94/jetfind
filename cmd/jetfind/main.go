package main

import (
	"fmt"
	"jetfind/internal/cli"
	"jetfind/internal/executor"
	"jetfind/internal/tui"
	"os"
)

func main() {
	config := cli.ParseArgs()
	
	selectedFile, err := tui.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
	
	exec := executor.New(config)
	if err := exec.Execute(selectedFile); err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		os.Exit(1)
	}
}
