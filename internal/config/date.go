package config

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"yaylog/internal/consts"
)

type DateFilter struct {
	StartDate    time.Time
	EndDate      time.Time
	IsExactMatch bool
}

func parseDateFilter(dateFilterInput string) (DateFilter, error) {
	if dateFilterInput == "" {
		return DateFilter{}, nil
	}

	if dateFilterInput == ":" {
		return DateFilter{}, fmt.Errorf("invalid date filter: ':' must be accompanied by a date")
	}

	pattern := `^(\d{4}-\d{2}-\d{2})?(?::(\d{4}-\d{2}-\d{2})?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(dateFilterInput)
	isExactMatch := !strings.Contains(dateFilterInput, ":")

	if matches == nil {
		return DateFilter{}, fmt.Errorf("invalid date filter format: %q", dateFilterInput)
	}

	startDate, err := parseDateMatch(matches[1], time.Time{})
	if err != nil {
		return DateFilter{}, err
	}

	endDate, err := parseDateMatch(matches[2], time.Now())
	if err != nil {
		return DateFilter{}, err
	}

	return DateFilter{
		startDate,
		endDate,
		isExactMatch,
	}, nil
}

func parseDateMatch(dateInput string, defaultDate time.Time) (time.Time, error) {
	if dateInput == "" {
		return defaultDate, nil
	}

	return parseValidDate(dateInput)
}

func parseValidDate(dateInput string) (time.Time, error) {
	parsedDate, err := time.Parse(consts.DateOnlyFormat, dateInput)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}
