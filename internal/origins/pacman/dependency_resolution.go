package pacman

import (
	"qp/internal/pipeline/meta"
	"qp/internal/pkgdata"
)

// TODO: we can do this concurrently. let's get on that.
func resolveDependencyGraph(
	pkgs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter, // TODO: Add progress reporting
) ([]*pkgdata.PkgInfo, error) {
	providesMap := buildProvidesMap(pkgs)
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

		pkg.OptionalFor = collapseRelations(
			walkFullOptGraph(name, optReverseShallow[name], reverseShallow),
		)

		pkg.OptDepends = collapseRelations(
			walkFullOptGraph(name, pkg.OptDepends, forwardShallow),
		)
	}

	return pkgs, nil
}

func buildProvidesMap(pkgs []*pkgdata.PkgInfo) map[string][]string {
	// key: provided library/package, value: package that provides it (provider)
	providesMap := make(map[string][]string)

	for _, pkg := range pkgs {
		for _, provided := range pkg.Provides {
			providesMap[provided.Name] = append(providesMap[provided.Name], pkg.Name)
		}
	}

	return providesMap
}

func buildShallowGraph(
	pkgs []*pkgdata.PkgInfo,
	providesMap map[string][]string,
) (
	forwardShallow map[string][]pkgdata.Relation,
	reverseShallow map[string][]pkgdata.Relation,
	optReverseShallow map[string][]pkgdata.Relation,
) {
	forwardShallow = make(map[string][]pkgdata.Relation)
	reverseShallow = make(map[string][]pkgdata.Relation)
	optReverseShallow = make(map[string][]pkgdata.Relation)

	for _, pkg := range pkgs {
		buildShallowTrees(pkg, pkg.Depends, providesMap, forwardShallow, reverseShallow)
		buildShallowTrees(pkg, pkg.OptDepends, providesMap, nil, optReverseShallow)
	}

	return forwardShallow, reverseShallow, optReverseShallow
}

func buildShallowTrees(
	pkg *pkgdata.PkgInfo,
	deps []pkgdata.Relation,
	providesMap map[string][]string,
	forwardTree map[string][]pkgdata.Relation,
	reverseTree map[string][]pkgdata.Relation,
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

			reverseRelation := pkgdata.Relation{
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
	tree map[string][]pkgdata.Relation,
	relation pkgdata.Relation,
) {
	tree[from] = append(tree[from], relation)
}

func resolveProvisions(
	depName string,
	version string,
	operator pkgdata.RelationOp,
	providesMap map[string][]string,
) []pkgdata.Relation {
	if providerNames, exists := providesMap[depName]; exists {
		provisions := make([]pkgdata.Relation, 0, len(providerNames))

		for _, providerName := range providerNames {
			provisions = append(provisions, pkgdata.Relation{
				Name:         depName,
				Version:      version,
				Operator:     operator,
				ProviderName: providerName,
				Depth:        1,
			})
		}

		return provisions
	}

	return []pkgdata.Relation{{Name: depName, Version: version, Operator: operator}}
}

// TODO: we can memoize this. we can also paralellize as well.
func walkFullGraph(
	name string,
	dependencyTree map[string][]pkgdata.Relation,
	visited map[string]int32,
	parentVirtualName string,
) []pkgdata.Relation {
	var results []pkgdata.Relation

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
	optionalGraph []pkgdata.Relation,
	hardGraph map[string][]pkgdata.Relation,
) []pkgdata.Relation {
	var results []pkgdata.Relation
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

func collapseRelations(relations []pkgdata.Relation) []pkgdata.Relation {
	seen := map[string]pkgdata.Relation{}
	for _, relation := range relations {
		existingRelation, ok := seen[relation.Name]

		if !ok || relation.Depth < existingRelation.Depth {
			seen[relation.Name] = relation
		}
	}

	result := make([]pkgdata.Relation, 0, len(seen))
	for _, relation := range seen {
		result = append(result, relation)
	}

	return result
}
