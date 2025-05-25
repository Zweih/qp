package rpm

import "strings"

type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

type Lexer struct {
	input string
	pos   int
	start int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		start: 0,
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: tokenEOF, Pos: l.pos}
	}

	l.start = l.pos
	var token Token

	switch {
	case l.input[l.pos] == opOpenParen:
		l.pos++
		token = Token{Type: tokenLParen, Value: valueOpenParen, Pos: l.start}

	case l.input[l.pos] == opCloseParen:
		l.pos++
		token = Token{Type: tokenRParen, Value: valueCloseParen, Pos: l.start}

	case l.hasPrefix(opAnd):
		l.pos += len(opAnd)
		token = Token{Type: tokenAnd, Value: valueAnd, Pos: l.start}

	case l.hasPrefix(opOr):
		l.pos += len(opOr)
		token = Token{Type: tokenOr, Value: valueOr, Pos: l.start}

	case l.hasPrefix(opIf):
		l.pos += len(opIf)
		token = Token{Type: tokenIf, Value: valueIf, Pos: l.start}

	case l.hasPrefix(opUnless):
		l.pos += len(opUnless)
		token = Token{Type: tokenUnless, Value: valueUnless, Pos: l.start}

	default:
		token = l.scanPackage()
	}

	return token
}

func (l *Lexer) hasPrefix(prefix string) bool {
	end := l.pos + len(prefix)
	if end > len(l.input) {
		return false
	}

	return l.input[l.pos:end] == prefix
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t') {
		l.pos++
	}
}

func (l *Lexer) scanPackage() Token {
	start := l.pos

	for l.pos < len(l.input) {
		if l.input[l.pos] == opOpenParen || l.input[l.pos] == opCloseParen {
			break
		}

		if l.hasPrefix(opAnd) || l.hasPrefix(opOr) ||
			l.hasPrefix(opIf) || l.hasPrefix(opUnless) {
			break
		}

		l.pos++
	}

	value := strings.TrimSpace(l.input[start:l.pos])
	return Token{Type: tokenPackage, Value: value, Pos: start}
}
