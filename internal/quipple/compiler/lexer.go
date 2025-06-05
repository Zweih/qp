package compiler

import (
	"fmt"
	"strings"
)

type TokenType int32

const (
	TokenQuery TokenType = iota
	TokenAnd
	TokenOr
	TokenNot
	TokenOpenParen
	TokenCloseParen
)

type QueryToken struct {
	Type  TokenType
	Value string
}

func tokenizeWhereTokens(rawTokens []string) ([]QueryToken, error) {
	var tokens []QueryToken
	lastWasQuery := false

	for _, raw := range rawTokens {
		switch strings.ToLower(raw) {
		case "and":
			tokens = append(tokens, QueryToken{Type: TokenAnd})
			lastWasQuery = false
		case "or":
			tokens = append(tokens, QueryToken{Type: TokenOr})
			lastWasQuery = false
		case "not":
			tokens = append(tokens, QueryToken{Type: TokenNot})
			lastWasQuery = false
		case "q":
			tokens = append(tokens, QueryToken{Type: TokenOpenParen})
			lastWasQuery = false
		case "p":
			tokens = append(tokens, QueryToken{Type: TokenCloseParen})
			lastWasQuery = false
		default:
			if lastWasQuery {
				return nil, fmt.Errorf(
					"missing logical operator between '%s' and '%s'. Use 'and', 'or', 'not', or group with 'q'/'p'",
					tokens[len(tokens)-1].Value, raw,
				)
			}

			tokens = append(tokens, QueryToken{
				Type:  TokenQuery,
				Value: raw,
			})
			lastWasQuery = true
		}
	}

	return tokens, nil
}
