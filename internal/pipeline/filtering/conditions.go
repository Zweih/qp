package filtering

import (
	"fmt"
	"qp/internal/config"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
)

type RangeSelector struct {
	Start   int64
	End     int64
	IsExact bool
}

type RangeMatcher map[bool]map[consts.MatchType]func(start int64, end int64) pkgdata.Filter

var DateMatchers = RangeMatcher{
	true: {
		consts.MatchFuzzy: func(start, _ int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.FuzzyDate(pkg, start)
			}
		},
		consts.MatchStrict: func(start, _ int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.StrictDate(pkg, start)
			}
		},
	},

	false: {
		consts.MatchFuzzy: func(start, end int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.FuzzyDateRange(pkg, start, end)
			}
		},
		consts.MatchStrict: func(start, end int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.StrictDateRange(pkg, start, end)
			}
		},
	},
}

var SizeMatchers = RangeMatcher{
	true: {
		consts.MatchFuzzy: func(start, _ int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.FuzzySize(pkg, start)
			}
		},
		consts.MatchStrict: func(start, _ int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.StrictSize(pkg, start)
			}
		},
	},
	false: {
		consts.MatchFuzzy: func(start, end int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.FuzzySizeRange(pkg, start, end)
			}
		},
		consts.MatchStrict: func(start, end int64) pkgdata.Filter {
			return func(pkg *pkgdata.PkgInfo) bool {
				return pkgdata.StrictSizeRange(pkg, start, end)
			}
		},
	},
}

var RangeMatchers = map[consts.FieldType]RangeMatcher{
	consts.FieldDate: DateMatchers,
	consts.FieldSize: SizeMatchers,
}

var StringMatchers = map[consts.MatchType]func(string, []string) bool{
	consts.MatchStrict: pkgdata.StrictStrings,
	consts.MatchFuzzy:  pkgdata.FuzzyStrings,
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
	query config.FieldQuery,
	selector RangeSelector,
) (*FilterCondition, error) {
	matchersByField, ok := RangeMatchers[query.Field]
	if !ok {
		return nil, fmt.Errorf("unsupported field type for range: %v", query.Field)
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
		return matcher(pkg) != mask
	}

	return &condition, nil
}
