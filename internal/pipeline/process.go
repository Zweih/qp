package pipeline

import (
	"fmt"
	"regexp"
	"strings"
	"yaylog/internal/config"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

type (
	PkgInfo         = pkgdata.PkgInfo
	FilterCondition = pkgdata.FilterCondition
)

var targetListRegex = regexp.MustCompile(`^([a-z0-9][a-z0-9._-]*[a-z0-9])(,([a-z0-9][a-z0-9._-]*[a-z0-9]))*$`)

func PreprocessFiltering(
	cfg config.Config,
	pkgPrts []*pkgdata.PkgInfo,
	reportProgress pkgdata.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if len(cfg.FilterQueries) == 0 {
		return pkgPrts, nil
	}

	filterConditions, err := queriesToConditions(cfg.FilterQueries)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	return pkgdata.FilterPackages(pkgPrts, filterConditions, reportProgress), nil
}

func queriesToConditions(filterQueries map[consts.FieldType]string) ([]pkgdata.FilterCondition, error) {
	conditions := make([]pkgdata.FilterCondition, 0, len(filterQueries))

	for fieldType, value := range filterQueries {
		var condition pkgdata.FilterCondition
		var err error

		switch fieldType {
		case consts.FieldDate:
			condition, err = parseDateFilterCondition(value)
		case consts.FieldSize:
			condition, err = parseSizeFilterCondition(value)
		case consts.FieldName,
			consts.FieldRequiredBy,
			consts.FieldDepends,
			consts.FieldProvides,
			consts.FieldConflicts,
			consts.FieldArch:
			condition, err = parsePackageFilterCondition(fieldType, value)
		case consts.FieldReason:
			condition, err = parseReasonFilterCondition(value)
		default:
			err = fmt.Errorf("unsupported filter type: %s", fieldType)
		}

		if err != nil {
			return []pkgdata.FilterCondition{}, err
		}

		conditions = append(conditions, condition)
	}

	return conditions, nil
}

func parsePackageFilterCondition(
	fieldType consts.FieldType,
	targetListInput string,
) (FilterCondition, error) {
	if !targetListRegex.MatchString(targetListInput) {
		return FilterCondition{}, fmt.Errorf("invalid package list: %s", targetListInput)
	}

	targetList := strings.Split(targetListInput, ",")
	return newPackageCondition(fieldType, targetList)
}

func parseReasonFilterCondition(installReason string) (FilterCondition, error) {
	if installReason != config.ReasonExplicit && installReason != config.ReasonDependency {
		return FilterCondition{}, fmt.Errorf("invalid install reason filter: %s", installReason)
	}

	return newReasonCondition(installReason), nil
}

func parseDateFilterCondition(value string) (FilterCondition, error) {
	dateFilter, err := parseDateFilter(value)
	if err != nil {
		return pkgdata.FilterCondition{}, fmt.Errorf("invalid date filter: %v", err)
	}

	if err = validateDateFilter(dateFilter); err != nil {
		return pkgdata.FilterCondition{}, err
	}

	return newDateCondition(dateFilter), nil
}

func parseSizeFilterCondition(value string) (FilterCondition, error) {
	sizeFilter, err := parseSizeFilter(value)
	if err != nil {
		return FilterCondition{}, fmt.Errorf("invalid size filter: %v", err)
	}

	if err = validateSizeFilter(sizeFilter); err != nil {
		return FilterCondition{}, err
	}

	return newSizeCondition(sizeFilter), nil
}
