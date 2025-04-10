package pkgdata

import (
	"qp/internal/pipeline/meta"
)

// TODO: we can do this concurrently. let's get on that.
func ResolveDependencyTree(
	pkgPtrs []*PkgInfo,
	_ meta.ProgressReporter, // TODO: Add progress reporting
) ([]*PkgInfo, error) {
	packagePointerMap := make(map[string]*PkgInfo)
	reverseDependencyTree := make(map[string][]Relation)
	forwardDependencyTree := make(map[string][]Relation)
	providesMap := make(map[string][]string)
	// key: provided library/package, value: package that provides it (provider)

	for _, pkg := range pkgPtrs {
		packagePointerMap[pkg.Name] = pkg

		// populate providesMap
		for _, provided := range pkg.Provides {
			providesMap[provided.Name] = append(providesMap[provided.Name], pkg.Name)
		}
	}

	for _, pkg := range pkgPtrs {
		for _, depPackage := range pkg.Depends {
			depName := depPackage.Name
			if depName == pkg.Name {
				continue // prevent checking self-referencing packages
			}

			targets := resolveProvisions(depName, depPackage.Version, depPackage.Operator, providesMap)

			for _, target := range targets {
				if target.Name == pkg.Name {
					continue
				}

				addToDependencyTree(pkg.Name, forwardDependencyTree, target)

				reverseKey := target.ProviderName
				if reverseKey == "" {
					reverseKey = target.Name
				}

				reverseRelation := Relation{
					Name:     pkg.Name,
					Version:  depPackage.Version,
					Operator: depPackage.Operator,
					Depth:    1,
				}

				addToDependencyTree(reverseKey, reverseDependencyTree, reverseRelation)
			}
		}
	}

	for _, pkg := range pkgPtrs {
		name := pkg.Name
		visitedReverse := map[string]int32{name: 0}
		visitedForward := map[string]int32{name: 0}

		fullReverseDepGraph := walkDependencyGraph(name, reverseDependencyTree, visitedReverse, "")
		fullForwardDepGraph := walkDependencyGraph(name, forwardDependencyTree, visitedForward, "")

		pkg.RequiredBy = dedupToLowestDepth(fullReverseDepGraph)
		pkg.Depends = dedupToLowestDepth(fullForwardDepGraph)
	}

	return pkgPtrs, nil
}

func addToDependencyTree(
	from string,
	tree map[string][]Relation,
	relation Relation,
) {
	tree[from] = append(tree[from], relation)
}

func resolveProvisions(
	depName string,
	version string,
	operator RelationOp,
	providesMap map[string][]string,
) []Relation {
	if providerNames, exists := providesMap[depName]; exists {
		provisions := make([]Relation, 0, len(providerNames))

		for _, providerName := range providerNames {
			provisions = append(provisions, Relation{
				Name:         depName,
				Version:      version,
				Operator:     operator,
				ProviderName: providerName,
				Depth:        1,
			})
		}

		return provisions
	}

	return []Relation{{Name: depName, Version: version, Operator: operator}}
}

// TODO: we can memoize this. we can also paralellize as well.
func walkDependencyGraph(
	name string,
	dependencyTree map[string][]Relation,
	visited map[string]int32,
	parentVirtualName string,
) []Relation {
	var results []Relation

	for _, relation := range dependencyTree[name] {
		newDepth := visited[name] + 1
		prevDepth, seen := visited[relation.Name]
		if seen && prevDepth <= newDepth {
			continue
		}

		visited[relation.Name] = newDepth

		relationCopy := relation
		relationCopy.Depth = newDepth

		if relationCopy.ProviderName == "" && parentVirtualName != "" {
			relationCopy.ProviderName = parentVirtualName
		}

		results = append(results, relationCopy)

		subTree := walkDependencyGraph(relation.Name, dependencyTree, visited, relationCopy.ProviderName)
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
