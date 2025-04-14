package filtering

import (
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
)

type RangeSelector struct {
	Start   int64
	End     int64
	IsExact bool
}

type ExactFilter func(pkg *PkgInfo, target int64) bool

type RangeFilter func(pkg *PkgInfo, start int64, end int64) bool

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
) (*FilterCondition, error) {
	condition := newCondition(field)
	matcher := getStringMatcher(match)

	for i, target := range targets {
		targets[i] = strings.ToLower(target)
	}

	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		return matcher(pkg.GetString(field), targets)
	}

	return &condition, nil
}

func newRelationCondition(
	field consts.FieldType,
	targets []string,
	depth int32,
	match consts.MatchType,
) (*FilterCondition, error) {
	condition := newCondition(field)
	matcher := getStringMatcher(match)

	for i, target := range targets {
		targets[i] = strings.ToLower(target)
	}

	condition.Filter = func(pkg *pkgdata.PkgInfo) bool {
		relationsAtDepth := pkgdata.GetRelationsByDepth(pkg.GetRelations(field), depth)
		for _, rel := range relationsAtDepth {
			if matcher(rel.Name, targets) {
				return true
			}
		}

		return false
	}

	return &condition, nil
}

func newRangeCondition(
	rangeSelector RangeSelector,
	fieldType consts.FieldType,
	exactFunc ExactFilter,
	rangeFunc RangeFilter,
) *FilterCondition {
	condition := newCondition(fieldType)

	if rangeSelector.IsExact {
		condition.Filter = func(pkg *PkgInfo) bool {
			return exactFunc(pkg, rangeSelector.Start)
		}

		return &condition
	}

	condition.Filter = func(pkg *PkgInfo) bool {
		return rangeFunc(pkg, rangeSelector.Start, rangeSelector.End)
	}

	return &condition
}

func newDateCondition(dateFilter RangeSelector) *FilterCondition {
	return newRangeCondition(
		dateFilter,
		consts.FieldDate,
		pkgdata.FilterByDate,
		pkgdata.FilterByDateRange,
	)
}

func newSizeCondition(sizeFilter RangeSelector) *FilterCondition {
	return newRangeCondition(
		sizeFilter,
		consts.FieldSize,
		pkgdata.FilterBySize, // TODO: maybe these two should be in maps that have the Field as a key
		pkgdata.FilterBySizeRange,
	)
}

func newReasonCondition(reason string) *FilterCondition {
	condition := newCondition(consts.FieldReason)
	condition.Filter = func(pkg *PkgInfo) bool {
		return pkgdata.FilterByReason(pkg.Reason, reason)
	}

	return &condition
}

func getStringMatcher(match consts.MatchType) func(string, []string) bool {
	switch match {
	case consts.MatchExact:
		return pkgdata.ExactStrings
	default:
		return pkgdata.FuzzyStrings
	}
}
