package cli

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	PostCmd string
	Help    bool
	Version bool
}

func ParseArgs() *Config {
	config := &Config{}

	flag.StringVar(&config.PostCmd, "post-cmd", "", "Command to execute after a file has been selected")
	flag.BoolVar(&config.Help, "help", false, "Show help message")
	flag.BoolVar(&config.Help, "h", false, "Show help message")
	flag.BoolVar(&config.Version, "version", false, "Show version information")
	flag.BoolVar(&config.Version, "v", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A configurable file finder with interactive selection.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s		      Select and print file path\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --post-cmd vim    Open selected file with vim\n", os.Args[0])
	}

	flag.Parse()

	if config.Help {
		flag.Usage()
		os.Exit(0)
	}

	if config.Version {
		fmt.Println("jetfind version 0.1.0")
		os.Exit(0)
	}

	return config
}

func (c *Config) HasPostCommand() bool {
	return c.PostCmd != ""
}

