package config

import (
	"fmt"
	"os"
	"regexp"
	"yaylog/internal/consts"

	"github.com/spf13/pflag"
)

const (
	ReasonExplicit   = "explicit"
	ReasonDependency = "dependencies"
)

type Config struct {
	Count             int
	AllPackages       bool
	ShowHelp          bool
	OutputJson        bool
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	SortBy            string
	Fields            []consts.FieldType
	FilterQueries     map[consts.FieldType]string
}

type ConfigProvider interface {
	GetConfig() (Config, error)
}

type CliConfigProvider struct{}

func (c *CliConfigProvider) GetConfig() (Config, error) {
	cfg, err := ParseFlags(os.Args[1:])
	if err != nil {
		return Config{}, err
	}

	if cfg.ShowHelp {
		PrintHelp()
		os.Exit(0)
	}

	return cfg, nil
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

	var filterInputs []string
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
	pflag.StringArrayVarP(&filterInputs, "filter", "f", []string{}, "Apply multiple filters (e.g. --filter size=2KB:3KB --filter name=vim)")

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

	err := validateFlagCombinations(
		columnsInput,
		addColumnsInput,
		hasAllColumns,
		explicitOnly,
		dependenciesOnly,
	)
	if err != nil {
		return Config{}, err
	}

	if allPackages {
		count = 0
	}

	columnsParsed, err := parseColumns(columnsInput, addColumnsInput, hasAllColumns)
	if err != nil {
		return Config{}, err
	}

	filterQueries, err := parseFilterQueries(filterInputs)
	if err != nil {
		return Config{}, err
	}

	filterQueries = convertLegacyFilters(
		filterQueries,
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
		OutputJson:        outputJson,
		HasNoHeaders:      hasNoHeaders,
		ShowFullTimestamp: showFullTimestamp,
		DisableProgress:   disableProgress,
		SortBy:            sortBy,
		Fields:            columnsParsed,
		FilterQueries:     filterQueries,
	}

	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func parseFilterQueries(filterInputs []string) (map[consts.FieldType]string, error) {
	filterQueries := make(map[consts.FieldType]string)
	filterRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+)=(.+)$`)

	for _, input := range filterInputs {
		matches := filterRegex.FindStringSubmatch(input)
		if matches == nil {
			return nil, fmt.Errorf("Invalid filter format: %q. Must be in form field=value", input)
		}

		field, value := matches[1], matches[2]
		fieldType, exists := consts.FieldTypeLookup[field]
		if !exists {
			return nil, fmt.Errorf("Unknown filter field: %s", field)
		}

		if fieldType == consts.FieldReason && value != ReasonExplicit && value != ReasonDependency {
			return nil, fmt.Errorf("Invalid reason filter value: %s. Allowed values are 'explicit' or 'dependency'", value)
		}

		filterQueries[fieldType] = value
	}

	return filterQueries, nil
}

func convertLegacyFilters(
	filterQueries map[consts.FieldType]string,
	dateFilter string,
	nameFilter string,
	sizeFilter string,
	requiredByFilter string,
	explicitOnly bool,
	dependenciesOnly bool,
) map[consts.FieldType]string {
	if dateFilter != "" {
		filterQueries[consts.FieldDate] = dateFilter
	}

	if nameFilter != "" {
		filterQueries[consts.FieldName] = nameFilter
	}

	if sizeFilter != "" {
		filterQueries[consts.FieldSize] = sizeFilter
	}

	if requiredByFilter != "" {
		filterQueries[consts.FieldRequiredBy] = requiredByFilter
	}

	if explicitOnly {
		filterQueries[consts.FieldReason] = ReasonExplicit
	}

	if dependenciesOnly {
		filterQueries[consts.FieldReason] = ReasonDependency
	}

	return filterQueries
}
