package pkgdata

import (
	"qp/internal/pipeline/meta"
)

// TODO: we can do this concurrently. let's get on that.
func ResolveDependencyGraph(
	pkgs []*PkgInfo,
	_ meta.ProgressReporter, // TODO: Add progress reporting
) ([]*PkgInfo, error) {
	providesMap, installedMap := collectPkgData(pkgs)
	forwardShallow, reverseShallow, optReverseShallow := buildShallowGraph(pkgs, providesMap)

	var visited map[string]int32

	for _, pkg := range pkgs {
		name := pkg.Name

		visited = map[string]int32{name: 0}
		pkg.RequiredBy = collapseRelations(
			walkFullGraph(name, reverseShallow, visited, ""),
		)

		visited = map[string]int32{name: 0}
		pkg.Depends = collapseRelations(
			walkFullGraph(name, forwardShallow, visited, ""),
		)

		pkg.OptionalFor = filterToKeys(pkg.OptionalFor, installedMap)
		pkg.OptionalFor = collapseRelations(
			walkFullOptGraph(name, optReverseShallow[name], reverseShallow),
		)

		pkg.OptDepends = collapseRelations(
			walkFullOptGraph(name, pkg.OptDepends, forwardShallow),
		)
	}

	return pkgs, nil
}

func filterToKeys(rels []Relation, keyset map[string]struct{}) []Relation {
	var filtered []Relation
	for _, rel := range rels {
		if _, ok := keyset[rel.Name]; ok {
			filtered = append(filtered, rel)
		}
	}

	return filtered
}

func collectPkgData(pkgs []*PkgInfo) (map[string][]string, map[string]struct{}) {
	// key: provided library/package, value: package that provides it (provider)
	providesMap := make(map[string][]string)
	installedMap := make(map[string]struct{}, len(pkgs))

	for _, pkg := range pkgs {
		for _, provided := range pkg.Provides {
			providesMap[provided.Name] = append(providesMap[provided.Name], pkg.Name)
		}

		installedMap[pkg.Name] = struct{}{}
	}

	return providesMap, installedMap
}

func buildShallowGraph(
	pkgs []*PkgInfo,
	providesMap map[string][]string,
) (
	forwardShallow map[string][]Relation,
	reverseShallow map[string][]Relation,
	optReverseShallow map[string][]Relation,
) {
	forwardShallow = make(map[string][]Relation)
	reverseShallow = make(map[string][]Relation)
	optReverseShallow = make(map[string][]Relation)

	for _, pkg := range pkgs {
		buildShallowTrees(pkg, pkg.Depends, providesMap, forwardShallow, reverseShallow)
		buildShallowTrees(pkg, pkg.OptDepends, providesMap, nil, optReverseShallow)
	}

	return forwardShallow, reverseShallow, optReverseShallow
}

func buildShallowTrees(
	pkg *PkgInfo,
	deps []Relation,
	providesMap map[string][]string,
	forwardTree map[string][]Relation,
	reverseTree map[string][]Relation,
) {
	for _, depPackage := range deps {
		depName := depPackage.Name
		if depName == pkg.Name {
			continue // prevent checking self-referencing packages
		}

		targets := resolveProvisions(depName, depPackage.Version, depPackage.Operator, providesMap)

		for _, target := range targets {
			if target.Name == pkg.Name {
				continue
			}

			if forwardTree != nil {
				addToDependencyTree(pkg.Name, forwardTree, target)
			}

			reverseKey := target.ProviderName
			if reverseKey == "" {
				reverseKey = target.Name
			}

			reverseRelation := Relation{
				Name:     pkg.Name,
				Version:  depPackage.Version,
				Operator: depPackage.Operator,
				Depth:    1,
				Why:      depPackage.Why,
			}

			addToDependencyTree(reverseKey, reverseTree, reverseRelation)
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
func walkFullGraph(
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

		subTree := walkFullGraph(relation.Name, dependencyTree, visited, relationCopy.ProviderName)
		results = append(results, subTree...)
	}

	return results
}

func walkFullOptGraph(
	name string,
	optionalGraph []Relation,
	hardGraph map[string][]Relation,
) []Relation {
	var results []Relation
	visited := map[string]int32{name: 0}

	for _, optRelation := range optionalGraph {
		visited[optRelation.Name] = 1
		optRelation.Depth = 1
		results = append(results, optRelation)

		subTree := walkFullGraph(optRelation.Name, hardGraph, visited, "")
		results = append(results, subTree...)
	}

	return results
}

func collapseRelations(relations []Relation) []Relation {
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
