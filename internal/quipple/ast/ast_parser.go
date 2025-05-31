package ast

import (
	"fmt"
	"qp/internal/quipple/query"
)

func ParseExprBlock(rawTokens []string) (Expr, error) {
	if len(rawTokens) == 0 {
		return nil, nil
	}

	tokens, err := tokenizeWhereTokens(rawTokens)
	if err != nil {
		return nil, err
	}

	return parseExpr(tokens)
}

func parseExpr(tokens []QueryToken) (Expr, error) {
	pos := 0
	return parseOrExpr(tokens, &pos)
}

func parseUnaryExpr(tokens []QueryToken, pos *int) (Expr, error) {
	if *pos >= len(tokens) {
		return nil, fmt.Errorf("unexpected end of query")
	}

	switch tokens[*pos].Type {
	case TokenNot:
		*pos++
		inner, err := parseUnaryExpr(tokens, pos)
		if err != nil {
			return nil, err
		}

		return &NotExpr{Inner: inner}, nil
	case TokenOpenParen:
		return parseGroupedExpr(tokens, pos)

	case TokenQuery:
		return parseQueryExpr(tokens, pos)
	}

	return nil, fmt.Errorf("unexpected token: %s", tokens[*pos].Value)
}

func parseQueryExpr(tokens []QueryToken, pos *int) (Expr, error) {
	raw := tokens[*pos].Value
	*pos++
	fieldQuery, err := query.ParseQueryInput(raw)
	if err != nil {
		return nil, err
	}

	return &QueryExpr{Query: fieldQuery}, nil
}

func parseGroupedExpr(tokens []QueryToken, pos *int) (Expr, error) {
	*pos++
	group, err := parseOrExpr(tokens, pos)
	if err != nil {
		return nil, err
	}

	if *pos >= len(tokens) || tokens[*pos].Type != TokenCloseParen {
		return nil, fmt.Errorf("missing closing p to your q. Hint: q's are open parentheses `(` and p's are close parentheses `)`.")
	}

	*pos++
	return group, nil
}

func parseOrExpr(tokens []QueryToken, pos *int) (Expr, error) {
	left, err := parseAndExpr(tokens, pos)
	if err != nil {
		return nil, err
	}

	for *pos < len(tokens) && tokens[*pos].Type == TokenOr {
		*pos++
		right, err := parseAndExpr(tokens, pos)
		if err != nil {
			return nil, err
		}

		left = &OrExpr{Left: left, Right: right}
	}

	return left, nil
}

func parseAndExpr(tokens []QueryToken, pos *int) (Expr, error) {
	left, err := parseUnaryExpr(tokens, pos)
	if err != nil {
		return nil, err
	}

	for *pos < len(tokens) && tokens[*pos].Type == TokenAnd {
		*pos++
		right, err := parseUnaryExpr(tokens, pos)
		if err != nil {
			return nil, err
		}

		left = &AndExpr{Left: left, Right: right}
	}

	return left, nil
}
