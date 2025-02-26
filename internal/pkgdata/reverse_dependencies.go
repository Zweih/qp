package pkgdata

import (
	"fmt"
	"regexp"
	"yaylog/internal/config"
)

func CalculateReverseDependencies(
	cfg config.Config,
	packages []PackageInfo,
	reportProgress ProgressReporter,
) []PackageInfo {
	packagePointerMap := make(map[string]*PackageInfo)
	packageDependencyMap := make(map[string][]string)
	re := regexp.MustCompile(`^([^<>=]+)`)

	for i := range packages {
		packagePointerMap[packages[i].Name] = &packages[i]

		if len(packages[i].Depends) > 0 {
			for _, depPackage := range packages[i].Depends {
				matches := re.FindStringSubmatch(depPackage)
				depPackageName := matches[1]

				if len(matches) >= 2 {
					packageDependencyMap[depPackageName] = append(packageDependencyMap[depPackageName], packages[i].Name)
				}
			}
		}
	}

	for name, requiredBy := range packageDependencyMap {
		if pkg, exists := packagePointerMap[name]; exists {
			pkg.RequiredBy = requiredBy
		} else {
			fmt.Printf("[COULDN'T FIND] Name: %s RequiredBy: %v", name, requiredBy)
		}
		// packagePointerMap[name].RequiredBy = requiredBy
	}

	return packages
}
