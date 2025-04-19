package config

import (
	"fmt"
	"qp/internal/syntax"
	"strings"

	"github.com/spf13/pflag"
)

func ParseFlags(args []string) (Config, error) {
	var flagCfg Config
	var legacyFieldInput, legacyAddFieldInput string
	var legacySortInput string
	var legacyFilterInputs []string

	var explicitOnly, dependenciesOnly, legacyHasAllFields bool
	var dateFilter, sizeFilter, nameFilter, requiredByFilter string

	registerCommonFlags(&flagCfg)
	registerLegacyFlags(
		&flagCfg, &legacyFieldInput, &legacyAddFieldInput,
		&legacySortInput, &legacyFilterInputs,
		&explicitOnly, &dependenciesOnly, &legacyHasAllFields,
		&dateFilter, &sizeFilter, &nameFilter, &requiredByFilter,
	)

	markHiddenFlags()

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %v", err)
	}

	if flagCfg.AllPackages {
		flagCfg.Count = 0
	}

	remainingArgs := pflag.Args()
	newSyntaxParser := func() (syntax.ParsedInput, error) {
		return syntax.ParseSyntax(remainingArgs)
	}

	parser := newSyntaxParser
	if !isNewSyntax(args) {
		if err := validateFlagCombinations(legacyFieldInput, legacyAddFieldInput, legacyHasAllFields, explicitOnly, dependenciesOnly); err != nil {
			return Config{}, err
		}

		parser = func() (syntax.ParsedInput, error) {
			return ParseLegacyConfig(
				legacyFieldInput,
				legacyAddFieldInput,
				legacyHasAllFields,
				legacySortInput,
				legacyFilterInputs,
				dateFilter, nameFilter, sizeFilter, requiredByFilter,
				explicitOnly, dependenciesOnly,
			)
		}
	}

	parsedInput, err := parser()
	if err != nil {
		return Config{}, err
	}

	mergeTopLevelOptions(&flagCfg, &parsedInput)
	return flagCfg, nil
}

func isNewSyntax(args []string) bool {
	for _, arg := range args {
		lower := strings.ToLower(arg)
		if expanded, exists := syntax.ShorthandMap[lower]; exists {
			lower = expanded
		}

		switch lower {
		case syntax.CmdSelect, syntax.CmdWhere, syntax.CmdOrder:
			return true
		default:
			return false
		}
	}

	return false
}

func mergeTopLevelOptions(dst *Config, src *syntax.ParsedInput) {
	dst.SortOption = src.SortOption
	dst.Fields = src.Fields
	dst.FieldQueries = src.FieldQueries
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

func registerCommonFlags(cfg *Config) {
	pflag.IntVarP(&cfg.Count, "limit", "l", 20, "Number of packages to show")
	pflag.BoolVarP(&cfg.AllPackages, "all", "a", false, "Show all packages")
	pflag.BoolVar(&cfg.HasNoHeaders, "no-headers", false, "Hide headers")
	pflag.BoolVar(&cfg.OutputJson, "json", false, "Output in JSON format")
	pflag.BoolVar(&cfg.ShowHelp, "help", false, "Show help")
	pflag.BoolVar(&cfg.ShowVersion, "version", false, "Show version")
	pflag.BoolVar(&cfg.ShowFullTimestamp, "full-timestamp", false, "Show full timestamp")
	pflag.BoolVar(&cfg.DisableProgress, "no-progress", false, "Disable progress bar")
	pflag.BoolVar(&cfg.NoCache, "no-cache", false, "Disable cache")
	pflag.BoolVar(&cfg.RegenCache, "regen-cache", false, "Force fresh cache")
}

func registerLegacyFlags(
	cfg *Config,
	fieldInput *string,
	addFieldInput *string,
	sortInput *string,
	filterInputs *[]string,
	explicitOnly *bool,
	dependenciesOnly *bool,
	hasAllFields *bool,
	dateFilter, sizeFilter, nameFilter, requiredByFilter *string,
) {
	pflag.BoolVarP(hasAllFields, "select-all", "A", false, "Show all fields")
	pflag.StringVarP(fieldInput, "select", "s", "", "Select fields")
	pflag.StringVarP(addFieldInput, "select-add", "S", "", "Add to selected fields")
	pflag.StringVarP(sortInput, "order", "O", "date", "Sort order")
	pflag.StringArrayVarP(filterInputs, "where", "w", []string{}, "Filter")

	// hidden legacy flags
	pflag.IntVarP(&cfg.Count, "number", "n", 20, "")
	pflag.StringVar(fieldInput, "columns", "", "")
	pflag.StringVar(addFieldInput, "add-columns", "", "")
	pflag.BoolVar(hasAllFields, "all-columns", false, "")
	pflag.StringArrayVarP(filterInputs, "filter", "f", []string{}, "")
	pflag.StringVar(sortInput, "sort", "date", "")

	pflag.BoolVarP(explicitOnly, "explicit", "e", false, "")
	pflag.BoolVarP(dependenciesOnly, "dependencies", "d", false, "")
	pflag.StringVar(dateFilter, "date", "", "")
	pflag.StringVar(sizeFilter, "size", "", "")
	pflag.StringVar(nameFilter, "name", "", "")
	pflag.StringVar(requiredByFilter, "required-by", "", "")
}
