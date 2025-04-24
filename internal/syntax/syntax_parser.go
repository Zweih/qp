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
	QueryExpr    Expr
	SortOption   SortOption
	Limit        int
	LimitMode    LimitMode
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
	limitMode := LimitStart
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
			limit, limitMode, err = parseLimit(token)
			if err != nil {
				return ParsedInput{}, err
			}

		default:
			return ParsedInput{}, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', 'order', or 'limit')", token)
		}
	}

	var queryExpr Expr
	if len(whereTokens) > 0 {
		expr, err := ParseExprBlock(whereTokens)
		if err != nil {
			return ParsedInput{}, err
		}
		queryExpr = expr
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
		QueryExpr:    queryExpr,
		SortOption:   sortOption,
		Limit:        limit,
		LimitMode:    limitMode,
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
	case CmdLimit:
		return BlockLimit
	default:
		return BlockNone
	}
}
