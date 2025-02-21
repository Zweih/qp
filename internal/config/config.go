package config

import (
	"fmt"
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

type SizeFilter struct {
	IsFilter    bool
	SizeInBytes int64
	Operator    string
}

type DateFilter struct {
	StartDate    time.Time
	EndDate      time.Time
	IsExactMatch bool
}

type Config struct {
	Count             int
	AllPackages       bool
	ShowHelp          bool
	ShowFullTimestamp bool
	DisableProgress   bool
	ExplicitOnly      bool
	DependenciesOnly  bool
	DateFilter        DateFilter
	SizeFilter        SizeFilter
	NameFilter        string
	SortBy            string
	OptionalColumns   []string
}

func ParseFlags(args []string) (Config, error) {
	var count int
	var allPackages bool
	var showHelp bool
	var showFullTimestamp bool
	var showVersion bool
	var disableProgress bool
	var explicitOnly bool
	var dependenciesOnly bool
	var dateFilter string
	var sizeFilter string
	var nameFilter string
	var sortBy string

	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")

	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -n)")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")
	pflag.BoolVarP(&showFullTimestamp, "full-timestamp", "", false, "Show full timestamp instead of just the date")
	pflag.BoolVarP(&showVersion, "", "v", false, "Show column for package versions")
	pflag.BoolVarP(&disableProgress, "no-progress", "", false, "Force suppress progress output")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")

	pflag.StringVar(&dateFilter, "date", "", "Filter packages installed on a specific date (YYYY-MM-DD)")
	pflag.StringVar(&sizeFilter, "size", "", "Filter packages by size, must be in quotes (e.g. \">20MB\", \"<1GB\")")
	pflag.StringVar(&nameFilter, "name", "", "Filter packages by name (or similar name)")
	pflag.StringVar(&sortBy, "sort", "date", "Sort packages by: 'date', 'alphabetical', 'size:desc', 'size:asc'")

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("Error parsing flags: %v", err)
	}

	if allPackages {
		count = 0
	}

	sizeFilterParsed, err := parseSizeFilter(sizeFilter)
	if err != nil {
		return Config{}, err
	}

	dateFilterParsed, err := parseDateFilter(dateFilter)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Count:             count,
		AllPackages:       allPackages,
		ShowHelp:          showHelp,
		ShowFullTimestamp: showFullTimestamp,
		DisableProgress:   disableProgress,
		ExplicitOnly:      explicitOnly,
		DependenciesOnly:  dependenciesOnly,
		DateFilter:        dateFilterParsed,
		SizeFilter:        sizeFilterParsed,
		NameFilter:        nameFilter,
		SortBy:            sortBy,
		OptionalColumns:   parseOptionalColumns(showVersion),
	}, nil
}

func parseDateFilter(dateFilterInput string) (DateFilter, error) {
	if dateFilterInput == "" {
		return DateFilter{}, fmt.Errorf("No date specified for --date flag")
	}

	dateParts := strings.Split(dateFilterInput, ":")

	var isExactMatch bool
	var startDate, endDate time.Time
	var err error

	switch {
	case len(dateParts) == 1:
		startDate, err = parseValidDate(dateParts[0])
		isExactMatch = true

	case dateParts[0] == "":
		endDate, err = parseValidDate(dateParts[1])

	case dateParts[1] == "":
		startDate, err = parseValidDate(dateParts[0])
		endDate = time.Now()

	default:
		startDate, err = parseValidDate(dateParts[0])
		if err == nil {
			endDate, err = parseValidDate(dateParts[1])
		}
	}

	if err != nil {
		return DateFilter{}, err
	}

	return DateFilter{
		startDate,
		endDate,
		isExactMatch,
	}, nil
}

func parseValidDate(dateInput string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", dateInput)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

func parseSizeFilter(sizeFilterInput string) (SizeFilter, error) {
	if sizeFilterInput != "" {
		sizeOperator, sizeInBytes, err := parseSizeInput(sizeFilterInput)
		if err != nil {
			return SizeFilter{}, fmt.Errorf("Invalid size filter: %v", err)
		}

		return SizeFilter{
			IsFilter:    true,
			SizeInBytes: sizeInBytes,
			Operator:    sizeOperator,
		}, nil
	}

	return SizeFilter{}, nil
}

func parseSizeInput(input string) (operator string, sizeInBytes int64, err error) {
	// matches for input of ">2KB" should be an array of [">2KB", ">", "2", "KB"]
	re := regexp.MustCompile(`(?i)^(<|>)?(\d+(?:\.\d+)?)(KB|MB|GB|B)?$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 1 {
		return "", 0, fmt.Errorf("invalid size filter format: %q", input)
	}

	operator = matches[1]

	if operator == "" {
		return "", 0, fmt.Errorf("invalid size operand: %q", operator)
	}

	sizeInBytes, err = parseSizeInBytes(matches[2], matches[3])
	if err != nil {
		return "", 0, err
	}

	return operator, sizeInBytes, nil
}

func parseSizeInBytes(valueInput string, unitInput string) (sizeInBytes int64, err error) {
	value, err := strconv.ParseFloat(valueInput, 64) // parseFloat for fractional input e.g. ">2.5KB"
	if err != nil {
		return sizeInBytes, fmt.Errorf("invalid size value")
	}

	unit := strings.ToUpper(unitInput)

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
		return sizeInBytes, fmt.Errorf("invalid size unit: %v", unit)
	}

	return sizeInBytes, nil
}

func parseOptionalColumns(showVersion bool) (optionalColumns []string) {
	if showVersion {
		optionalColumns = append(optionalColumns, "version")
	}

	return optionalColumns
}

func PrintHelp() {
	fmt.Println("Usage: yaylog [options]")

	fmt.Println("\nOptions:")
	pflag.PrintDefaults()

	fmt.Println("\nSorting Options:")
	fmt.Println("  --sort date          Sort packages by installation date (default)")
	fmt.Println("  --sort alphabetical  Sort packages alphabetically")
	fmt.Println("  --sort size:desc     Sort packages by size in descending order")
	fmt.Println("  --sort size:asc      Sort packages by size in ascending order")

	fmt.Println("\nFiltering Options:")
	fmt.Println("  --date YYYY-MM-DD    Filter packages installed on a specific date")
	fmt.Println("  --size <filter>      Filter packages by size with operators like:")
	fmt.Println("                         >20MB  (greater than 20 megabytes)")
	fmt.Println("                         <1GB   (less than 1 gigabyte)")

	fmt.Println("\nExamples:")
	fmt.Println("  yaylog --sort size:asc           # Sort by size (smallest to largest)")
	fmt.Println("  yaylog --size '>50MB'            # Show packages larger than 50MB")
	fmt.Println("  yaylog --date 2024-12-28         # Show packages installed on Dec 28, 2024")
	fmt.Println("  yaylog --size '<100MB' --sort alphabetical  # Filter packages smaller than 100MB and sort alphabetically")
}
