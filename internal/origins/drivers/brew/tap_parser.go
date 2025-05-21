package brew

import (
	"bytes"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	"strings"
	"unicode"
)

// TODO: let's remove the magic strings
// TODO: let's do byte operations here like we do in pacman/parser.go, we can beat the performance of strings.HasPrefix
func inferTapMetadata(pkg *pkgdata.PkgInfo, path string) {
	rubyPath := filepath.Join(path, ".brew", pkg.Name+".rb")
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
			case strings.HasPrefix(line, "desc "):
				pkg.Description = extractQuotedValue(line)

			case strings.HasPrefix(line, "homepage "):
				pkg.Url = extractQuotedValue(line)

			case strings.HasPrefix(line, "license "):
				if !strings.Contains(line, ":") && !strings.Contains(line, "[") {
					pkg.License = extractQuotedValue(line)
					break
				}
				inLicenseBlock = true
			}

			if inLicenseBlock {
				licenseLines = append(licenseLines, line)
				braceDepth += strings.Count(line, "[") + strings.Count(line, "{")
				braceDepth -= strings.Count(line, "]") + strings.Count(line, "}")

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
	if strings.HasPrefix(str, "(") && strings.HasSuffix(str, ")") {
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

	for _, r := range input {
		switch {
		case r == '"' && !inString:
			inString = true
			current.WriteRune(r)

		case r == '"' && inString:
			current.WriteRune(r)
			tokens = append(tokens, current.String())
			current.Reset()
			inString = false

		case inString:
			current.WriteRune(r)

		case unicode.IsSpace(r) || r == ',':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		case strings.ContainsRune("[]{}:", r):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))

		default:
			current.WriteRune(r)
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
		case "all_of":
			inner, offset := parseLicenseGroup(tokens[i+2:], "AND")
			parts = append(parts, "("+inner+")")
			i += offset + 2

		case "any_of":
			inner, offset := parseLicenseGroup(tokens[i+2:], "OR")
			parts = append(parts, "("+inner+")")
			i += offset + 2

		case "{":
			group, offset := parseLicenseTokens(tokens[i+1:])
			parts = append(parts, group)
			i += offset + 2

		case "[":
			i++

		case "]", "}", ":":
			return strings.Join(parts, " AND "), i

		default:
			if strings.HasPrefix(tokens[i], "\"") {
				parts = append(parts, strings.Trim(tokens[i], "\""))
			}
			i++
		}
	}

	return strings.Join(parts, " AND "), i
}

func parseLicenseGroup(tokens []string, joiner string) (string, int) {
	var items []string
	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case "{":
			expr, offset := parseLicenseTokens(tokens[i+1:])
			items = append(items, expr)
			i += offset + 2

		case "]":
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
