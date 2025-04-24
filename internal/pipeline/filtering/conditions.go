package filtering

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"qp/internal/query"
	"strings"
)

type RangeSelector struct {
	Start   int64
	End     int64
	IsExact bool
}

func newCondition(fieldType consts.FieldType) FilterCondition {
	return FilterCondition{
		PhaseName: "Filtering by " + consts.FieldNameLookup[fieldType],
		FieldType: fieldType,
	}
}

func newStringCondition(
	field consts.FieldType,
	targets []string,
	match consts.MatchType,
	mask bool,
) (*FilterCondition, error) {
	condition := newCondition(field)
	matcher := StringMatchers[match]

	for i, target := range targets {
		targets[i] = strings.ToLower(target)
	}

	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		return matcher(pkg.GetString(field), targets) != mask
	}

	return &condition, nil
}

func newStringExistsCondition(field consts.FieldType, mask bool) (*FilterCondition, error) {
	condition := newCondition(field)
	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		return pkgdata.StringExists(pkg.GetString(field)) != mask
	}

	return &condition, nil
}

func newRelationCondition(
	field consts.FieldType,
	targets []string,
	depth int32,
	match consts.MatchType,
	mask bool,
) (*FilterCondition, error) {
	condition := newCondition(field)
	matcher := StringMatchers[match]

	for i, target := range targets {
		targets[i] = strings.ToLower(target)
	}

	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		relationsAtDepth := pkgdata.GetRelationsByDepth(pkg.GetRelations(field), depth)
		for _, rel := range relationsAtDepth {
			if matcher(rel.Name, targets) {
				return !mask
			}
		}

		return mask
	}

	return &condition, nil
}

func newRelationExistsCondition(
	field consts.FieldType,
	depth int32,
	mask bool,
) (*FilterCondition, error) {
	condition := newCondition(field)

	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		relationsAtDepth := pkgdata.GetRelationsByDepth(pkg.GetRelations(field), depth)
		return pkgdata.RelationExists(relationsAtDepth) != mask
	}

	return &condition, nil
}

func newRangeCondition(
	query query.FieldQuery,
	selector RangeSelector,
) (*FilterCondition, error) {
	matchersByField, ok := RangeMatchers[query.Field]
	if !ok {
		return nil, fmt.Errorf(
			"unsupported field type for range: %v", consts.FieldNameLookup[query.Field],
		)
	}

	matchersByExact, ok := matchersByField[selector.IsExact]
	if !ok {
		return nil, fmt.Errorf("internal error: missing exactness entry")
	}

	builder, ok := matchersByExact[query.Match]
	if !ok {
		return nil, fmt.Errorf("unsupported match type: %v", query.Match)
	}

	matcher := builder(selector.Start, selector.End)
	condition := newCondition(query.Field)
	mask := query.Negate

	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		return matcher(pkg.GetInt(query.Field)) != mask
	}

	return &condition, nil
}
