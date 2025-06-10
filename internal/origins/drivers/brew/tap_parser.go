package brew

import (
	"bytes"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	"strings"
	"unicode"
)

// TODO: let's do byte operations here like we do in pacman/parser.go, we can beat the performance of strings.HasPrefix
func inferTapMetadata(pkg *pkgdata.PkgInfo, path string) {
	rubyPath := filepath.Join(path, dotBrew, pkg.Name+dotRuby)
	data, err := os.ReadFile(rubyPath)
	if err != nil {
		return
	}

	var licenseLines []string
	inLicenseBlock := false
	braceDepth := 0

	start := 0
	end := 0
	length := len(data)

	for end <= length {
		if end == length || data[end] == '\n' {
			line := string(bytes.TrimSpace(data[start:end]))

			switch {
			case strings.HasPrefix(line, fieldDesc):
				pkg.Description = extractQuotedValue(line)

			case strings.HasPrefix(line, fieldHomepage):
				pkg.Url = extractQuotedValue(line)

			case strings.HasPrefix(line, fieldLicense):
				if !strings.Contains(line, ":") && !strings.Contains(line, openBracket) {
					pkg.License = extractQuotedValue(line)
					break
				}
				inLicenseBlock = true
			}

			if inLicenseBlock {
				licenseLines = append(licenseLines, line)
				braceDepth += strings.Count(line, openBracket) + strings.Count(line, openCurly)
				braceDepth -= strings.Count(line, closeBracket) + strings.Count(line, closeCurly)

				if braceDepth <= 0 {
					inLicenseBlock = false
					pkg.License = parseLicenseBlock(strings.Join(licenseLines, " "))
				}
			}

			start = end + 1
		}
		end++
	}
}

func extractQuotedValue(line string) string {
	parts := strings.SplitN(line, "\"", 3)
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

func parseLicenseBlock(input string) string {
	tokens := tokenizeLicenseBlock(input)
	expr, _ := parseLicenseTokens(tokens)
	return trimOuterParens(expr)
}

func trimOuterParens(str string) string {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, openParen) && strings.HasSuffix(str, closeParen) {
		depth := 0

		for i := range len(str) {
			switch str[i] {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 && i < len(str)-1 {
					return str
				}
			}
		}

		return str[1 : len(str)-1]
	}

	return str
}

func tokenizeLicenseBlock(input string) []string {
	var tokens []string
	var current strings.Builder

	inString := false

	for _, rune := range input {
		switch {
		case rune == '"' && !inString:
			inString = true
			current.WriteRune(rune)

		case rune == '"' && inString:
			current.WriteRune(rune)
			tokens = append(tokens, current.String())
			current.Reset()
			inString = false

		case inString:
			current.WriteRune(rune)

		case unicode.IsSpace(rune) || rune == ',':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		case strings.ContainsRune("[]{}:", rune):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(rune))

		default:
			current.WriteRune(rune)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func parseLicenseTokens(tokens []string) (string, int) {
	var parts []string

	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case fieldAllOf:
			inner, offset := parseLicenseGroup(tokens[i+2:], trimAnd)
			parts = append(parts, openParen+inner+closeParen)
			i += offset + 2

		case fieldAnyOf:
			inner, offset := parseLicenseGroup(tokens[i+2:], trimOr)
			parts = append(parts, openParen+inner+closeParen)
			i += offset + 2

		case openCurly:
			group, offset := parseLicenseTokens(tokens[i+1:])
			parts = append(parts, group)
			i += offset + 2

		case openBracket:
			i++

		case closeBracket, closeCurly, ":":
			return strings.Join(parts, spaceAnd), i

		default:
			if strings.HasPrefix(tokens[i], "\"") {
				parts = append(parts, strings.Trim(tokens[i], "\""))
			}
			i++
		}
	}

	return strings.Join(parts, spaceAnd), i
}

func parseLicenseGroup(tokens []string, joiner string) (string, int) {
	var items []string

	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case openCurly:
			expr, offset := parseLicenseTokens(tokens[i+1:])
			items = append(items, expr)
			i += offset + 2

		case closeBracket:
			return strings.Join(items, " "+joiner+" "), i + 1

		default:
			if strings.HasPrefix(tokens[i], "\"") {
				items = append(items, strings.Trim(tokens[i], "\""))
			}

			i++
		}
	}

	return strings.Join(items, " "+joiner+" "), i
}
