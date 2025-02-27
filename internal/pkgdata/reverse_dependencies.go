package pkgdata

import (
	"regexp"
	"slices"
	"yaylog/internal/config"
	"yaylog/internal/consts"
)

func CalculateReverseDependencies(
	cfg config.Config,
	packages []PackageInfo,
	reportProgress ProgressReporter,
) []PackageInfo {
	if !slices.Contains(cfg.ColumnNames, consts.RequiredBy) {
		return packages
	}

	packagePointerMap := make(map[string]*PackageInfo)
	packageDependencyMap := make(map[string][]string)
	providesMap := make(map[string]string) // key: provided library/package, value: package that providers it (provider)
	re := regexp.MustCompile(`^([^<>=]+)`) // pulls package name out of `package-name>=2.0.1`

	for i := range packages {
		pkg := &packages[i]
		packagePointerMap[pkg.Name] = pkg

		// populate providesMap
		for _, provided := range pkg.Provides {
			matches := re.FindStringSubmatch(provided)
			if len(matches) >= 2 {
				providesMap[matches[1]] = pkg.Name
			}
		}
	}

	for _, pkg := range packages {
		for _, depPackage := range pkg.Depends {
			matches := re.FindStringSubmatch(depPackage)

			if len(matches) >= 2 {
				depName := matches[1]

				if provider, exists := providesMap[depName]; exists {
					depName = provider
				}

				if depName == pkg.Name {
					continue // skip if a package names itself as a dependency
				}

				packageDependencyMap[depName] = append(packageDependencyMap[depName], pkg.Name)
			}
		}
	}

	for name, requiredBy := range packageDependencyMap {
		if pkg, exists := packagePointerMap[name]; exists {
			pkg.RequiredBy = requiredBy
		}
	}

	return packages
}
