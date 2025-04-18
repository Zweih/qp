package pkgdata

import (
	"slices"
	"strings"
	"time"
)

const fuzzySizeTolerancePercent = 0.3

func FilterByReason(installReason string, targetReason string) bool {
	return installReason == targetReason
}

func GetRelationsByDepth(relations []Relation, targetDepth int32) []Relation {
	filteredRelations := []Relation{}

	for _, relation := range relations {
		if relation.Depth == targetDepth {
			filteredRelations = append(filteredRelations, relation)
		}
	}

	return filteredRelations
}

func FuzzyDate(value int64, date int64) bool {
	pkgDate := time.Unix(value, 0)
	targetDate := time.Unix(date, 0) // TODO: we can pull this out to the top level
	return pkgDate.Year() == targetDate.Year() && pkgDate.YearDay() == targetDate.YearDay()
}

func FuzzyDateRange(value int64, start int64, end int64) bool {
	pkgDate := time.Unix(value, 0).Truncate(24 * time.Hour)
	startDate := time.Unix(start, 0).Truncate(24 * time.Hour)
	endDate := time.Unix(end, 0).Truncate(24 * time.Hour)

	return (pkgDate.Equal(startDate) || pkgDate.After(startDate)) &&
		(pkgDate.Equal(endDate) || pkgDate.Before(endDate))
}

func StrictDate(value int64, targetDate int64) bool {
	return value == targetDate
}

func StrictDateRange(value int64, start int64, end int64) bool {
	return !(value < start || value > end)
}

func FuzzySizeTolerance(targetSize int64) int64 {
	return int64(float64(targetSize) * fuzzySizeTolerancePercent / 100.0)
}

func FuzzySize(value int64, targetSize int64) bool {
	tolerance := FuzzySizeTolerance(targetSize)
	result := value - targetSize
	return max(result, -result) <= tolerance
}

func FuzzySizeRange(value, start int64, end int64) bool {
	toleranceStart := FuzzySizeTolerance(start)
	toleranceEnd := FuzzySizeTolerance(end)

	return value >= (start-toleranceStart) && value <= (end+toleranceEnd)
}

func StrictSize(value int64, targetSize int64) bool {
	return value == targetSize
}

func StrictSizeRange(value int64, startSize int64, endSize int64) bool {
	return !(value < startSize || value > endSize)
}

func FilterSliceByStrings(pkgStrings []string, targetStrings []string) bool {
	for _, pkgString := range pkgStrings {
		if FuzzyStrings(pkgString, targetStrings) {
			return true
		}
	}

	return false
}

func FuzzyStrings(pkgString string, targetStrings []string) bool {
	pkgString = strings.ToLower(pkgString)

	for _, targetString := range targetStrings {
		if strings.Contains(pkgString, targetString) {
			return true
		}
	}

	return false
}

func StrictStrings(pkgString string, targetStrings []string) bool {
	pkgString = strings.ToLower(pkgString)

	return slices.Contains(targetStrings, pkgString)
}

func RelationExists(relations []Relation) bool {
	return len(relations) > 0
}

func StringExists(pkgString string) bool {
	return pkgString != ""
}
