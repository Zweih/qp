package filtering

import (
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type RangeFilter func(value int64) bool

type RangeMatcher map[bool]map[consts.MatchType]func(start int64, end int64) RangeFilter

var DateMatchers = RangeMatcher{
	true: {
		consts.MatchFuzzy: func(start, _ int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.FuzzyDate(value, start)
			}
		},
		consts.MatchStrict: func(start, _ int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.StrictDate(value, start)
			}
		},
	},

	false: {
		consts.MatchFuzzy: func(start, end int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.FuzzyDateRange(value, start, end)
			}
		},
		consts.MatchStrict: func(start, end int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.StrictDateRange(value, start, end)
			}
		},
	},
}

var SizeMatchers = RangeMatcher{
	true: {
		consts.MatchFuzzy: func(start, _ int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.FuzzySize(value, start)
			}
		},
		consts.MatchStrict: func(start, _ int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.StrictSize(value, start)
			}
		},
	},
	false: {
		consts.MatchFuzzy: func(start, end int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.FuzzySizeRange(value, start, end)
			}
		},
		consts.MatchStrict: func(start, end int64) RangeFilter {
			return func(value int64) bool {
				return pkgdata.StrictSizeRange(value, start, end)
			}
		},
	},
}

var RangeMatchers = map[consts.FieldType]RangeMatcher{
	consts.FieldDate:      DateMatchers,
	consts.FieldBuildDate: DateMatchers,
	consts.FieldSize:      SizeMatchers,
}

var StringMatchers = map[consts.MatchType]func(string, []string) bool{
	consts.MatchStrict: pkgdata.StrictStrings,
	consts.MatchFuzzy:  pkgdata.FuzzyStrings,
}
