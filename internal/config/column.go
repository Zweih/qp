package config

import (
	"strings"
	"yaylog/internal/consts"
)

func parseColumns(columnsInput string, addColumnsInput string, hasAllColumns bool) []string {
	var specifiedColumnsRaw string
	var columns []string

	switch {
	case columnsInput != "":
		specifiedColumnsRaw = columnsInput
	case addColumnsInput != "":
		specifiedColumnsRaw = addColumnsInput
		fallthrough
	default:
		if hasAllColumns {
			columns = consts.ValidColumns
		} else {
			columns = consts.DefaultColumns
		}
	}

	specifiedColumns := strings.Split(strings.ToLower(strings.TrimSpace(specifiedColumnsRaw)), ",")
	for i, column := range specifiedColumns {
		specifiedColumns[i] = strings.TrimSpace(column)
	}

	columns = append(columns, specifiedColumns...)
	return cleanColumns(columns)
}

// in case of input like ",,"
func cleanColumns(columns []string) []string {
	cleanedColumns := []string{}
	for _, column := range columns {
		if column != "" {
			cleanedColumns = append(cleanedColumns, column)
		}
	}

	return cleanedColumns
}
