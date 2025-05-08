package deb

import (
	"bytes"
	"errors"
	"fmt"
	"qp/internal/consts"
	"qp/internal/origins/formats/debstyle"
	"qp/internal/pkgdata"
	"strconv"
	"strings"
)

func parseStatusFile(data []byte, origin string, reasonMap map[string]string) ([]*pkgdata.PkgInfo, error) {
	var collectedErrors []error
	pkgs := []*pkgdata.PkgInfo{}

	for block := range bytes.SplitSeq(data, []byte("\n\n")) {
		if len(block) < 1 {
			continue
		}

		fields := debstyle.ParseStatusFields(block)
		if fields[fieldStatus] != "install ok installed" {
			continue
		}

		pkg, err := parseStatusBlock(fields, reasonMap, origin)
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

func parseStatusBlock(
	fields map[string]string,
	reasonMap map[string]string,
	origin string,
) (*pkgdata.PkgInfo, error) {
	var collected []error
	pkg := &pkgdata.PkgInfo{}
	meta := map[string]string{}

	for key, value := range fields {
		switch key {
		case fieldInstalledSize:
			size, err := strconv.ParseInt(value, 10, 64)
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
			pkg.Conflicts = append(pkg.Conflicts, debstyle.ParseRelations(value)...)

		case fieldReplaces:
			pkg.Replaces = debstyle.ParseRelations(value)

		case fieldDepends, fieldPreDepends:
			pkg.Depends = append(pkg.Depends, debstyle.ParseRelations(value)...)

		case fieldRecommends, fieldSuggests:
			pkg.OptDepends = append(pkg.OptDepends, debstyle.ParseRelations(value)...)

		case fieldProvides:
			pkg.Provides = debstyle.ParseRelations(value)

		case fieldPriority, fieldEssential:
			meta[key] = value
		}
	}

	if err := getInstallTime(pkg); err != nil {
		collected = append(collected, err)
	}
	_ = extractLicense(pkg)

	pkg.Origin = origin

	switch {
	case isSystem(meta):
		pkg.Reason = consts.ReasonExplicit
	case reasonMap[pkg.Name] == consts.ReasonDependency:
		pkg.Reason = consts.ReasonDependency
	default:
		pkg.Reason = consts.ReasonExplicit
	}

	return pkg, errors.Join(collected...)
}

func isSystem(meta map[string]string) bool {
	priority := strings.ToLower(meta[fieldPriority])
	essential := strings.ToLower(meta[fieldEssential])

	return priority == "required" || essential == "yes"
}
