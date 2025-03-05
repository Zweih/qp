package config

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"yaylog/internal/consts"
)

type SizeFilter struct {
	StartSize    int64
	EndSize      int64
	IsExactMatch bool
}

func parseSizeFilter(sizeFilterInput string) (SizeFilter, error) {
	if sizeFilterInput == "" {
		return SizeFilter{}, nil
	}

	if sizeFilterInput == ":" {
		return SizeFilter{}, fmt.Errorf("invalid size filter: ':' must be accompanied by a value")
	}

	// valid size format: "10MB", "5GB:", ":20KB", "1.5MB:2GB" (value + unit, optional range)
	pattern := `(?i)^(?:(\d+(?:\.\d+)?)(B|KB|MB|GB))?(?::(?:(\d+(?:\.\d+)?)(B|KB|MB|GB))?)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(sizeFilterInput)
	isExactMatch := !strings.Contains(sizeFilterInput, ":")

	if matches == nil {
		return SizeFilter{}, fmt.Errorf("invalid size filter format: %q", sizeFilterInput)
	}

	startSize, err := parseSizeMatch(matches[1], matches[2], 0)
	if err != nil {
		return SizeFilter{}, err
	}

	endSize, err := parseSizeMatch(matches[3], matches[4], math.MaxInt64)
	if err != nil {
		return SizeFilter{}, err
	}

	return SizeFilter{
		startSize,
		endSize,
		isExactMatch,
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
