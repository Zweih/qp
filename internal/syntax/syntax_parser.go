package syntax

import (
	"fmt"
	"qp/internal/consts"
	"strings"
)

type CmdType int

const (
	BlockNone CmdType = iota
	BlockSelect
	BlockWhere
	BlockOrder
)

const (
	CmdSelect = "select"
	CmdWhere  = "where"
	CmdOrder  = "order"
)

type ParsedInput struct {
	Fields       []consts.FieldType
	FieldQueries []FieldQuery
	SortOption   SortOption
}

func ParseSyntax(args []string) (ParsedInput, error) {
	preprocessedArgs, err := Preprocess(args)
	if err != nil {
		return ParsedInput{}, err
	}

	var fields []consts.FieldType
	var queries []FieldQuery
	var sortOption SortOption

	currentBlock := BlockNone

	for _, token := range preprocessedArgs {
		cmd := lookupCommand(token)
		if cmd != BlockNone {
			currentBlock = cmd
			continue
		}

		switch currentBlock {
		case BlockSelect:
			fieldTokens := strings.Split(token, ",")
			for _, fieldStr := range fieldTokens {
				fieldStr = strings.TrimSpace(fieldStr)

				fieldType, ok := consts.FieldTypeLookup[fieldStr]
				if !ok {
					return ParsedInput{}, fmt.Errorf("unknown field: %q", fieldStr)
				}

				fields = append(fields, fieldType)
			}

		case BlockWhere:
			query, err := parseQueryInput(token)
			if err != nil {
				return ParsedInput{}, err
			}
			queries = append(queries, query)

		case BlockOrder:
			opt, err := ParseSortOption(token)
			if err != nil {
				return ParsedInput{}, err
			}
			sortOption = opt

		default:
			return ParsedInput{}, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', or 'order')", token)
		}
	}

	if len(fields) == 0 {
		fields = consts.DefaultFields
	}

	if sortOption == (SortOption{}) {
		sortOption = SortOption{
			Field: consts.FieldDate,
			Asc:   true,
		}
	}

	return ParsedInput{
		Fields:       fields,
		FieldQueries: queries,
		SortOption:   sortOption,
	}, nil
}

func lookupCommand(input string) CmdType {
	switch strings.ToLower(input) {
	case CmdSelect:
		return BlockSelect
	case CmdWhere:
		return BlockWhere
	case CmdOrder:
		return BlockOrder
	default:
		return BlockNone
	}
}
