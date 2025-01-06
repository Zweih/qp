package pkgdata

import "time"

// TODO: combine functions to allow for mixed arguments

func FilterExplicit(pkgs []PackageInfo) []PackageInfo {
	var explicitPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Reason == "explicit" {
			explicitPackages = append(explicitPackages, pkg)
		}
	}

	return explicitPackages
}

func FilterDependencies(pkgs []PackageInfo) []PackageInfo {
	var dependencyPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Reason == "dependency" {
			dependencyPackages = append(dependencyPackages, pkg)
		}
	}

	return dependencyPackages
}

// filters packages installed on specific date
func FilterByDate(pkgs []PackageInfo, date time.Time) []PackageInfo {
	var filteredPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Timestamp.Year() == date.Year() && pkg.Timestamp.YearDay() == date.YearDay() {
			filteredPackages = append(filteredPackages, pkg)
		}
	}

	return filteredPackages
}
