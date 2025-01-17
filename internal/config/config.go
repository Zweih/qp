package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type Config struct {
	Count            int
	AllPackages      bool
	ShowHelp         bool
	ExplicitOnly     bool
	DependenciesOnly bool
	DateFilter       time.Time
	SortBy           string
}

// reads cli arguments and populates a Config
func ParseFlags(args []string) Config {
	var count int
	var allPackages bool
	var showHelp bool
	var explicitOnly bool
	var dependenciesOnly bool
	var dateFilter string
	var sortBy string

	// flags
	// pflag.*VarP specifies a long name, a short name, and a default value
	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")
	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -n)")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")
	pflag.StringVar(&dateFilter, "date", "", "Filter packages installed on a specific date (YYYY-MM-DD)")
	pflag.StringVar(&sortBy, "sort", "date", "Sort by date/alphabetical/size")

	if allPackages {
		count = 0
	}

	var parsedDate time.Time

	if dateFilter != "" {
		var err error
		parsedDate, err = time.Parse("2006-01-02", dateFilter)
		if err != nil {
			log.Fatalf("Invalid date format: %v\n", err)
		}
	}

	// parse the flags in args
	if err := pflag.CommandLine.Parse(args); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	return Config{
		Count:            count,
		AllPackages:      allPackages,
		ShowHelp:         showHelp,
		ExplicitOnly:     explicitOnly,
		DependenciesOnly: dependenciesOnly,
		DateFilter:       parsedDate,
		SortBy:           sortBy,
	}
}

func PrintHelp() {
	fmt.Println(`Usage: yaylog [options]

Options:
  -n, --number <number>   Display the specified number of recent packages (default: 20)
  -a, --all               Show all installed packages (ignores -n)
  -e, --explicit          Show only explicitly installed packages
  -d, --dependencies      Show only packages installed as dependencies
      --date <YYYY-MM-DD> Filter packages installed on a specific date
      --sort <mode>       Sort by date (default) or "alphabetical" or "size"
  -h, --help              Display this help message`)
}
