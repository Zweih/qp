package rpm

import (
	"strconv"
	"strings"
	"unicode"
)

func compareVersions(first, second string) int {
	epoch1, version1 := parseEpochAndVersion(first)
	epoch2, version2 := parseEpochAndVersion(second)

	if epoch1 < epoch2 {
		return -1
	}

	if epoch1 > epoch2 {
		return 1
	}

	return comparePostfix(version1, version2)
}

func parseEpochAndVersion(version string) (int, string) {
	if idx := strings.Index(version, ":"); idx != -1 {
		epochStr := version[:idx]
		if epoch, err := strconv.Atoi(epochStr); err == nil {
			return epoch, version[idx+1:]
		}
	}

	return 0, version
}

func comparePostfix(first string, second string) int {
	for first != "" || second != "" {
		first = skipJunk(first)
		second = skipJunk(second)

		if strings.HasPrefix(first, versionTilde) {
			if !strings.HasPrefix(second, versionTilde) {
				return -1
			}

			first = first[1:]
			second = second[1:]
			continue
		}

		if strings.HasPrefix(second, versionTilde) {
			return 1
		}

		if first == "" || second == "" {
			break
		}

		seg1, remaining1, isNum1 := extractSegment(first)
		seg2, remaining2, isNum2 := extractSegment(second)

		if isNum1 && !isNum2 {
			return 1
		}

		if !isNum1 && isNum2 {
			return -1
		}

		var cmp int
		if isNum1 {
			cmp = compareNumericSegments(seg1, seg2)
		} else {
			cmp = compareAlphaSegments(seg1, seg2)
		}

		if cmp != 0 {
			return cmp
		}

		first = remaining1
		second = remaining2
	}

	if first == "" && second == "" {
		return 0
	}

	if first != "" {
		return 1
	}

	return -1
}

func skipJunk(version string) string {
	for len(version) > 0 {
		char := rune(version[0])
		if isAlphaNumericOrTilde(char) {
			break
		}

		version = version[1:]
	}

	return version
}

func extractSegment(version string) (segment, remaining string, isNumeric bool) {
	if version == "" {
		return "", "", false
	}

	start := 0
	firstChar := rune(version[0])
	isNumeric = unicode.IsDigit(firstChar)

	for i, char := range version {
		if isNumeric {
			if !unicode.IsDigit(char) {
				return version[start:i], version[i:], true
			}
		} else {
			if !unicode.IsLetter(char) {
				return version[start:i], version[i:], false
			}
		}
	}

	return version[start:], "", isNumeric
}

func compareNumericSegments(segment1 string, segment2 string) int {
	segment1 = strings.TrimLeft(segment1, digitZero)
	segment2 = strings.TrimLeft(segment2, digitZero)

	if segment1 == "" {
		segment1 = digitZero
	}

	if segment2 == "" {
		segment2 = digitZero
	}

	len1, len2 := len(segment1), len(segment2)
	if len1 < len2 {
		return -1
	}

	if len1 > len2 {
		return 1
	}

	if segment1 < segment2 {
		return -1
	}

	if segment1 > segment2 {
		return 1
	}

	return 0
}

func compareAlphaSegments(segments1 string, segments2 string) int {
	if segments1 < segments2 {
		return -1
	}

	if segments1 > segments2 {
		return 1
	}

	return 0
}

func isAlphaNumericOrTilde(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '~'
}
