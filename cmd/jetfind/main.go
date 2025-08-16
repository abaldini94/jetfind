package main

import (
	"fmt"
	"jetfind/internal/cli"
	"jetfind/internal/tui"
	"os"
)

func main() {
	cliFalgs := cli.ParseArgs()

	selectedFile, err := tui.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}

	exec := cli.NewExecutor(cliFalgs)
	if err := exec.Execute(selectedFile); err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		os.Exit(1)
	}
}
