package pipeline

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"yaylog/internal/consts"
)

func parseSizeFilter(sizeFilterInput string) (RangeSelector, error) {
	if sizeFilterInput == "" {
		return RangeSelector{}, nil
	}

	if sizeFilterInput == ":" {
		return RangeSelector{}, fmt.Errorf("invalid size filter: ':' must be accompanied by a value")
	}

	// valid size format: "10MB", "5GB:", ":20KB", "1.5MB:2GB" (value + unit, optional range)
	pattern := `(?i)^(?:(\d+(?:\.\d+)?)(B|KB|MB|GB))?(?::(?:(\d+(?:\.\d+)?)(B|KB|MB|GB))?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(sizeFilterInput)
	isExact := !strings.Contains(sizeFilterInput, ":")

	if matches == nil {
		return RangeSelector{}, fmt.Errorf("invalid size filter format: %q", sizeFilterInput)
	}

	start, err := parseSizeMatch(matches[1], matches[2], 0)
	if err != nil {
		return RangeSelector{}, err
	}

	end, err := parseSizeMatch(matches[3], matches[4], math.MaxInt64)
	if err != nil {
		return RangeSelector{}, err
	}

	return RangeSelector{
		start,
		end,
		isExact,
	}, nil
}

func parseSizeMatch(value string, unit string, defaultSize int64) (int64, error) {
	if value == "" {
		return defaultSize, nil
	}

	return parseSizeInBytes(value, unit)
}

func parseSizeInBytes(valueInput string, unitInput string) (sizeInBytes int64, err error) {
	value, err := strconv.ParseFloat(valueInput, 64) // parseFloat for fractional input e.g. ">2.5KB"
	if err != nil {
		return 0, fmt.Errorf("invalid size value")
	}

	unit := strings.ToUpper(unitInput)

	switch unit {
	case "KB":
		sizeInBytes = int64(value * consts.KB)
	case "MB":
		sizeInBytes = int64(value * consts.MB)
	case "GB":
		sizeInBytes = int64(value * consts.GB)
	case "B":
		sizeInBytes = int64(value)
	default:
		return 0, fmt.Errorf("invalid size unit: %v", unit)
	}

	return sizeInBytes, nil
}

func validateSizeFilter(sizeFilter RangeSelector) error {
	if sizeFilter.Start > 0 && sizeFilter.End > 0 {
		if sizeFilter.Start > sizeFilter.End {
			return fmt.Errorf("Error: invalid size range. Start size cannot be greater than the end size")
		}
	}

	return nil
}
