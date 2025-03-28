package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

func PrintHelp() {
	fmt.Println("Usage: yaylog [options]")

	fmt.Println("\nOptions:")
	pflag.PrintDefaults()

	fmt.Println("\nQuerying Options:")
	fmt.Println("  -w, --where <field>=<value> Apply queries to refine package listings. Can be used multiple times.")
	fmt.Println("                               Example: --where size=100MB:1GB --where name=firefox")

	fmt.Println("\n  Available queries:")
	fmt.Println("    date=<YYYY-MM-DD>               Show packages installed on a specific date")
	fmt.Println("    date=<YYYY-MM-DD>:              Show packages installed on or after the given date")
	fmt.Println("    date=:<YYYY-MM-DD>              Show packages installed up to the given date")
	fmt.Println("    date=<YYYY-MM-DD>:<YYYY-MM-DD>  Show packages installed in a date range")
	fmt.Println("    size=10MB:                      Show packages larger than 10MB")
	fmt.Println("    size=:500KB                     Show packages up to 500KB")
	fmt.Println("    size=1GB:5GB                    Show packages between 1GB and 5GB")
	fmt.Println("    name=firefox              Query packages by names (substring match)")
	fmt.Println("    reason=explicit           Show only explicitly installed packages")
	fmt.Println("    reason=dependencies       Show only packages installed as dependencies")
	fmt.Println("    required-by=vlc           Show packages required by specified packages")
	fmt.Println("    depends=glibc             Show packages that depend upon specified packages")
	fmt.Println("    provides=awk              Show packages that provide specified libraries, programs, or packages")
	fmt.Println("    conflicts=fuse            Show packages that conflict with the specified packages.")
	fmt.Println("    arch=x86_64               Show packages built for the specified architectures. \"any\" is a valid category of architecture.")

	fmt.Println("\nSorting Options:")
	fmt.Println("  -O, --order <type> Apply sorting to package output.")
	fmt.Println("  --order date                 Sort packages by installation date (default)")
	fmt.Println("  --order alphabetical         Sort packages alphabetically")
	fmt.Println("  --order size:desc            Sort packages by size in descending order")
	fmt.Println("  --order size:asc             Sort packages by size in ascending order")

	fmt.Println("\nOutput Options:")
	fmt.Println("  --json                      Output results in JSON format")
	fmt.Println("  --no-headers                Hide headers in table output (useful for scripts)")
	fmt.Println("  -s, --select <list>         Specify a comma-separated list of fields to display")
	fmt.Println("  -S, --select-add <list>     Add fields to the default view")
	fmt.Println("  -A, --select-all            Display all available fields")
	fmt.Println("  --full-timestamp            Show full timestamps (date + time) for package installations")

	fmt.Println("\nAvailable Fields:")
	fmt.Println("  date         Installation date of the package")
	fmt.Println("  name         Package name")
	fmt.Println("  reason       Installation reason (explicit/dependency)")
	fmt.Println("  size         Package size on disk")
	fmt.Println("  version      Installed package version")
	fmt.Println("  depends      List of dependencies (output can be long)")
	fmt.Println("  required-by  List of packages that depend on this package (output can be long)")
	fmt.Println("  provides     List of alternative package names or shared libraries provided (output can be long)")
	fmt.Println("  conflicts    List of packages that conflict, or cause problems, with the package")
	fmt.Println("  arch         Architecture the package was built for")
	fmt.Println("  license      Package software license")
	fmt.Println("  url          URL of the official site of the software being packaged")

	fmt.Println("\nExamples:")
	fmt.Println("  yaylog -l 10                      # Show the last 10 installed packages")
	fmt.Println("  yaylog -a -w reason=explicit      # Show all explicitly installed packages")
	fmt.Println("  yaylog -w reason=dependencies     # Show only dependencies")
	fmt.Println("  yaylog -w date=2024-12-25         # Show packages installed on a specific date")
	fmt.Println("  yaylog -w size=100MB:1GB          # Show packages between 100MB and 1GB")
	fmt.Println("  yaylog -w required-by=vlc         # Show packages required by VLC")
	fmt.Println("  yaylog --json                     # Output package data in JSON format")
	fmt.Println("  yaylog -w name=sqlite --json      # Output details for SQLite in JSON")
	fmt.Println("  yaylog --no-headers -s name,size  # Show package names and sizes without headers")

	fmt.Println("\nFor more details, see the manpage: man yaylog")
	fmt.Println("Or check the README on the GitHub repo.")
}
