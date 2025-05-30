package config

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/syntax"

	"github.com/spf13/pflag"
)

const internalCacheWorker = "internal-cache-worker"

func ParseFlags(args []string) (Config, error) {
	var cfg Config
	registerCommonFlags(&cfg)

	if err := pflag.CommandLine.Parse(args); err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %v", err)
	}

	// exit asap, we don't need user syntax
	if cfg.CacheOnly != "" {
		return cfg, nil
	}

	remainingArgs := pflag.Args()
	parsedInput, err := syntax.ParseSyntax(remainingArgs)
	if err != nil {
		return Config{}, err
	}

	mergeTopLevelOptions(&cfg, &parsedInput)
	return cfg, nil
}

func mergeTopLevelOptions(dst *Config, src *syntax.ParsedInput) {
	dst.SortOption = src.SortOption
	dst.Fields = src.Fields
	dst.FieldQueries = src.FieldQueries
	dst.QueryExpr = src.QueryExpr
	dst.Limit = src.Limit
	dst.LimitMode = src.LimitMode
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
	pflag.StringVar(&cfg.CacheOnly, "cache-only", "", "Update cache only for specifed origin.")
	pflag.StringVar(&cfg.CacheWorker, internalCacheWorker, "", "Internal flag for background cache operations - do not use directly")

	_ = pflag.CommandLine.MarkHidden(internalCacheWorker)

	var legacyJSON bool
	pflag.BoolVar(&legacyJSON, "json", false, "")
	_ = pflag.CommandLine.MarkHidden("json")
	_ = pflag.CommandLine.MarkDeprecated("json", "use \"--output json\" instead.")
	if legacyJSON {
		cfg.OutputFormat = consts.OutputJSON
	}
}
