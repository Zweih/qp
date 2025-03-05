package config

import (
	"fmt"
	"strings"
	"yaylog/internal/consts"

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

	pflag.CommandLine.SortFlags = false

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

func parseColumns(columnsInput string, addColumnsInput string, hasAllColumns bool) ([]string, error) {
	if columnsInput != "" && addColumnsInput != "" {
		return nil, fmt.Errorf("cannot use --columns and --add-columns together. Use --columns to fully define the columns you want")
	}

	if hasAllColumns {
		return consts.ValidColumns, nil
	}

	var specifiedColumnsRaw string
	var columns []string

	switch {
	case columnsInput != "":
		specifiedColumnsRaw = columnsInput
	case addColumnsInput != "":
		specifiedColumnsRaw = addColumnsInput
		fallthrough
	default:
		columns = consts.DefaultColumns
	}

	specifiedColumns, err := validateColumns(strings.ToLower(specifiedColumnsRaw))
	if err != nil {
		return nil, err
	}

	columns = append(columns, specifiedColumns...)

	if len(columns) < 1 {
		return nil, fmt.Errorf("no columns selected: use --columns to specify at least one column")
	}

	return columns, nil
}

func validateColumns(columnInput string) ([]string, error) {
	if columnInput == "" {
		return []string{}, nil
	}

	validColumnsSet := map[string]bool{}
	for _, columnName := range consts.ValidColumns {
		validColumnsSet[columnName] = true
	}

	var columns []string

	for _, column := range strings.Split(columnInput, ",") {
		cleanColumn := strings.TrimSpace(column)

		if !validColumnsSet[strings.TrimSpace(column)] {
			return nil, fmt.Errorf("%s is not a valid column", cleanColumn)
		}

		columns = append(columns, cleanColumn)
	}

	return columns, nil
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
	fmt.Println("  --date <filter>      Filter packages by installation date. Supports:")
	fmt.Println("                         YYYY-MM-DD       (exact date match)")
	fmt.Println("                         YYYY-MM-DD:      (installed on or after the date)")
	fmt.Println("                         :YYYY-MM-DD      (installed up to the date)")
	fmt.Println("                         YYYY-MM-DD:YYYY-MM-DD  (installed within a date range)")

	fmt.Println("  --size <filter>      Filter packages by size on disk. Supports:")
	fmt.Println("                         10MB       (exactly 10MB)")
	fmt.Println("                         5GB:       (5GB and larger)")
	fmt.Println("                         :20KB      (up to 20KB)")
	fmt.Println("                         1.5MB:2GB  (between 1.5MB and 2GB)")

	fmt.Println("  --name <search-term> Filter packages by name (substring match)")
	fmt.Println("                         Example: 'gtk' matches 'gtk3', 'libgtk', etc")

	fmt.Println("  --required-by <name> Show only packages that are required by the specified package")
	fmt.Println("                         Example: 'yaylog --required-by firefox' lists packages that firefox depends on")

	fmt.Println("\nColumn Options:")
	fmt.Println("  --columns <list>     Comma-separated list of columns to display (overrides defaults)")
	fmt.Println("  --add-columns <list> Comma-separated list of columns to add to defaults")
	fmt.Println("  --all-columns        Display all available columns")
	fmt.Println("  --no-headers         Omit column headers in output (useful for scripts and automation)")

	fmt.Println("\nAvailable Columns:")
	fmt.Println("  date         - Installation date of the package")
	fmt.Println("  name         - Package name")
	fmt.Println("  reason       - Installation reason (explicit/dependency)")
	fmt.Println("  size         - Package size on disk")
	fmt.Println("  version      - Installed package version")
	fmt.Println("  depends      - List of dependencies (output can be long)")
	fmt.Println("  required-by  - List of packages required by the package and are dependent on it (output can be long)")
	fmt.Println("  provides     - List of alternative package names or shared libraries provided by package (output can be long)")

	fmt.Println("\nCaveat:")
	fmt.Println("  The 'depends', 'provides', and 'required-by' columns output can be lengthy. It's recommended to use `less` for better readability:")
	fmt.Println("  yaylog --columns name,depends | less")

	fmt.Println("\nExamples:")
	fmt.Println("  yaylog --size 50MB --date 2024-12-28             # Show 50MB packages installed on Dec 28, 2024")
	fmt.Println("  yaylog --size 100MB: --date :2024-06-30          # Show packages >100MB installed up to June 30, 2024")
	fmt.Println("  yaylog --size 10MB:1GB --date 2023-01-01:        # Packages 10MB-1GB installed after Jan 1, 2023")
	fmt.Println("  yaylog --sort size:desc --date 2024-01-01:       # Sort by largest, installed on/after Jan 1, 2024")
	fmt.Println("  yaylog --size :50MB --sort alphabetical          # Sort small packages alphabetically")
	fmt.Println("  yaylog --name python                             # Show installed packages containing 'python'")
	fmt.Println("  yaylog --name gtk --size 5MB: --date 2023-01-01: # Packages with 'gtk', >5MB, installed after Jan 1, 2023")
	fmt.Println("  yaylog --columns name,version,size               # Show packages with name, version, and size")
	fmt.Println("  yaylog --columns name,depends | less             # Show package names and dependencies with less for readability")
}
