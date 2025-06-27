package display

import (
	"bytes"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
)

type kvEntry struct {
	Key   string
	Value string
}

func (o *OutputManager) renderKeyValue(pkgs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	dateFormat := consts.DateTimeFormat
	ctx := tableContext{DateFormat: dateFormat}

	var buffer bytes.Buffer

	for i, pkg := range pkgs {
		entries := make([]kvEntry, len(fields))
		maxKeyLen := 0

		for _, field := range fields {
			value := getTableValue(pkg, field, ctx)
			if value == "" {
				continue
			}

			key := columnHeaders[field]
			if len(key) > maxKeyLen {
				maxKeyLen = len(key)
			}

			entries = append(entries, kvEntry{key, value})
		}

		for _, entry := range entries {
			lines := wrapKeyValue(entry.Key, entry.Value, maxKeyLen, o.terminalWidth)
			for _, line := range lines {
				buffer.WriteString(line + "\n")
			}
		}

		if i+1 < len(pkgs) {
			buffer.WriteString("\n")
		}
	}

	o.write(buffer.String())
}

func wrapKeyValue(key, value string, keyWidth, termWidth int) []string {
	prefix := fmt.Sprintf("%-*s : ", keyWidth, key)
	indent := strings.Repeat(" ", len(prefix))

	maxValueWidth := termWidth - len(prefix)
	if maxValueWidth <= 0 {
		maxValueWidth = 10
	}

	words := strings.Fields(value)
	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len()+len(word)+1 > maxValueWidth {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteByte(' ')
		}

		currentLine.WriteString(word)
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	output := make([]string, len(lines))
	for i, line := range lines {
		if i == 0 {
			output[i] = prefix + line
			continue
		}

		output[i] = indent + line
	}

	return output
}
