package config

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yaylog/internal/consts"

	"github.com/spf13/pflag"
)

type SizeFilter struct {
	StartSize    int64
	EndSize      int64
	IsExactMatch bool
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
	NoDefaults        bool
	DateFilter        DateFilter
	SizeFilter        SizeFilter
	NameFilter        string
	SortBy            string
	ColumnNames       []string
}

func ParseFlags(args []string) (Config, error) {
	var count int
	var allPackages bool
	var showHelp bool
	var showFullTimestamp bool
	var disableProgress bool
	var explicitOnly bool
	var dependenciesOnly bool
	var dateFilter string
	var sizeFilter string
	var nameFilter string
	var sortBy string
	var columnsInput string
	var addColumnsInput string

	pflag.IntVarP(&count, "number", "n", 20, "Number of packages to show")

	pflag.BoolVarP(&allPackages, "all", "a", false, "Show all packages (ignores -n)")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Display help")
	pflag.BoolVarP(&showFullTimestamp, "full-timestamp", "", false, "Show full timestamp instead of just the date")
	pflag.BoolVarP(&disableProgress, "no-progress", "", false, "Force suppress progress output")
	pflag.BoolVarP(&explicitOnly, "explicit", "e", false, "Show only explicitly installed packages")
	pflag.BoolVarP(&dependenciesOnly, "dependencies", "d", false, "Show only packages installed as dependencies")

	pflag.StringVar(&dateFilter, "date", "", "Filter packages by installation date. Supports exact dates (YYYY-MM-DD), ranges (YYYY-MM-DD:YYYY-MM-DD), and open-ended filters (:YYYY-MM-DD or YYYY-MM-DD:).")
	pflag.StringVar(&sizeFilter, "size", "", "Filter packages by size. Supports ranges (e.g., 10MB:20GB), exact matches (e.g., 5MB), and open-ended values (e.g., :2GB or 500KB:)")
	pflag.StringVar(&nameFilter, "name", "", "Filter packages by name (or similar name)")
	pflag.StringVar(&sortBy, "sort", "date", "Sort packages by: 'date', 'alphabetical', 'size:desc', 'size:asc'")
	pflag.StringVar(&columnsInput, "columns", "", "Comma-separated list of columns to display (overrides defaults)")
	pflag.StringVar(&addColumnsInput, "add-columns", "", "Comma-separated list of columns to add to defaults")

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

	columnsParsed, err := parseColumns(columnsInput, addColumnsInput)
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
		ColumnNames:       columnsParsed,
	}, nil
}

func parseDateFilter(dateFilterInput string) (DateFilter, error) {
	if dateFilterInput == "" {
		return DateFilter{}, nil
	}

	if dateFilterInput == ":" {
		return DateFilter{}, fmt.Errorf("invalid date filter: ':' must be accompanied by a date")
	}

	pattern := `^(\d{4}-\d{2}-\d{2})?(?::(\d{4}-\d{2}-\d{2})?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(dateFilterInput)
	isExactMatch := !strings.Contains(dateFilterInput, ":")

	if matches == nil {
		return DateFilter{}, fmt.Errorf("invalid date filter format: %q", dateFilterInput)
	}

	startDate, err := parseDateMatch(matches[1], time.Time{})
	if err != nil {
		return DateFilter{}, err
	}

	endDate, err := parseDateMatch(matches[2], time.Now())
	if err != nil {
		return DateFilter{}, err
	}

	return DateFilter{
		startDate,
		endDate,
		isExactMatch,
	}, nil
}

func parseDateMatch(dateInput string, defaultDate time.Time) (time.Time, error) {
	if dateInput == "" {
		return defaultDate, nil
	}

	return parseValidDate(dateInput)
}

func parseValidDate(dateInput string) (time.Time, error) {
	parsedDate, err := time.Parse(consts.DateOnlyFormat, dateInput)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

func parseSizeFilter(sizeFilterInput string) (SizeFilter, error) {
	if sizeFilterInput == "" {
		return SizeFilter{}, nil
	}

	if sizeFilterInput == ":" {
		return SizeFilter{}, fmt.Errorf("invalid size filter: ':' must be accompanied by a value")
	}

	// valid size format: "10MB", "5GB:", ":20KB", "1.5MB:2GB" (value + unit, optional range)
	pattern := `(?i)^(?:(\d+(?:\.\d+)?)(B|KB|MB|GB))?(?::(?:(\d+(?:\.\d+)?)(B|KB|MB|GB))?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(sizeFilterInput)
	isExactMatch := !strings.Contains(sizeFilterInput, ":")

	if matches == nil {
		return SizeFilter{}, fmt.Errorf("invalid size filter format: %q", sizeFilterInput)
	}

	startSize, err := parseSizeMatch(matches[1], matches[2], 0)
	if err != nil {
		return SizeFilter{}, err
	}

	endSize, err := parseSizeMatch(matches[3], matches[4], math.MaxInt64)
	if err != nil {
		return SizeFilter{}, err
	}

	return SizeFilter{
		startSize,
		endSize,
		isExactMatch,
	}, nil
}

func parseSizeMatch(value string, unit string, defaultSize int64) (int64, error) {
	if value == "" {
		return defaultSize, nil
	}

	return parseSizeInBytes(value, unit)
}

func parseSizeInBytes(valueInput string, unitInput string) (sizeInBytes int64, err error) {
	value, err := strconv.ParseFloat(valueInput, 64) // parseFloat for fractional input e.g. ">2.5KB"
	if err != nil {
		return 0, fmt.Errorf("invalid size value")
	}

	unit := strings.ToUpper(unitInput)

	switch unit {
	case "KB":
		sizeInBytes = int64(value * consts.KB)
	case "MB":
		sizeInBytes = int64(value * consts.MB)
	case "GB":
		sizeInBytes = int64(value * consts.GB)
	case "B":
		sizeInBytes = int64(value)
	default:
		return 0, fmt.Errorf("invalid size unit: %v", unit)
	}

	return sizeInBytes, nil
}

func parseColumns(columnsInput string, addColumnsInput string) ([]string, error) {
	if columnsInput != "" && addColumnsInput != "" {
		return nil, fmt.Errorf("cannot use --columns and --add-columns together. Use --columns to fully define the columns you want")
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

	validColumns := map[string]bool{
		consts.Date:       true,
		consts.Name:       true,
		consts.Reason:     true,
		consts.Size:       true,
		consts.Version:    true,
		consts.Depends:    true,
		consts.RequiredBy: true,
		consts.Provides:   true,
	}

	var columns []string

	for _, column := range strings.Split(columnInput, ",") {
		cleanColumn := strings.TrimSpace(column)

		if !validColumns[strings.TrimSpace(column)] {
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

	fmt.Println("\nColumn Options:")
	fmt.Println("  --columns <list>     Comma-separated list of columns to display (overrides defaults)")
	fmt.Println("  --add-columns <list> Comma-separated list of columns to add to defaults")

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
