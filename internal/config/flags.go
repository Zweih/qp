package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

func ParseFlags(args []string) (Config, error) {
	var count int
	var allPackages bool
	var hasAllFields bool
	var showHelp bool
	var showVersion bool
	var outputJson bool
	var hasNoHeaders bool
	var showFullTimestamp bool
	var disableProgress bool
	var noCache bool
	var regenCache bool

	var filterInputs []string
	var sortInput string
	var fieldInput string
	var addFieldInput string

	// legacy hidden flags
	var explicitOnly bool
	var dependenciesOnly bool
	var dateFilter string
	var sizeFilter string
	var nameFilter string
	var requiredByFilter string

	pflag.CommandLine.SortFlags = false

	pflag.IntVarP(&count, "limit", "l", 20, "Number of packages to show")
	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -l)")
	pflag.StringArrayVarP(&filterInputs, "where", "w", []string{}, "Query by one or more fields (e.g. -w size=2KB:3KB -w name=vim)")
	pflag.StringVarP(&sortInput, "order", "O", "date", "Order results by field")

	pflag.BoolVar(&hasNoHeaders, "no-headers", false, "Hide headers for table output (useful for scripts/automation)")
	pflag.BoolVarP(&hasAllFields, "select-all", "A", false, "Display all available fields")
	pflag.StringVarP(&fieldInput, "select", "s", "", "Select exact fields to display")
	pflag.StringVarP(&addFieldInput, "select-add", "S", "", "Add fields to the default output")

	pflag.BoolVar(&showFullTimestamp, "full-timestamp", false, "Show full timestamp instead of just the date")
	pflag.BoolVar(&outputJson, "json", false, "Output results in JSON format")
	pflag.BoolVar(&disableProgress, "no-progress", false, "Force suppress progress output")
	pflag.BoolVar(&noCache, "no-cache", false, "Disable cache loading/saving and force fresh package data loading")
	pflag.BoolVar(&regenCache, "regen-cache", false, "Disable cache loading, force fresh package data loading, and save fresh cache")

	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help information")
	pflag.BoolVar(&showVersion, "version", false, "Show version, author, and license information")

	// legacy hidden flags (still functional but hidden)
	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")
	pflag.StringArrayVarP(&filterInputs, "filter", "f", []string{}, "Apply multiple filters (e.g. --filter size=2KB:3KB --filter name=vim)")
	pflag.StringVar(&sortInput, "sort", "date", "Sort packages by: 'date', 'alphabetical', 'size:desc', 'size:asc'")
	pflag.BoolVar(&hasAllFields, "all-columns", false, "Show all available columns/fields in the output (overrides defaults)")
	pflag.StringVar(&fieldInput, "columns", "", "Comma-separated list of columns to display (overrides defaults)")
	pflag.StringVar(&addFieldInput, "add-columns", "", "Comma-separated list of columns to add to defaults")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")
	pflag.StringVar(&dateFilter, "date", "", "Filter packages by installation date")
	pflag.StringVar(&sizeFilter, "size", "", "Filter packages by size")
	pflag.StringVar(&nameFilter, "name", "", "Filter packages by name")
	pflag.StringVar(&requiredByFilter, "required-by", "", "Show only packages required by the specified package")

	markHiddenFlags()

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %v", err)
	}

	if err := validateFlagCombinations(fieldInput, addFieldInput, hasAllFields, explicitOnly, dependenciesOnly); err != nil {
		return Config{}, err
	}

	if allPackages {
		count = 0
	}

	fieldsParsed, err := parseSelection(fieldInput, addFieldInput, hasAllFields)
	if err != nil {
		return Config{}, err
	}

	sortOption, err := parseSortOption(sortInput)
	if err != nil {
		return Config{}, err
	}

	fieldQueries, err := parseQueries(filterInputs)
	if err != nil {
		return Config{}, err
	}

	fieldQueries = convertLegacyQueries(
		fieldQueries,
		dateFilter,
		nameFilter,
		sizeFilter,
		requiredByFilter,
		explicitOnly,
		dependenciesOnly,
	)

	cfg := Config{
		Count:             count,
		AllPackages:       allPackages,
		ShowHelp:          showHelp,
		ShowVersion:       showVersion,
		OutputJson:        outputJson,
		HasNoHeaders:      hasNoHeaders,
		ShowFullTimestamp: showFullTimestamp,
		NoCache:           noCache,
		RegenCache:        regenCache,
		DisableProgress:   disableProgress,
		SortOption:        sortOption,
		Fields:            fieldsParsed,
		FieldQueries:      fieldQueries,
	}

	return cfg, nil
}

func markHiddenFlags() {
	hiddenFlags := []string{
		"number",
		"filter",
		"sort",
		"all-columns",
		"columns",
		"add-columns",
		"explicit",
		"dependencies",
		"date",
		"size",
		"name",
		"required-by",
	}

	for _, flag := range hiddenFlags {
		_ = pflag.CommandLine.MarkHidden(flag)
	}
}
