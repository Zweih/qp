package filtering

import (
	"fmt"
	"qp/internal/consts"
	"regexp"
	"strings"
	"time"
)

func parseDateFilter(dateFilterInput string) (RangeSelector, error) {
	if dateFilterInput == "" {
		return RangeSelector{}, nil
	}

	if dateFilterInput == ":" {
		return RangeSelector{}, fmt.Errorf("invalid date filter: ':' must be accompanied by a date")
	}

	pattern := `^(\d{4}-\d{2}-\d{2})?(?::(\d{4}-\d{2}-\d{2})?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(dateFilterInput)
	isExact := !strings.Contains(dateFilterInput, ":")

	if matches == nil {
		return RangeSelector{}, fmt.Errorf("invalid date filter format: %q", dateFilterInput)
	}

	start, err := parseDateMatch(matches[1], 0)
	if err != nil {
		return RangeSelector{}, err
	}

	end, err := parseDateMatch(matches[2], time.Now().Unix())
	if err != nil {
		return RangeSelector{}, err
	}

	end += int64(time.Hour * 24 / time.Second)

	return RangeSelector{
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

func validateDateFilter(dateFilter RangeSelector) error {
	if dateFilter.Start > 0 && dateFilter.End > 0 {
		if dateFilter.Start > dateFilter.End {
			return fmt.Errorf("Error invalid date range. The start date cannot be after the end date")
		}
	}

	return nil
}
