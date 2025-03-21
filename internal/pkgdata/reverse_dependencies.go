package pkgdata

import (
	"regexp"
	"slices"
	"yaylog/internal/config"
	"yaylog/internal/consts"
)

// pulls package name out of `package-name>=2.0.1`
var packageNameRegex = regexp.MustCompile(`^([^<>=]+)`)

// TODO: we can do this concurrently. let's get on that.
func CalculateReverseDependencies(
	cfg config.Config,
	pkgPtrs []*PkgInfo,
	_ ProgressReporter, // TODO: Add progress reporting
) ([]*PkgInfo, error) {
	_, hasRequiredByFilter := cfg.FilterQueries[consts.FieldRequiredBy]

	if !slices.Contains(cfg.Fields, consts.FieldRequiredBy) && !hasRequiredByFilter {
		return pkgPtrs, nil
	}

	packagePointerMap := make(map[string]*PkgInfo)
	packageDependencyMap := make(map[string][]string)
	providesMap := make(map[string]string)
	// key: provided library/package, value: package that provides it (provider)

	for _, pkg := range pkgPtrs {
		packagePointerMap[pkg.Name] = pkg

		// populate providesMap
		for _, provided := range pkg.Provides {
			matches := packageNameRegex.FindStringSubmatch(provided)
			if len(matches) >= 2 {
				providesMap[matches[1]] = pkg.Name
			}
		}
	}

	for _, pkg := range pkgPtrs {
		for _, depPackage := range pkg.Depends {
			matches := packageNameRegex.FindStringSubmatch(depPackage)

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

	return pkgPtrs, nil
}
