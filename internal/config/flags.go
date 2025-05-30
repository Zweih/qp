package config

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/preprocess"
	"qp/internal/syntax"
	"strings"

	"github.com/spf13/pflag"
)

const internalCacheWorker = "internal-cache-worker"

func ParseFlags(args []string) (Config, error) {
	var flagCfg Config
	var legacyFieldInput, legacyAddFieldInput string
	var legacySortInput string
	var legacyFilterInputs []string
	var legacyCount int

	var explicitOnly, dependenciesOnly, legacyHasAllFields, legacyAllPackages bool
	var dateFilter, sizeFilter, nameFilter, requiredByFilter string

	registerCommonFlags(&flagCfg)
	registerLegacyFlags(
		&flagCfg, &legacyFieldInput, &legacyAddFieldInput,
		&legacySortInput, &legacyFilterInputs,
		&explicitOnly, &dependenciesOnly, &legacyHasAllFields,
		&dateFilter, &sizeFilter, &nameFilter, &requiredByFilter,
		&legacyAllPackages, &legacyCount,
	)

	markHiddenFlags()

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %v", err)
	}

	// exit asap, we don't need user syntax
	if flagCfg.CacheOnly != "" {
		return flagCfg, nil
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
				legacyAllPackages, legacyCount,
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
		if expanded, exists := preprocess.ShorthandMap[lower]; exists {
			lower = expanded
		}

		switch lower {
		case consts.CmdSelect, consts.CmdWhere,
			consts.CmdOrder, consts.CmdLimit:
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
	dst.QueryExpr = src.QueryExpr
	dst.Limit = src.Limit
	dst.LimitMode = src.LimitMode
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
		"select-all",
		"select",
		"select-add",
		"order",
		"where",
		"limit",
		"all",
	}

	for _, flag := range hiddenFlags {
		_ = pflag.CommandLine.MarkHidden(flag)
		pflag.CommandLine.MarkDeprecated(flag, "please use command syntax instead. Run `man qp` for more info. Flag syntax will be removed in the future.")
	}
}

func registerCommonFlags(cfg *Config) {
	pflag.BoolVar(&cfg.HasNoHeaders, "no-headers", false, "Hide headers")
	pflag.StringVar(&cfg.OutputFormat, "output", "table", "Output format: \"table\", \"json\", or \"kv\" (key-value)")
	pflag.BoolVarP(&cfg.ShowHelp, "help", "h", false, "Show help")
	pflag.BoolVar(&cfg.ShowVersion, "version", false, "Show version")
	pflag.BoolVar(&cfg.ShowFullTimestamp, "full-timestamp", false, "Show full timestamp")
	pflag.BoolVar(&cfg.DisableProgress, "no-progress", false, "Disable progress bar")
	pflag.BoolVar(&cfg.NoCache, "no-cache", false, "Disable cache")
	pflag.BoolVar(&cfg.RegenCache, "regen-cache", false, "Force fresh cache")
	pflag.StringVar(&cfg.CacheOnly, "cache-only", "", "Update cache only and nothing else. Specify origin ('pacman', 'brew', 'deb') or 'all'")
	pflag.StringVar(&cfg.CacheWorker, internalCacheWorker, "", "Internal flag for background cache operations - do not use directly")

	_ = pflag.CommandLine.MarkHidden(internalCacheWorker)
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
	allPackages *bool,
	count *int,
) {
	pflag.BoolVarP(hasAllFields, "select-all", "A", false, "Show all fields")
	pflag.StringVarP(fieldInput, "select", "s", "", "Select fields")
	pflag.StringVarP(addFieldInput, "select-add", "S", "", "Add to selected fields")
	pflag.StringVarP(sortInput, "order", "O", "date", "Sort order")
	pflag.StringArrayVarP(filterInputs, "where", "w", []string{}, "Filter")
	pflag.IntVarP(count, "limit", "l", 20, "Number of packages to show")
	pflag.BoolVarP(allPackages, "all", "a", false, "Show all packages")

	// hidden legacy flags
	pflag.IntVarP(&cfg.Limit, "number", "n", 20, "")
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

	var legacyJSON bool
	pflag.BoolVar(&legacyJSON, "json", false, "")
	_ = pflag.CommandLine.MarkHidden("json")
	_ = pflag.CommandLine.MarkDeprecated("json", "use \"--output json\" instead.")
	if legacyJSON {
		cfg.OutputFormat = consts.OutputJSON
	}
}
