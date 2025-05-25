package rpm

import (
	"fmt"
	"qp/internal/pkgdata"
	"strings"
)

type Parser struct {
	lexer  *Lexer
	tokens []Token
	pos    int
}

func NewParser(input string) *Parser {
	lexer := NewLexer(input)
	var tokens []Token

	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == tokenEOF {
			break
		}
	}

	return &Parser{
		lexer:  lexer,
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) ParseExpression() (ExprNode, error) {
	return p.parseConditional()
}

func (p *Parser) parseConditional() (ExprNode, error) {
	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == tokenIf {
		p.advance()
		condition, err := p.parseOr()
		if err != nil {
			return nil, err
		}

		return &CondExpr{
			MainExpr: expr,
			Cond:     condition,
			IsUnless: false,
		}, nil
	}

	if p.current().Type == tokenUnless {
		p.advance()
		condition, err := p.parseOr()
		if err != nil {
			return nil, err
		}

		return &CondExpr{
			MainExpr: expr,
			Cond:     condition,
			IsUnless: true,
		}, nil
	}

	return expr, nil
}

func (p *Parser) parseOr() (ExprNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokenOr {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}

		left = &OrExpr{Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAnd() (ExprNode, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokenAnd {
		p.advance()
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		left = &AndExpr{Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (ExprNode, error) {
	token := p.current()

	switch token.Type {
	case tokenLParen:
		p.advance()
		expr, err := p.parseConditional()
		if err != nil {
			return nil, err
		}

		if p.current().Type != tokenRParen {
			return nil, fmt.Errorf("expected ')', got %v", p.current())
		}
		p.advance()
		return expr, nil

	case tokenPackage:
		p.advance()
		return parsePackageExpression(token.Value), nil

	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
	}
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: tokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func parsePackageExpression(packageStr string) *PkgExpr {
	operators := []string{opGreaterEqual, opLessEqual, opGreater, opLess, opEqual}

	// TODO: this should just split into 3
	for _, op := range operators {
		if idx := strings.Index(packageStr, " "+op+" "); idx != -1 {
			name := strings.TrimSpace(packageStr[:idx])
			version := strings.TrimSpace(packageStr[idx+len(op)+2:])
			return &PkgExpr{
				Name:     name,
				Version:  version,
				Operator: pkgdata.StringToOperator(op),
			}
		}
	}

	return &PkgExpr{
		Name:     packageStr,
		Operator: pkgdata.OpNone,
	}
}
