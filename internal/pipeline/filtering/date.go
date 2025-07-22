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

	pattern := `:(\d{4}-\d{2}-\d{2})`
	re := regexp.MustCompile(pattern)
	rangeMatch := re.FindStringIndex(dateFilterInput)

	var startStr, endStr string
	var isExact bool

	if rangeMatch != nil {
		isExact = false
		startStr = dateFilterInput[:rangeMatch[0]]
		endStr = dateFilterInput[rangeMatch[0]+1:]
	} else {
		isExact = true
		startStr = dateFilterInput
	}

	start, err := parseDateMatch(startStr, 0)
	if err != nil {
		return RangeSelector{}, err
	}

	end, err := parseDateMatch(endStr, time.Now().Unix())
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

func parseValidDate(dateInput string) (int64, error) {
	formatsNoTz := []string{
		consts.DateOnlyFormat,
		consts.DateTimeMinuteFormat,
		consts.DateTimeSecondFormat,
		consts.DateTime12HourMinuteFormat,
		consts.DateTime12HourSecondFormat,
	}

	for _, format := range formatsNoTz {
		if parsedDate, err := time.ParseInLocation(format, dateInput, time.Local); err == nil {
			return parsedDate.Unix(), nil
		}
	}

	formatsTz := []string{
		consts.DateOnlyTzFormat,
		consts.DateTimeMinuteTzFormat,
		consts.DateTimeSecondTzFormat,
		consts.DateTime12HourMinuteTzFormat,
		consts.DateTime12HourSecondTzFormat,
	}

	for _, format := range formatsTz {
		if parsedDate, err := time.Parse(format, dateInput); err == nil {
			return parsedDate.Unix(), nil
		}
	}

	return 0, fmt.Errorf("invalid date format: %q (supported formats: YYYY-MM-DD, YYYY-MM-DD HH:MM, YYYY-MM-DD HH:MM:SS, YYYY-MM-DD H:MM AM/PM, with optional timezone like MST or MDT).\nIf your format includes a space, you must surround it with quotes.", dateInput)
}

func parseDateMatch(dateInput string, defaultDate int64) (int64, error) {
	dateInput = strings.TrimSpace(dateInput)
	if dateInput == "" {
		return defaultDate, nil
	}

	return parseValidDate(dateInput)
}

func validateDateFilter(dateFilter RangeSelector) error {
	if dateFilter.Start > 0 && dateFilter.End > 0 {
		if dateFilter.Start > dateFilter.End {
			return fmt.Errorf("Error invalid date range. The start date cannot be after the end date")
		}
	}

	return nil
}
