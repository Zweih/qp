package syntax

import (
	"fmt"
	"qp/internal/consts"
	"strings"
)

type SortOption struct {
	Field consts.FieldType
	Asc   bool
}

func ParseSortOption(sortInput string) (SortOption, error) {
	parts := strings.Split(sortInput, ":")
	fieldKey := strings.ToLower(parts[0])
	fieldType, exists := consts.FieldTypeLookup[fieldKey]
	if !exists {
		return SortOption{}, fmt.Errorf("invalid sort field: %s", fieldKey)
	}

	asc := true
	if len(parts) > 1 {
		switch parts[1] {
		case "desc":
			asc = false
		case "asc":
			asc = true
		default:
			return SortOption{}, fmt.Errorf("invalid sort direction: %s", parts[1])
		}
	}

	return SortOption{
		Field: fieldType,
		Asc:   asc,
	}, nil
}
