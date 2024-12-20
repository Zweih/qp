package pkgdata

// FilterExplicit filters packages, returning only those with Reason == "explicit"
func FilterExplicit(pkgs []PackageInfo) []PackageInfo {
  var explicitPackages []PackageInfo

  for _, pkg := range pkgs {
    if pkg.Reason == "explicit" {
       explicitPackages = append(explicitPackages, pkg)
    }
  }

  return explicitPackages
}

// FilterDependencies filters packages, returning only those with Reason == "dependency"
func FilterDependencies(pkgs []PackageInfo) []PackageInfo {
  var dependencyPackages []PackageInfo

  for _, pkg := range pkgs {
    if pkg.Reason == "dependency" {
      dependencyPackages = append(dependencyPackages, pkg)
    }
  }

  return dependencyPackages
}
