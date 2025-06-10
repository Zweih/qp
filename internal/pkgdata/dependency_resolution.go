package pkgdata

import (
	"strings"
)

// TODO: we can do this concurrently. let's get on that.
func ResolveDependencyGraph(
	pkgs []*PkgInfo,
	complexEvaluator func(Relation, map[string][]string, map[string]*PkgInfo) []Relation,
) ([]*PkgInfo, error) {
	providesMap, installedMap := collectPkgData(pkgs)
	normalizeOptionalPkgs(pkgs, installedMap)

	var curriedEval func(Relation) []Relation
	if complexEvaluator != nil {
		curriedEval = func(rel Relation) []Relation {
			return complexEvaluator(rel, providesMap, installedMap)
		}
	}

	forwardShallow, reverseShallow, optReverseShallow := buildShallowGraph(
		pkgs,
		providesMap,
		curriedEval,
	)

	var visited map[string]int32

	for _, pkg := range pkgs {
		key := pkg.Key()

		visited = map[string]int32{key: 0}
		pkg.RequiredBy = collapseRelations(
			walkFullGraph(key, reverseShallow, visited, ""),
		)

		visited = map[string]int32{key: 0}
		pkg.Depends = collapseRelations(
			walkFullGraph(key, forwardShallow, visited, ""),
		)

		pkg.OptionalFor = collapseRelations(
			walkFullOptGraph(key, optReverseShallow[key], reverseShallow),
		)

		pkg.OptDepends = collapseRelations(
			walkFullOptGraph(key, pkg.OptDepends, forwardShallow),
		)
	}

	for _, pkg := range pkgs {
		pkg.Freeable = calculateFreeable(pkg, installedMap)
	}

	return pkgs, nil
}

func normalizeOptionalPkgs(pkgs []*PkgInfo, installedMap map[string]*PkgInfo) {
	for _, pkg := range pkgs {
		optForRels := filterToInstalled(pkg.OptionalFor, installedMap)
		for _, optForRel := range optForRels {
			if inst, ok := installedMap[optForRel.Key()]; ok {
				inst.OptDepends = append(
					inst.OptDepends, Relation{
						Name:    pkg.Name,
						Depth:   1,
						PkgType: pkg.PkgType,
					},
				)
			}
		}
	}
}

func filterToInstalled(rels []Relation, installedMap map[string]*PkgInfo) []Relation {
	var filteredRels []Relation

	for _, rel := range rels {
		if _, exists := installedMap[rel.Key()]; exists {
			filteredRels = append(filteredRels, rel)
		}
	}

	return filteredRels
}

func collectPkgData(pkgs []*PkgInfo) (map[string][]string, map[string]*PkgInfo) {
	// key: provided library/package, value: package that provides it (provider)
	providesMap := make(map[string][]string)
	installedMap := make(map[string]*PkgInfo, len(pkgs))

	for _, pkg := range pkgs {
		for _, provided := range pkg.Provides {
			key := provided.Key()
			providesMap[key] = append(providesMap[key], pkg.Name)
		}

		installedMap[pkg.Key()] = pkg
	}

	return providesMap, installedMap
}

func buildShallowGraph(
	pkgs []*PkgInfo,
	providesMap map[string][]string,
	complexEvaluator func(Relation) []Relation,
) (
	forwardShallow map[string][]Relation,
	reverseShallow map[string][]Relation,
	optReverseShallow map[string][]Relation,
) {
	forwardShallow = make(map[string][]Relation)
	reverseShallow = make(map[string][]Relation)
	optReverseShallow = make(map[string][]Relation)

	for _, pkg := range pkgs {
		buildShallowTrees(
			pkg,
			pkg.Depends,
			providesMap,
			forwardShallow,
			reverseShallow,
			complexEvaluator,
		)
		buildShallowTrees(pkg, pkg.OptDepends, providesMap, nil, optReverseShallow, nil)
	}

	return forwardShallow, reverseShallow, optReverseShallow
}

func buildShallowTrees(
	pkg *PkgInfo,
	deps []Relation,
	providesMap map[string][]string,
	forwardTree map[string][]Relation,
	reverseTree map[string][]Relation,
	complexEvaluator func(Relation) []Relation,
) {
	pkgKey := pkg.Key()
	for _, dep := range deps {
		depKey := dep.Key()
		if depKey == pkgKey {
			continue // prevent checking self-referencing packages
		}

		var depsToProcess []Relation

		if dep.IsComplex && complexEvaluator != nil {
			evaluatedDeps := complexEvaluator(dep)
			depsToProcess = evaluatedDeps
		} else {
			depsToProcess = []Relation{dep}
		}

		for _, processedDep := range depsToProcess {
			targets := resolveProvisions(processedDep.Key(), processedDep.Version, processedDep.Operator, providesMap)

			for _, target := range targets {
				if target.Name == pkg.Name {
					continue
				}

				targetCopy := target
				targetCopy.PkgType = processedDep.PkgType

				if forwardTree != nil {
					addToDependencyTree(pkgKey, forwardTree, targetCopy)
				}

				reverseKey := targetCopy.ProviderKey()
				if reverseKey == "" {
					reverseKey = targetCopy.Key()
				}

				reverseRelation := Relation{
					Name:     pkg.Name,
					Version:  processedDep.Version,
					Operator: processedDep.Operator,
					Depth:    1,
					Why:      processedDep.Why,
					PkgType:  pkg.PkgType,
				}

				addToDependencyTree(reverseKey, reverseTree, reverseRelation)
			}
		}
	}
}

func addToDependencyTree(
	from string,
	tree map[string][]Relation,
	relation Relation,
) {
	tree[from] = append(tree[from], relation)
}

func resolveProvisions(
	depKey string,
	version string,
	operator RelationOp,
	providesMap map[string][]string,
) []Relation {
	displayName := depKey
	if i := strings.IndexRune(depKey, ':'); i != -1 {
		displayName = depKey[i+1:]
	}

	if providerNames, exists := providesMap[depKey]; exists {
		provisions := make([]Relation, 0, len(providerNames))

		for _, providerName := range providerNames {
			provisions = append(provisions, Relation{
				Name:         displayName,
				Version:      version,
				Operator:     operator,
				ProviderName: providerName,
				Depth:        1,
			})
		}

		return provisions
	}

	return []Relation{{
		Name:     displayName,
		Version:  version,
		Operator: operator,
		Depth:    1,
	}}
}

// TODO: we can memoize this. we can also paralellize as well.
func walkFullGraph(
	key string,
	dependencyTree map[string][]Relation,
	visited map[string]int32,
	parentVirtualName string,
) []Relation {
	var results []Relation

	for _, relation := range dependencyTree[key] {
		relationKey := relation.Key()
		newDepth := visited[key] + 1
		prevDepth, seen := visited[relationKey]
		if seen && prevDepth <= newDepth {
			continue
		}

		visited[relationKey] = newDepth

		relationCopy := relation
		relationCopy.Depth = newDepth

		if relationCopy.ProviderName == "" && parentVirtualName != "" {
			relationCopy.ProviderName = parentVirtualName
		}

		results = append(results, relationCopy)

		subTree := walkFullGraph(relationKey, dependencyTree, visited, relationCopy.ProviderName)
		results = append(results, subTree...)
	}

	return results
}

func walkFullOptGraph(
	key string,
	optionalGraph []Relation,
	hardGraph map[string][]Relation,
) []Relation {
	var results []Relation
	visited := map[string]int32{key: 0}

	for _, optRelation := range optionalGraph {
		relKey := optRelation.Key()
		visited[relKey] = 1
		optRelationCopy := optRelation
		optRelationCopy.Depth = 1
		results = append(results, optRelationCopy)

		subTree := walkFullGraph(relKey, hardGraph, visited, "")
		results = append(results, subTree...)
	}

	return results
}

func collapseRelations(relations []Relation) []Relation {
	seen := map[string]Relation{}
	for _, relation := range relations {
		key := relation.Key()
		existingRelation, ok := seen[key]

		if !ok || relation.Depth < existingRelation.Depth {
			seen[key] = relation
		}
	}

	result := make([]Relation, 0, len(seen))
	for _, relation := range seen {
		result = append(result, relation)
	}

	return result
}
