package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Config struct {
	Count             int
	AllPackages       bool
	ShowHelp          bool
	OutputJson        bool
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	ExplicitOnly      bool
	DependenciesOnly  bool
	DateFilter        DateFilter
	SizeFilter        SizeFilter
	NameFilter        string
	RequiredByFilter  string
	SortBy            string
	ColumnNames       []string
}

func ParseFlags(args []string) (Config, error) {
	var count int

	var allPackages bool
	var hasAllColumns bool
	var showHelp bool
	var outputJson bool
	var hasNoHeaders bool
	var showFullTimestamp bool
	var disableProgress bool
	var explicitOnly bool
	var dependenciesOnly bool

	var dateFilter string
	var sizeFilter string
	var nameFilter string
	var requiredByFilter string
	var sortBy string
	var columnsInput string
	var addColumnsInput string

	pflag.CommandLine.SortFlags = false // keeps the help output in the order we define below

	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")
	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -n)")

	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")

	pflag.StringVar(&dateFilter, "date", "", "Filter packages by installation date. Supports exact dates (YYYY-MM-DD), ranges (YYYY-MM-DD:YYYY-MM-DD), and open-ended filters (:YYYY-MM-DD or YYYY-MM-DD:).")
	pflag.StringVar(&sizeFilter, "size", "", "Filter packages by size. Supports ranges (e.g., 10MB:20GB), exact matches (e.g., 5MB), and open-ended values (e.g., :2GB or 500KB:)")
	pflag.StringVar(&nameFilter, "name", "", "Filter packages by name (or similar name)")

	pflag.StringVar(&sortBy, "sort", "date", "Sort packages by: 'date', 'alphabetical', 'size:desc', 'size:asc'")

	pflag.BoolVarP(&hasNoHeaders, "no-headers", "", false, "Hide headers for columns (useful for scripts/automation)")

	pflag.BoolVarP(&hasAllColumns, "all-columns", "", false, "Show all available columns/fields in the output (overrides defaults)")
	pflag.StringVar(&columnsInput, "columns", "", "Comma-separated list of columns to display (overrides defaults)")
	pflag.StringVar(&addColumnsInput, "add-columns", "", "Comma-separated list of columns to add to defaults")

	pflag.BoolVarP(&showFullTimestamp, "full-timestamp", "", false, "Show full timestamp instead of just the date")
	pflag.BoolVarP(&outputJson, "json", "", false, "Output results in JSON format")
	pflag.BoolVarP(&disableProgress, "no-progress", "", false, "Force suppress progress output")
	pflag.StringVar(&requiredByFilter, "required-by", "", "Show only packages that are required by the specified package")

	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("Error parsing flags: %v", err)
	}

	err := validateFlagCombinations(columnsInput, addColumnsInput, hasAllColumns, explicitOnly, dependenciesOnly)
	if err != nil {
		return Config{}, err
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

	columnsParsed, err := parseColumns(columnsInput, addColumnsInput, hasAllColumns)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Count:             count,
		AllPackages:       allPackages,
		ShowHelp:          showHelp,
		OutputJson:        outputJson,
		HasNoHeaders:      hasNoHeaders,
		ShowFullTimestamp: showFullTimestamp,
		DisableProgress:   disableProgress,
		ExplicitOnly:      explicitOnly,
		DependenciesOnly:  dependenciesOnly,
		DateFilter:        dateFilterParsed,
		SizeFilter:        sizeFilterParsed,
		NameFilter:        nameFilter,
		RequiredByFilter:  requiredByFilter,
		SortBy:            sortBy,
		ColumnNames:       columnsParsed,
	}, nil
}
