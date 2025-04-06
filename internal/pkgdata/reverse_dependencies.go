package pkgdata

import (
	"qp/internal/pipeline/meta"
)

// TODO: we can do this concurrently. let's get on that.
func CalculateReverseDependencies(
	pkgPtrs []*PkgInfo,
	_ meta.ProgressReporter, // TODO: Add progress reporting
) ([]*PkgInfo, error) {
	packagePointerMap := make(map[string]*PkgInfo)
	reverseDependencyTree := make(map[string][]Relation)
	providesMap := make(map[string]string)
	// key: provided library/package, value: package that provides it (provider)

	for _, pkg := range pkgPtrs {
		packagePointerMap[pkg.Name] = pkg

		// populate providesMap
		for _, provided := range pkg.Provides {
			// TODO: this assumes only one provider exists to satisfy a dependency, but that's not always true. this should be an array of names, not just one name
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

			reverseDependencyTree[depName] = append(
				reverseDependencyTree[depName],
				Relation{
					Name:     pkg.Name,
					Version:  depPackage.Version,
					Operator: depPackage.Operator,
					Depth:    1,
				})
		}
	}

	for name := range reverseDependencyTree {
		if pkg, exists := packagePointerMap[name]; exists {
			visited := map[string]int32{name: 0}
			fullDependencyTree := walkReverseDeps(name, reverseDependencyTree, visited)
			pkg.RequiredBy = dedupToLowestDepth(fullDependencyTree)
		}
	}

	return pkgPtrs, nil
}

// TODO: we can memoize this. we can also paralellize as well.
func walkReverseDeps(
	name string,
	reverseMap map[string][]Relation,
	visited map[string]int32,
) []Relation {
	var results []Relation

	for _, relation := range reverseMap[name] {
		newDepth := visited[name] + 1
		prevDepth, seen := visited[relation.Name]
		if seen && prevDepth <= newDepth {
			continue
		}

		visited[relation.Name] = newDepth

		relationCopy := relation
		relationCopy.Depth = newDepth
		results = append(results, relationCopy)

		subTree := walkReverseDeps(relation.Name, reverseMap, visited)
		results = append(results, subTree...)
	}

	return results
}

func dedupToLowestDepth(relations []Relation) []Relation {
	seen := map[string]Relation{}
	for _, relation := range relations {
		existingRelation, ok := seen[relation.Name]
		if !ok || relation.Depth < existingRelation.Depth {
			seen[relation.Name] = relation
		}
	}

	result := make([]Relation, 0, len(seen))
	for _, relation := range seen {
		result = append(result, relation)
	}

	return result
}
