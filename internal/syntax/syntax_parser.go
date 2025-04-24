package syntax

import (
	"fmt"
	"qp/internal/ast"
	"qp/internal/consts"
	"qp/internal/preprocess"
	"qp/internal/query"
	"strings"
)

type ParsedInput struct {
	Fields       []consts.FieldType
	FieldQueries []query.FieldQuery
	QueryExpr    ast.Expr
	SortOption   SortOption
	Limit        int
	LimitMode    LimitMode
}

func ParseSyntax(args []string) (ParsedInput, error) {
	preprocessedArgs, err := preprocess.Preprocess(args)
	if err != nil {
		return ParsedInput{}, err
	}

	var fields []consts.FieldType
	var queries []query.FieldQuery
	var sortOption SortOption
	var whereTokens []string
	limitMode := LimitStart
	limit := 20

	currentBlock := consts.BlockNone
	blockSeen := map[consts.CmdType]bool{}

	for _, token := range preprocessedArgs {
		cmd := consts.CmdTypeLookup[token]
		if cmd != consts.BlockNone {
			currentBlock = cmd

			if blockSeen[cmd] {
				return ParsedInput{}, fmt.Errorf("duplicate block: '%s'", consts.CmdNameLookup[cmd])
			}
			blockSeen[cmd] = true

			continue
		}

		switch currentBlock {
		case consts.BlockSelect:
			fieldTokens := strings.Split(token, ",")
			for _, fieldStr := range fieldTokens {
				fieldStr = strings.TrimSpace(fieldStr)

				fieldType, ok := consts.FieldTypeLookup[fieldStr]
				if !ok {
					return ParsedInput{}, fmt.Errorf("unknown field: %q", fieldStr)
				}

				fields = append(fields, fieldType)
			}

		case consts.BlockWhere:
			whereTokens = append(whereTokens, token)

		case consts.BlockOrder:
			sortOption, err = ParseSortOption(token)
			if err != nil {
				return ParsedInput{}, err
			}

		case consts.BlockLimit:
			limit, limitMode, err = parseLimit(token)
			if err != nil {
				return ParsedInput{}, err
			}

		default:
			return ParsedInput{}, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', 'order', or 'limit')", token)
		}
	}

	var queryExpr ast.Expr
	if len(whereTokens) > 0 {
		expr, err := ast.ParseExprBlock(whereTokens)
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
