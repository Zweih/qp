package pacman

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"qp/internal/pkgdata"
	"strconv"
	"strings"
)

func parseDescFile(descPath string) (*pkgdata.PkgInfo, error) {
	file, err := os.Open(descPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()

	// the average desc file is 103.13 lines, reading the entire file into memory is more efficient than using bufio.Scanner
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var pkg pkgdata.PkgInfo
	var currentField string
	start := 0
	end := 0
	length := len(data)

	for end <= length {
		if end == length || data[end] == '\n' {
			line := string(bytes.TrimSpace(data[start:end]))

			switch line {

			case fieldInstallDate, fieldBuildDate, fieldSize,
				fieldName, fieldReason, fieldVersion, fieldArch,
				fieldLicense, fieldPkgBase, fieldDescription,
				fieldUrl, fieldValidation, fieldPackager:
				currentField = line

			case fieldGroups, fieldConflicts, fieldReplaces,
				fieldDepends, fieldOptDepends, fieldProvides, fieldXData:
				currentField = line
				block, next := collectBlockBytes(data, end+1)

				applyMultiLineField(&pkg, currentField, block)
				end = next
				start = next

				continue

			case "":
				currentField = ""

			default:
				if err := applySingleLineField(&pkg, currentField, line); err != nil {
					return nil, fmt.Errorf("error reading desc file %s: %w", descPath, err)
				}
			}

			start = end + 1
		}

		end++
	}

	if pkg.Name == "" {
		return nil, fmt.Errorf("package name is missing in file: %s", descPath)
	}

	if pkg.Reason == "" {
		pkg.Reason = "explicit"
	}

	return &pkg, nil
}

func collectBlockBytes(data []byte, start int) ([]string, int) {
	var block []string
	i := start

	for i < len(data) {
		j := i

		for j < len(data) && data[j] != '\n' {
			j++
		}

		line := bytes.TrimSpace(data[i:j])

		if len(line) == 0 {
			break
		}

		block = append(block, string(line))
		i = j + 1
	}

	return block, i
}

func applySingleLineField(pkg *pkgdata.PkgInfo, field string, value string) error {
	switch field {
	case fieldInstallDate, fieldBuildDate, fieldSize:
		err := applyIntField(pkg, field, value)
		if err != nil {
			return err
		}

	case fieldName:
		pkg.Name = value

	case fieldReason:
		if value == "1" {
			pkg.Reason = "dependency"
		} else {
			pkg.Reason = "explicit"
		}

	case fieldVersion:
		pkg.Version = value

	case fieldArch:
		pkg.Arch = value

	case fieldLicense:
		pkg.License = value

	case fieldUrl:
		pkg.Url = value

	case fieldValidation:
		pkg.Validation = value

	case fieldDescription:
		pkg.Description = value

	case fieldPkgBase:
		pkg.PkgBase = value

	case fieldPackager:
		if value != "Unknown Packager" {
			pkg.Packager = value
		}

	default:
		// ignore unknown fields
	}
	return nil
}

func applyIntField(pkg *pkgdata.PkgInfo, field string, value string) error {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid %s value %q: %w", field, value, err)
	}

	switch field {
	case fieldInstallDate:
		pkg.InstallTimestamp = parsedValue
	case fieldBuildDate:
		pkg.BuildTimestamp = parsedValue
	case fieldSize:
		pkg.Size = parsedValue
	}

	return nil
}

func applyMultiLineField(pkg *pkgdata.PkgInfo, field string, block []string) {
	switch field {
	case fieldGroups:
		pkg.Groups = block
	case fieldConflicts, fieldReplaces,
		fieldDepends, fieldOptDepends, fieldProvides:
		applyRelations(pkg, field, block)
	case fieldXData:
		applyXData(pkg, block)
	}
}

func applyXData(pkg *pkgdata.PkgInfo, block []string) {
	for _, line := range block {
		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			subfield, value := parts[0], parts[1]

			switch subfield {
			case subfieldPkgType:
				pkg.PkgType = value
			}
		}
	}
}

func applyRelations(pkg *pkgdata.PkgInfo, field string, block []string) {
	relations := make([]pkgdata.Relation, 0, len(block))

	for _, line := range block {
		relations = append(relations, parseRelation(line))
	}

	switch field {
	case fieldConflicts:
		pkg.Conflicts = relations
	case fieldReplaces:
		pkg.Replaces = relations
	case fieldDepends:
		pkg.Depends = relations
	case fieldOptDepends:
		pkg.OptDepends = relations
	case fieldProvides:
		pkg.Provides = relations
	}
}

func parseRelation(input string) pkgdata.Relation {
	opStart := 0
	var depth int32 = 1

	for i := range input {
		switch input[i] {
		case '=', '<', '>':
			opStart = i
			goto parseVersion
		case ':':
			opStart = i
			goto parseWhy
		}
	}

	return pkgdata.Relation{Name: input, Depth: depth}

parseWhy:
	return pkgdata.Relation{
		Name:  input[:opStart],
		Why:   strings.TrimSpace(input[(opStart + 1):]),
		Depth: depth,
	}

parseVersion:
	name := input[:opStart]
	opEnd := opStart + 1

	if opEnd < len(input) {
		switch input[opEnd] {
		case '=', '<', '>':
			opEnd++
		}
	}

	operator := pkgdata.StringToOperator(input[opStart:opEnd])
	var version string

	if opEnd < len(input) {
		version = input[opEnd:]
	}

	return pkgdata.Relation{
		Name:     name,
		Operator: operator,
		Version:  version,
		Depth:    depth,
	}
}
