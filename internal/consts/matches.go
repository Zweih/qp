package consts

type (
	MatchType int32
)

const (
	MatchFuzzy MatchType = iota
	MatchStrict
)
