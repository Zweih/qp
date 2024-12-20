package config

import (
	"fmt"
	//  "log"
	//  "time"
	//  "os"

	"github.com/spf13/pflag"
)

type Config struct {
	Count            int
	ShowHelp         bool
	ExplicitOnly     bool
	DependenciesOnly bool
}

// reads cli arguments and populates a Config
func ParseFlags(args []string) Config {
	var count int
	var showHelp bool
	var explicitOnly bool
	var dependenciesOnly bool

	// flags
	// pflag.*VarP specifies a long name, a short name, and a default value
	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")

	// parse the flags in args
	pflag.CommandLine.Parse(args)

	return Config{
		Count:            count,
		ShowHelp:         showHelp,
		ExplicitOnly:     explicitOnly,
		DependenciesOnly: dependenciesOnly,
	}
}

// help message
func PrintHelp() {
	fmt.Println(`Usage: yaylog [options]

Options:
  -n, --number <number>   Display the specified number of recent packages (default: 20)
  -h, --help              Display this help message
`)
}
