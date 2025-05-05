package deb

import (
	"bytes"
	"errors"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strconv"
	"strings"
)

func parseStatusFile(data []byte, origin string) ([]*pkgdata.PkgInfo, error) {
	reasonMap, err := loadInstallReasons()
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	var collectedErrors []error
	pkgs := []*pkgdata.PkgInfo{}

	for block := range bytes.SplitSeq(data, []byte("\n\n")) {
		if len(block) < 1 {
			continue
		}

		pkg, err := parseStatusBlock(block, reasonMap, origin)
		if err != nil {
			collectedErrors = append(collectedErrors, err)
		}

		pkgs = append(pkgs, pkg)
	}

	if len(collectedErrors) > 0 {
		return pkgs, errors.Join(collectedErrors...)
	}

	return pkgs, nil
}

func parseStatusFields(block []byte) map[string]string {
	fields := make(map[string]string)

	for line := range bytes.SplitSeq(block, []byte("\n")) {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(string(parts[0]))
		value := strings.TrimSpace(string(parts[1]))
		fields[key] = value
	}

	return fields
}

func parseStatusBlock(block []byte, reasonMap map[string]string, origin string) (*pkgdata.PkgInfo, error) {
	fields := parseStatusFields(block)
	var collected []error
	pkg := &pkgdata.PkgInfo{}
	meta := map[string]string{}

	for key, value := range fields {
		switch key {
		case fieldInstalledSize:
			size, err := strconv.Atoi(value)
			if err != nil {
				collected = append(collected, fmt.Errorf("invalid install size for %s: %v", pkg.Name, err))
				continue
			}
			pkg.Size = consts.KB * int64(size)

		case fieldPackage:
			pkg.Name = value

		case fieldVersion:
			pkg.Version = value

		case fieldArchitecture:
			pkg.Arch = value

		case fieldDescription:
			pkg.Description = value

		case fieldHomepage:
			pkg.Url = value

		case fieldMaintainer:
			pkg.Packager = value

		case fieldConflicts, fieldBreaks:
			pkg.Conflicts = append(pkg.Conflicts, parseRelations(value)...)

		case fieldReplaces:
			pkg.Replaces = parseRelations(value)

		case fieldDepends, fieldPreDepends:
			pkg.Depends = append(pkg.Depends, parseRelations(value)...)

		case fieldReccommends, fieldSuggests:
			pkg.OptDepends = append(pkg.OptDepends, parseRelations(value)...)

		case fieldProvides:
			pkg.Provides = parseRelations(value)

		case fieldPriority, fieldEssential:
			meta[key] = value
		}
	}

	if err := getInstallTime(pkg); err != nil {
		collected = append(collected, err)
	}
	_ = extractLicense(pkg)

	pkg.Origin = origin

	// TODO: for dpkg-only systems, perhaps return "unknown" for non-system packages
	if isSystem(meta) {
		pkg.Reason = "system"
	} else if reasonMap[pkg.Name] == "dependency" {
		pkg.Reason = "dependency"
	} else {
		pkg.Reason = "explicit"
	}

	return pkg, errors.Join(collected...)
}

func isSystem(meta map[string]string) bool {
	priority := strings.ToLower(meta[fieldPriority])
	essential := strings.ToLower(meta[fieldEssential])

	return priority == "required" || essential == "yes"
}
