package pkgdata

// TODO: combine functions to allow for mixed arguments

// filters packages only those listed as FilterExplicit
func FilterExplicit(pkgs []PackageInfo) []PackageInfo {
	var explicitPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Reason == "explicit" {
			explicitPackages = append(explicitPackages, pkg)
		}
	}

	return explicitPackages
}

// filters packages only listed as dependencies
func FilterDependencies(pkgs []PackageInfo) []PackageInfo {
	var dependencyPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Reason == "dependency" {
			dependencyPackages = append(dependencyPackages, pkg)
		}
	}

	return dependencyPackages
}
