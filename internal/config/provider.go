package config

import (
	"fmt"
	"os"
	"qp/internal/about"
	"qp/internal/quipple/completion"
)

func (c *CliConfigProvider) GetConfig() (*Config, error) {
	cfg, err := ParseFlags(os.Args[1:])
	if err != nil {
		return &Config{}, err
	}

	if cfg.ShowHelp {
		PrintHelp()
		os.Exit(0)
	}

	if cfg.ShowVersion {
		about.PrintVersionInfo()
		os.Exit(0)
	}

	if cfg.ShowCompletion {
		fmt.Print(completion.GetCompletions())
		os.Exit(0)
	}

	return &cfg, nil
}
