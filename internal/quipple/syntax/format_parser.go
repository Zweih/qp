package syntax

import (
	"fmt"
	"qp/internal/consts"
)

func parseFormat(token string) (string, error) {
	switch token {
	case consts.OutputTable:
		return consts.OutputTable, nil
	case consts.OutputKeyValue:
		return consts.OutputKeyValue, nil
	case consts.OutputJSON:
		return consts.OutputJSON, nil
	default:
		return "", fmt.Errorf("invalid output format: %s", token)
	}
}
