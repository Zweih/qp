package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

func PrintHelp() {
	fmt.Println("Usage: yaylog [options]")

	fmt.Println("\nOptions:")
	pflag.PrintDefaults()

	fmt.Println("\nFiltering Options:")
	fmt.Println("  -f, --filter <field>=<value> Apply filters to refine package listings. Can be used multiple times.")
	fmt.Println("                               Example: --filter size=100MB:1GB --filter name=firefox")

	fmt.Println("\n  Available filters:")
	fmt.Println("    date=<YYYY-MM-DD>               Show packages installed on a specific date")
	fmt.Println("    date=<YYYY-MM-DD>:              Show packages installed on or after the given date")
	fmt.Println("    date=:<YYYY-MM-DD>              Show packages installed up to the given date")
	fmt.Println("    date=<YYYY-MM-DD>:<YYYY-MM-DD>  Show packages installed in a date range")
	fmt.Println("    size=10MB:                      Show packages larger than 10MB")
	fmt.Println("    size=:500KB                     Show packages up to 500KB")
	fmt.Println("    size=1GB:5GB                    Show packages between 1GB and 5GB")
	fmt.Println("    name=firefox              Filter packages by name (substring match)")
	fmt.Println("    reason=explicit           Show only explicitly installed packages")
	fmt.Println("    reason=dependencies       Show only packages installed as dependencies")
	fmt.Println("    required-by=vlc           Show packages required by the specified package")
	fmt.Println("    depends=glibc             Show packages that depend upon a specific package")
	fmt.Println("    provides=awk              Show packages that provide a specific library, program, or package")

	fmt.Println("\nSorting Options:")
	fmt.Println("  --sort date                 Sort packages by installation date (default)")
	fmt.Println("  --sort alphabetical         Sort packages alphabetically")
	fmt.Println("  --sort size:desc            Sort packages by size in descending order")
	fmt.Println("  --sort size:asc             Sort packages by size in ascending order")

	fmt.Println("\nOutput Options:")
	fmt.Println("  --json                      Output results in JSON format")
	fmt.Println("  --no-headers                Hide headers in table output (useful for scripts)")
	fmt.Println("  --columns <list>            Specify a comma-separated list of columns to display")
	fmt.Println("  --add-columns <list>        Add columns to the default view")
	fmt.Println("  --all-columns               Display all available columns")
	fmt.Println("  --full-timestamp            Show full timestamps (date + time) for package installations")

	fmt.Println("\nAvailable Columns:")
	fmt.Println("  date         Installation date of the package")
	fmt.Println("  name         Package name")
	fmt.Println("  reason       Installation reason (explicit/dependency)")
	fmt.Println("  size         Package size on disk")
	fmt.Println("  version      Installed package version")
	fmt.Println("  depends      List of dependencies (output can be long)")
	fmt.Println("  required-by  List of packages that depend on this package (output can be long)")
	fmt.Println("  provides     List of alternative package names or shared libraries provided")

	fmt.Println("\nDeprecated Legacy Options (Use --filter Instead):")
	fmt.Println("  -e, --explicit              (Deprecated) Show only explicitly installed packages")
	fmt.Println("                                Equivalent to: --filter reason=explicit")
	fmt.Println("  -d, --dependencies          (Deprecated) Show only packages installed as dependencies")
	fmt.Println("                                Equivalent to: --filter reason=dependencies")
	fmt.Println("  --date <filter>             (Deprecated) Filter by installation date")
	fmt.Println("                                Equivalent to: --filter date=YYYY-MM-DD")
	fmt.Println("  --size <filter>             (Deprecated) Filter by size")
	fmt.Println("                                Equivalent to: --filter size=100MB:1GB")
	fmt.Println("  --name <search-term>        (Deprecated) Filter by name")
	fmt.Println("                                Equivalent to: --filter name=vim")
	fmt.Println("  --required-by <package>     (Deprecated) Show packages required by another package")
	fmt.Println("                                Equivalent to: --filter required-by=firefox")

	fmt.Println("\nExamples:")
	fmt.Println("  yaylog -n 10                        # Show the last 10 installed packages")
	fmt.Println("  yaylog -f reason=explicit           # Show all explicitly installed packages")
	fmt.Println("  yaylog -f reason=dependencies       # Show only dependencies")
	fmt.Println("  yaylog -f date=2024-12-25           # Show packages installed on a specific date")
	fmt.Println("  yaylog -f size=100MB:1GB            # Show packages between 100MB and 1GB")
	fmt.Println("  yaylog -f required-by=vlc           # Show packages required by VLC")
	fmt.Println("  yaylog --json                       # Output package data in JSON format")
	fmt.Println("  yaylog -f name=sqlite --json        # Output details for SQLite in JSON")
	fmt.Println("  yaylog --no-headers --columns name,size  # Show package names and sizes without headers")

	fmt.Println("\nFor more details, see the manpage: man yaylog")
	fmt.Println("Or check the README on the GitHub repo.")
}
