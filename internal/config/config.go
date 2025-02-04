package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

const (
	KB = 1024
	MB = KB * KB
	GB = MB * MB
)

// TODO: make this more readable
func ParseSizeFilter(input string) (operator string, sizeInBytes int64, err error) {
	re := regexp.MustCompile(`(?i)^(<|>)?(\d+(?:\.\d+)?)(KB|MB|GB|B)?$`)
	matches := re.FindStringSubmatch(input)

	// matches for input of ">2KB" should be an array of [">2KB", ">", "2", "KB"]

	if len(matches) < 1 {
		return "", 0, fmt.Errorf("invalid size filter format")
	}

	operator = matches[1]

	if operator == "" {
		operator = ">" // default to greater than
		// TODO: implement greater/less than or equal to
	}

	value, err := strconv.ParseFloat(matches[2], 64) // parseFloat for fractional input e.g. ">2.5KB"
	if err != nil {
		return "", 0, fmt.Errorf("invalid size value")
	}

	unit := strings.ToUpper(matches[3])

	switch unit {
	case "KB":
		sizeInBytes = int64(value * KB)
	case "MB":
		sizeInBytes = int64(value * MB)
	case "GB":
		sizeInBytes = int64(value * GB)
	case "B":
		sizeInBytes = int64(value)
	default:
		return "", 0, fmt.Errorf("invalid size unit: %v", unit)
	}

	return operator, sizeInBytes, nil
}

type SizeFilter struct {
	IsFilter    bool
	SizeInBytes int64
	Operator    string
}

type Config struct {
	Count            int
	AllPackages      bool
	ShowHelp         bool
	ExplicitOnly     bool
	DependenciesOnly bool
	DateFilter       time.Time
	SizeFilter       SizeFilter
	SortBy           string
}

func ParseFlags(args []string) Config {
	var count int
	var allPackages bool
	var showHelp bool
	var explicitOnly bool
	var dependenciesOnly bool
	var dateFilter string
	var sizeFilter string
	var sortBy string

	// TODO: make this more readable
	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")
	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -n)")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")
	pflag.StringVar(&dateFilter, "date", "", "Filter packages installed on a specific date (YYYY-MM-DD)")
	pflag.StringVar(&sizeFilter, "size", "", "Filter pacakges by size (e.g. >20MB, <1GB)")
	pflag.StringVar(&sortBy, "sort", "date", "Sort by date/alphabetical/size")

	if err := pflag.CommandLine.Parse(args); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if allPackages {
		count = 0
	}

	var sizeFilterParsed SizeFilter

	if sizeFilter != "" {
		var err error
		sizeOperator, sizeInBytes, err := ParseSizeFilter(sizeFilter)
		if err != nil {
			log.Fatalf("Invalid size filter: %v\n", err)
		}

		sizeFilterParsed = SizeFilter{
			IsFilter:    true,
			SizeInBytes: sizeInBytes,
			Operator:    sizeOperator,
		}
	}

	var parsedDate time.Time

	if dateFilter != "" {
		var err error
		parsedDate, err = time.Parse("2006-01-02", dateFilter)
		if err != nil {
			log.Fatalf("Invalid date format: %v\n", err)
		}
	}

	return Config{
		Count:            count,
		AllPackages:      allPackages,
		ShowHelp:         showHelp,
		ExplicitOnly:     explicitOnly,
		DependenciesOnly: dependenciesOnly,
		DateFilter:       parsedDate,
		SizeFilter:       sizeFilterParsed,
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
      --sort <mode>       Sort by date (default) or "alphabetical" or "size:asc"/"size:desc"
      --size <filter>     Filter packages by size (e.g., >10MB, <1GB)
  -h, --help              Display this help message`)
}
