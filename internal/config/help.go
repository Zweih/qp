package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

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
	fmt.Println("  --columns <list>     Comma-separated list of columns to display (cannot use with --all-columns or --add-columns)")
	fmt.Println("  --add-columns <list> Comma-separated list of columns to add to defaults or --all-columns")
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
