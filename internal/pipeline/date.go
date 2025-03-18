package pipeline

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"yaylog/internal/consts"
)

func parseDateFilter(dateFilterInput string) (RangeFilter, error) {
	if dateFilterInput == "" {
		return RangeFilter{}, nil
	}

	if dateFilterInput == ":" {
		return RangeFilter{}, fmt.Errorf("invalid date filter: ':' must be accompanied by a date")
	}

	pattern := `^(\d{4}-\d{2}-\d{2})?(?::(\d{4}-\d{2}-\d{2})?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(dateFilterInput)
	isExact := !strings.Contains(dateFilterInput, ":")

	if matches == nil {
		return RangeFilter{}, fmt.Errorf("invalid date filter format: %q", dateFilterInput)
	}

	start, err := parseDateMatch(matches[1], 0)
	if err != nil {
		return RangeFilter{}, err
	}

	end, err := parseDateMatch(matches[2], time.Now().Unix())
	if err != nil {
		return RangeFilter{}, err
	}

	return RangeFilter{
		start,
		end,
		isExact,
	}, nil
}

func parseDateMatch(dateInput string, defaultDate int64) (int64, error) {
	if dateInput == "" {
		return defaultDate, nil
	}

	return parseValidDate(dateInput)
}

func parseValidDate(dateInput string) (int64, error) {
	parsedDate, err := time.Parse(consts.DateOnlyFormat, dateInput)
	if err != nil {
		return 0, err
	}

	return parsedDate.Unix(), nil
}

func validateDateFilter(dateFilter RangeFilter) error {
	if dateFilter.Start > 0 && dateFilter.End > 0 {
		if dateFilter.Start > dateFilter.End {
			return fmt.Errorf("Error invalid date range. The start date cannot be after the end date")
		}
	}

	return nil
}
