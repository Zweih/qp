package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

func PrintHelp() {
	fmt.Println("Usage: qp [options]")

	fmt.Println("\nOptions:")
	pflag.PrintDefaults()

	fmt.Println("\nQuerying Options:")
	fmt.Println("  -w, --where <field>=<value> Apply queries to refine package listings. Can be used multiple times.")
	fmt.Println("                               Example: --where size=100MB:1GB --where name=firefox")

	fmt.Println("\n  Available queries:")
	fmt.Println("    date=<YYYY-MM-DD>               Query packages installed on a specific date")
	fmt.Println("    date=<YYYY-MM-DD>:              Query packages installed on or after the given date")
	fmt.Println("    date=:<YYYY-MM-DD>              Query packages installed up to the given date")
	fmt.Println("    date=<YYYY-MM-DD>:<YYYY-MM-DD>  Query packages installed in a date range")
	fmt.Println("    size=10MB:                      Query packages larger than 10MB")
	fmt.Println("    size=:500KB                     Query packages up to 500KB")
	fmt.Println("    size=1GB:5GB                    Query packages between 1GB and 5GB")
	fmt.Println("    name=firefox              Query packages by names (substring match)")
	fmt.Println("    reason=explicit           Query only explicitly installed packages")
	fmt.Println("    reason=dependencies       Query only packages installed as dependencies")
	fmt.Println("    required-by=vlc           Query packages required by specified packages")
	fmt.Println("    depends=glibc             Query packages that depend upon specified packages")
	fmt.Println("    provides=awk              Query packages that provide specified libraries, programs, or packages")
	fmt.Println("    conflicts=fuse            Query packages that conflict with the specified packages")
	fmt.Println("    arch=x86_64               Show packages built for the specified architectures. \"any\" is a valid category of architecture.")
	fmt.Println("    description=x86_64        Query packages by description (substring match)")

	fmt.Println("\nSorting Options:")
	fmt.Println("  -O, --order <type>:<direction> Apply sorting to package output")
	fmt.Println("                                 Default sort is date:asc")
	fmt.Println("  --order date                   Sort packages by installation date")
	fmt.Println("  --order name                   Sort packages alphabetically by package name")
	fmt.Println("  --order size                   Sort packages by size in descending order")
	fmt.Println("  --order license                Sort packages alphabetically by package license")

	fmt.Println("\nOutput Options:")
	fmt.Println("  --json                      Output results in JSON format")
	fmt.Println("  --no-headers                Hide headers in table output (useful for scripts)")
	fmt.Println("  -s, --select <list>         Specify a comma-separated list of fields to display")
	fmt.Println("  -S, --select-add <list>     Add fields to the default view")
	fmt.Println("  -A, --select-all            Display all available fields")
	fmt.Println("  --full-timestamp            Show full timestamps (date + time) for package installations")

	fmt.Println("\nAvailable Fields:")
	fmt.Println("  date         Installation date of the package")
	fmt.Println("  build-date   Date the package was built")
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
	fmt.Println("  description  Package description")
	fmt.Println("  pkgbase      Name of the base package used to group split packages; for non-split packages, it is the same as the package name.")
	fmt.Println("  pkgtype      Type of the package (standard, split, debug, source, unknown)")
	fmt.Println("               Note: Older packages may show \"unknown\" pkgtype if built before pacman introduced XDATA.")
	fmt.Println("  groups       Package groups or categories (e.g., base, gnome, xfce4)")

	fmt.Println("\nExamples:")
	fmt.Println("  qp -l 10                      # Show the last 10 installed packages")
	fmt.Println("  qp -a -w reason=explicit      # Show all explicitly installed packages")
	fmt.Println("  qp -w reason=dependencies     # Show only dependencies")
	fmt.Println("  qp -w date=2024-12-25         # Show packages installed on a specific date")
	fmt.Println("  qp -w size=100MB:1GB          # Show packages between 100MB and 1GB")
	fmt.Println("  qp -w required-by=vlc         # Show packages required by VLC")
	fmt.Println("  qp --json                     # Output package data in JSON format")
	fmt.Println("  qp -w name=sqlite --json      # Output details for SQLite in JSON")
	fmt.Println("  qp --no-headers -s name,size  # Show package names and sizes without headers")

	fmt.Println("\nFor more details, see the manpage: man qp")
	fmt.Println("Or check the README on the GitHub repo.")
}
