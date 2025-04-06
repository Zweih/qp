package config

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"strings"

	"github.com/spf13/pflag"
)

const (
	ReasonExplicit   = "explicit"
	ReasonDependency = "dependency"
)

type Config struct {
	Count             int
	AllPackages       bool
	ShowHelp          bool
	OutputJson        bool
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	SortOption        SortOption
	Fields            []consts.FieldType
	FieldQueries      FieldQueries
}

type (
	FieldQueries    map[consts.FieldType]SubfieldQueries
	SubfieldQueries map[consts.SubfieldType]string
)

type SortOption struct {
	Field consts.FieldType
	Asc   bool
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
	var hasAllFields bool
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
	var sortInput string
	var fieldInput string
	var addFieldInput string

	pflag.CommandLine.SortFlags = false // keeps the help output in the order we define below

	pflag.IntVarP(&count, "limit", "l", 20, "Number of packages to show")
	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -l)")

	pflag.StringArrayVarP(&filterInputs, "where", "w", []string{}, "Query by one or more fields (e.g. -w size=2KB:3KB -w name=vim)")
	pflag.StringVarP(&sortInput, "order", "O", "date", "Order results by field")

	pflag.BoolVarP(&hasNoHeaders, "no-headers", "", false, "Hide headers for table ouput (useful for scripts/automation)")
	pflag.BoolVarP(&hasAllFields, "select-all", "A", false, "Display all available fields")
	pflag.StringVarP(&fieldInput, "select", "s", "", "Select exact fields to display")
	pflag.StringVarP(&addFieldInput, "select-add", "S", "", "Add fields to the default output")

	pflag.BoolVarP(&showFullTimestamp, "full-timestamp", "", false, "Show full timestamp instead of just the date")
	pflag.BoolVarP(&outputJson, "json", "", false, "Output results in JSON format")
	pflag.BoolVarP(&disableProgress, "no-progress", "", false, "Force suppress progress output")

	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")

	// deprecated legacy flags, hidden but still functioning
	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")
	pflag.StringArrayVarP(&filterInputs, "filter", "f", []string{}, "Apply multiple filters (e.g. --filter size=2KB:3KB --filter name=vim)")
	pflag.StringVar(&sortInput, "sort", "date", "Sort packages by: 'date', 'alphabetical', 'size:desc', 'size:asc'")
	pflag.BoolVarP(&hasAllFields, "all-columns", "", false, "Show all available columns/fields in the output (overrides defaults)")
	pflag.StringVar(&fieldInput, "columns", "", "Comma-separated list of columns to display (overrides defaults)")
	pflag.StringVar(&addFieldInput, "add-columns", "", "Comma-separated list of columns to add to defaults")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")
	pflag.StringVar(&dateFilter, "date", "", "Filter packages by installation date. Supports exact dates (YYYY-MM-DD), ranges (YYYY-MM-DD:YYYY-MM-DD), and open-ended filters (:YYYY-MM-DD or YYYY-MM-DD:).")
	pflag.StringVar(&sizeFilter, "size", "", "Filter packages by size. Supports ranges (e.g., 10MB:20GB), exact matches (e.g., 5MB), and open-ended values (e.g., :2GB or 500KB:)")
	pflag.StringVar(&nameFilter, "name", "", "Filter packages by name (or similar name)")
	pflag.StringVar(&requiredByFilter, "required-by", "", "Show only packages that are required by the specified package")

	_ = pflag.CommandLine.MarkHidden("number")
	_ = pflag.CommandLine.MarkHidden("filter")
	_ = pflag.CommandLine.MarkHidden("sort")
	_ = pflag.CommandLine.MarkHidden("all-columns")
	_ = pflag.CommandLine.MarkHidden("columns")
	_ = pflag.CommandLine.MarkHidden("add-columns")
	_ = pflag.CommandLine.MarkHidden("explicit")
	_ = pflag.CommandLine.MarkHidden("dependencies")
	_ = pflag.CommandLine.MarkHidden("date")
	_ = pflag.CommandLine.MarkHidden("size")
	_ = pflag.CommandLine.MarkHidden("name")
	_ = pflag.CommandLine.MarkHidden("required-by")

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("Error parsing flags: %v", err)
	}

	err := validateFlagCombinations(
		fieldInput,
		addFieldInput,
		hasAllFields,
		explicitOnly,
		dependenciesOnly,
	)
	if err != nil {
		return Config{}, err
	}

	if allPackages {
		count = 0
	}

	fieldsParsed, err := parseFields(fieldInput, addFieldInput, hasAllFields)
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
		OutputJson:        outputJson,
		HasNoHeaders:      hasNoHeaders,
		ShowFullTimestamp: showFullTimestamp,
		DisableProgress:   disableProgress,
		SortOption:        sortOption,
		Fields:            fieldsParsed,
		FieldQueries:      fieldQueries,
	}

	return cfg, nil
}

func parseSortOption(sortInput string) (SortOption, error) {
	parts := strings.Split(sortInput, ":")
	fieldKey := strings.ToLower(parts[0])
	fieldType, exists := consts.FieldTypeLookup[fieldKey]
	if !exists {
		return SortOption{}, fmt.Errorf("invalid sort field: %s", fieldKey)
	}

	asc := true
	if len(parts) > 1 {
		switch parts[1] {
		case "desc":
			asc = false
		case "asc":
			asc = true
		default:
			return SortOption{}, fmt.Errorf("invalid sort direction: %s", parts[1])
		}
	}

	return SortOption{
		Field: fieldType,
		Asc:   asc,
	}, nil
}

func parseQueries(queryInputs []string) (FieldQueries, error) {
	queries := make(FieldQueries)

	for _, queryInput := range queryInputs {
		fieldPart, value, err := parseQueryInput(queryInput)
		if err != nil {
			return nil, err
		}

		field, subfield, err := parseFieldPart(fieldPart)
		if err != nil {
			return nil, err
		}

		addQuery(queries, field, subfield, value)
	}

	return queries, nil
}

func parseQueryInput(input string) (string, string, error) {
	queryParts := strings.SplitN(input, "=", 2)
	if len(queryParts) != 2 {
		return "",
			"",
			fmt.Errorf("invalid query format: %s. Must be in form fireld.subfield=value", input)
	}

	return queryParts[0], queryParts[1], nil
}

func parseFieldPart(fieldPart string) (consts.FieldType, consts.SubfieldType, error) {
	fieldParts := strings.SplitN(fieldPart, ".", 2)
	fieldName := fieldParts[0]
	field, exists := consts.FieldTypeLookup[fieldName]
	if !exists {
		return 0, 0, fmt.Errorf("unknown query field: %s", fieldName)
	}

	subfieldName := ""
	if len(fieldParts) == 2 {
		subfieldName = fieldParts[1]
	}

	subfield, exists := consts.SubfieldTypeLookup[subfieldName]
	if !exists {
		return 0, 0, fmt.Errorf("unknown query subfield: %s", subfieldName)
	}

	return field, subfield, nil
}

func addQuery(
	queries FieldQueries,
	field consts.FieldType,
	subfield consts.SubfieldType,
	value string,
) {
	subfields, exists := queries[field]
	if !exists {
		subfields = make(SubfieldQueries)
	}

	subfields[subfield] = value
	queries[field] = subfields
}

func addLegacyQuery(
	queries FieldQueries,
	field consts.FieldType,
	subfield consts.SubfieldType,
	value string,
) {
	if value == "" {
		return
	}

	subfields := queries[field]
	if subfields == nil {
		subfields = make(SubfieldQueries)
	}

	subfields[subfield] = value
	queries[field] = subfields
}

func convertLegacyQueries(
	queries FieldQueries,
	dateFilter string,
	nameFilter string,
	sizeFilter string,
	requiredByFilter string,
	explicitOnly bool,
	dependenciesOnly bool,
) FieldQueries {
	addLegacyQuery(queries, consts.FieldDate, consts.SubfieldTarget, dateFilter)
	addLegacyQuery(queries, consts.FieldName, consts.SubfieldTarget, nameFilter)
	addLegacyQuery(queries, consts.FieldSize, consts.SubfieldTarget, sizeFilter)
	addLegacyQuery(queries, consts.FieldRequiredBy, consts.SubfieldTarget, requiredByFilter)

	if explicitOnly {
		addLegacyQuery(queries, consts.FieldReason, consts.SubfieldTarget, ReasonExplicit)
	}

	if dependenciesOnly {
		addLegacyQuery(queries, consts.FieldReason, consts.SubfieldTarget, ReasonDependency)
	}

	return queries
}
