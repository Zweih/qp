package syntax

import (
	"fmt"
	"strconv"
)

type LimitMode int32

const (
	LimitStart LimitMode = iota
	LimitMid
	LimitEnd
)

func parseLimit(token string) (int, LimitMode, error) {
	limitMode := LimitStart
	value := token
	var err error

	for i := range token {
		if token[i] == ':' {
			value, limitMode, err = parsePrefix(token, i)
			if err != nil {
				return -1, -1, err
			}

			break
		}
	}

	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		return -1, -1, fmt.Errorf("not a whole number (integer): %v", err)
	}

	return parsedInt, limitMode, err
}

func parsePrefix(token string, colonIdx int) (string, LimitMode, error) {
	prefix := token[:colonIdx]
	suffix := token[colonIdx+1:]

	switch prefix {
	case "mid":
		return suffix, LimitMid, nil
	case "end":
		return suffix, LimitEnd, nil
	default:
		return "", -1, fmt.Errorf("unknown limit prefix: %s", prefix)
	}
}
