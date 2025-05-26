package rpm

import (
	"fmt"
	"qp/internal/pkgdata"
	"strconv"
	"strings"
)

func evaluateComplexDependency(
	relation pkgdata.Relation,
	providesMap map[string][]string,
	installedMap map[string]*pkgdata.PkgInfo,
) []pkgdata.Relation {
	if !relation.IsComplex {
		return []pkgdata.Relation{relation}
	}

	expr := strings.Trim(relation.Name, "()")

	parser := NewParser(expr)
	ast, err := parser.ParseExpression()
	if err != nil {
		return []pkgdata.Relation{{
			Name:      relation.Name,
			Depth:     relation.Depth,
			IsComplex: false,
		}}
	}

	results := ast.Evaluate(installedMap)

	for i := range results {
		results[i].Depth = relation.Depth
		results[i].PkgType = relation.PkgType
	}

	return results
}

func parseRelationList(rawRels []string) []pkgdata.Relation {
	if len(rawRels) < 1 {
		return nil
	}

	rels := make([]pkgdata.Relation, 0, len(rawRels))
	for _, rawRel := range rawRels {
		if rel, err := parseRelation(rawRel); err == nil {
			rels = append(rels, rel)
		}
	}

	return rels
}

func parseRelation(rawRel string) (pkgdata.Relation, error) {
	trimmed := strings.TrimSpace(rawRel)
	if trimmed == "" {
		return pkgdata.Relation{}, fmt.Errorf("relation empty")
	}

	isRichDep := trimmed[0] == '(' && trimmed[len(trimmed)-1] == ')'
	if isRichDep {
		return pkgdata.Relation{
			Name:      trimmed,
			Depth:     1,
			IsComplex: true,
		}, nil
	}

	rel := pkgdata.Relation{
		Name:  trimmed,
		Depth: 1,
	}

	if strings.Contains(rel.Name, " ") {
		parts := strings.Fields(rel.Name)
		if len(parts) == 3 {
			rel.Name = parts[0]
			rel.Operator = pkgdata.StringToOperator(parts[1])
			rel.Version = parts[2]
		}
	}

	return rel, nil
}

func parseVersion(rpmEpoch *int, rpmVersion string) string {
	if rpmEpoch != nil && *rpmEpoch != 0 {
		return strconv.Itoa(*rpmEpoch) + ":" + rpmVersion
	}

	return rpmVersion
}
