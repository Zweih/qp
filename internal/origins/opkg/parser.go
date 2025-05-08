package opkg

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/formats/debstyle"
	"qp/internal/pkgdata"
	"strconv"
)

func parseStatusFile(data []byte, origin string) ([]*pkgdata.PkgInfo, error) {
	var collectedErrors []error
	pkgs := []*pkgdata.PkgInfo{}

	for block := range bytes.SplitSeq(data, []byte("\n\n")) {
		if len(block) < 1 {
			continue
		}

		fields := debstyle.ParseStatusFields(block)
		if fields[fieldStatus] == "install ok installed" {
			continue
		}

		pkg, err := parseStatusBlock(fields, origin)
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

func parseStatusBlock(fields map[string]string, origin string) (*pkgdata.PkgInfo, error) {
	var collected []error
	pkg := &pkgdata.PkgInfo{}
	meta := map[string]string{}

	for key, value := range fields {
		switch key {
		case fieldInstalledTime:
			time, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				collected = append(collected, err)
				continue
			}

			pkg.InstallTimestamp = time

		case fieldPackage:
			pkg.Name = value

		case fieldAutoInstalled:
			if value == "yes" {
				pkg.Reason = consts.ReasonDependency
			}

		case fieldVersion:
			pkg.Version = value

		case fieldArchitecture:
			pkg.Arch = value

		case fieldConflicts:
			pkg.Conflicts = debstyle.ParseRelations(value)

		case fieldDepends:
			pkg.Depends = debstyle.ParseRelations(value)

		case fieldProvides:
			pkg.Provides = debstyle.ParseRelations(value)

		case fieldEssential:
			meta[key] = value
		}
	}

	if pkg.Reason == "" {
		pkg.Reason = consts.ReasonExplicit
	}

	extractMetadata(pkg)
	pkg.Origin = origin

	return pkg, errors.Join(collected...)
}

func extractMetadata(pkg *pkgdata.PkgInfo) error {
	controlPath := filepath.Join(opkgInfoRoot, pkg.Name+".control")

	data, err := os.ReadFile(controlPath)
	if err != nil {
		return fmt.Errorf("control file missing for %s: %w", pkg.Name, err)
	}

	var collected []error
	fields := debstyle.ParseStatusFields(data)

	for key, value := range fields {
		switch key {
		case fieldLicense:
			pkg.License = value

		case fieldDescription:
			pkg.Description = value

		case fieldInstalledSize:
			kb, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				collected = append(collected, err)
				continue
			}

			pkg.Size = kb * consts.KB

		case fieldSourceDateEpoch:
			time, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				collected = append(collected, fmt.Errorf("invalid Installed-Size for %s: %v", pkg.Name, err))
				continue
			}

			pkg.BuildTimestamp = time
		}
	}

	return errors.Join(collected...)
}
