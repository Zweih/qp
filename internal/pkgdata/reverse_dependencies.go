package pkgdata

import (
	"slices"
	"yaylog/internal/config"
	"yaylog/internal/consts"
	"yaylog/internal/pipeline/meta"
)

// TODO: we can do this concurrently. let's get on that.
func CalculateReverseDependencies(
	cfg config.Config,
	pkgPtrs []*PkgInfo,
	_ meta.ProgressReporter, // TODO: Add progress reporting
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if pipelineCtx.UsedCache {
		return pkgPtrs, nil
	}

	_, hasRequiredByFilter := cfg.FilterQueries[consts.FieldRequiredBy]
	hasRequiredByField := slices.Contains(cfg.Fields, consts.FieldRequiredBy)

	if !hasRequiredByField && !hasRequiredByFilter {
		return pkgPtrs, nil
	}

	packagePointerMap := make(map[string]*PkgInfo)
	packageDependencyMap := make(map[string][]Relation)
	providesMap := make(map[string]string)
	// key: provided library/package, value: package that provides it (provider)

	for _, pkg := range pkgPtrs {
		packagePointerMap[pkg.Name] = pkg

		// populate providesMap
		for _, provided := range pkg.Provides {
			providesMap[provided.Name] = pkg.Name
		}
	}

	for _, pkg := range pkgPtrs {
		for _, depPackage := range pkg.Depends {
			depName := depPackage.Name

			if providerName, exists := providesMap[depName]; exists {
				depName = providerName
			}

			if depName == pkg.Name {
				continue // skip if a package names itself as a dependency
			}

			packageDependencyMap[depName] = append(packageDependencyMap[depName], Relation{Name: pkg.Name})
		}
	}

	for name, requiredBy := range packageDependencyMap {
		if pkg, exists := packagePointerMap[name]; exists {
			pkg.RequiredBy = requiredBy
		}
	}

	return pkgPtrs, nil
}
