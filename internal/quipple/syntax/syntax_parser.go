package syntax

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/quipple"
	"qp/internal/quipple/compiler"
	"qp/internal/quipple/preprocess"
	"qp/internal/quipple/query"
	"strings"
)

type ParsedInput struct {
	Fields       []consts.FieldType
	FieldQueries []query.FieldQuery
	QueryExpr    compiler.Expr
	SortOption   SortOption
	OutputFormat string
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
	var formatMode string
	var whereTokens []string
	limitMode := LimitStart
	limit := 20

	currentBlock := quipple.BlockNone
	blockSeen := map[quipple.CmdType]bool{}

	for _, token := range preprocessedArgs {
		cmd := quipple.CmdTypeLookup[strings.ToLower(token)]
		if cmd != quipple.BlockNone {
			currentBlock = cmd

			if blockSeen[cmd] {
				return ParsedInput{}, fmt.Errorf("duplicate block: '%s'", quipple.CmdNameLookup[cmd])
			}
			blockSeen[cmd] = true

			continue
		}

		switch currentBlock {
		case quipple.BlockSelect:
			fieldTokens := strings.Split(token, ",")
			for _, fieldStr := range fieldTokens {
				fieldStr = strings.TrimSpace(fieldStr)

				fieldType, ok := consts.FieldTypeLookup[fieldStr]
				if !ok {
					return ParsedInput{}, fmt.Errorf("unknown field: %q", fieldStr)
				}

				fields = append(fields, fieldType)
			}

		case quipple.BlockWhere:
			whereTokens = append(whereTokens, token)

		case quipple.BlockOrder:
			sortOption, err = ParseSortOption(token)
			if err != nil {
				return ParsedInput{}, err
			}

		case quipple.BlockLimit:
			limit, limitMode, err = parseLimit(token)
			if err != nil {
				return ParsedInput{}, err
			}

		case quipple.BlockFormat:
			formatMode, err = parseFormat(token)
			if err != nil {
				return ParsedInput{}, err
			}

		default:
			return ParsedInput{}, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', 'order',  'limit', or 'format')", token)
		}
	}

	var queryExpr compiler.Expr
	if len(whereTokens) > 0 {
		expr, err := compiler.ParseExprBlock(whereTokens)
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
			Field: consts.FieldUpdated,
			Asc:   true,
		}
	}

	if formatMode == "" {
		formatMode = consts.OutputTable
	}

	return ParsedInput{
		Fields:       fields,
		FieldQueries: queries,
		QueryExpr:    queryExpr,
		SortOption:   sortOption,
		OutputFormat: formatMode,
		Limit:        limit,
		LimitMode:    limitMode,
	}, nil
}
