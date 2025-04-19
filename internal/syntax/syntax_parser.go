package syntax

import (
	"fmt"
	"qp/internal/consts"
	"strconv"
	"strings"
)

type CmdType int

const (
	BlockNone CmdType = iota
	BlockSelect
	BlockWhere
	BlockOrder
	BlockLimit
)

const (
	CmdSelect = "select"
	CmdWhere  = "where"
	CmdOrder  = "order"
	CmdLimit  = "limit"
)

type ParsedInput struct {
	Fields       []consts.FieldType
	FieldQueries []FieldQuery
	SortOption   SortOption
	Limit        int
}

func ParseSyntax(args []string) (ParsedInput, error) {
	preprocessedArgs, err := Preprocess(args)
	if err != nil {
		return ParsedInput{}, err
	}

	var fields []consts.FieldType
	var queries []FieldQuery
	var sortOption SortOption
	var whereTokens []string
	limit := 20

	currentBlock := BlockNone
	blockSeen := map[CmdType]bool{}

	for _, token := range preprocessedArgs {
		cmd := lookupCommand(token)
		if cmd != BlockNone {
			currentBlock = cmd

			if blockSeen[cmd] {
				return ParsedInput{}, fmt.Errorf("duplicate block: '%s'", cmdTypeName(cmd))
			}
			blockSeen[cmd] = true

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
			whereTokens = append(whereTokens, token)

		case BlockOrder:
			sortOption, err = ParseSortOption(token)
			if err != nil {
				return ParsedInput{}, err
			}

		case BlockLimit:
			limit, err = parseLimit(token)
			if err != nil {
				return ParsedInput{}, err
			}

		default:
			return ParsedInput{}, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', 'order', or 'limit')", token)
		}
	}

	if len(whereTokens) > 0 {
		parsedQueries, err := ParseQueriesBlock(whereTokens)
		if err != nil {
			return ParsedInput{}, err
		}

		queries = append(queries, parsedQueries...)
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
		Limit:        limit,
	}, nil
}

func parseLimit(token string) (int, error) {
	parsedInt, err := strconv.Atoi(token)
	if err != nil {
		return -2, fmt.Errorf("not a whole number (integer): %v", err)
	}

	return parsedInt, err
}

func lookupCommand(input string) CmdType {
	switch strings.ToLower(input) {
	case CmdSelect:
		return BlockSelect
	case CmdWhere:
		return BlockWhere
	case CmdOrder:
		return BlockOrder
	case CmdLimit:
		return BlockLimit
	default:
		return BlockNone
	}
}
