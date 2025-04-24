package syntax

import (
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

	for _, raw := range rawTokens {
		switch strings.ToLower(raw) {
		case "and":
			tokens = append(tokens, QueryToken{Type: TokenAnd})
		case "or":
			tokens = append(tokens, QueryToken{Type: TokenOr})
		case "not":
			tokens = append(tokens, QueryToken{Type: TokenNot})
		case "q":
			tokens = append(tokens, QueryToken{Type: TokenOpenParen})
		case "p":
			tokens = append(tokens, QueryToken{Type: TokenCloseParen})
		default:
			tokens = append(tokens, QueryToken{
				Type:  TokenQuery,
				Value: raw,
			})
		}
	}

	return tokens, nil
}
