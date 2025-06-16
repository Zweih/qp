package deb

import (
	"bytes"
	"errors"
	"fmt"
	"qp/internal/consts"
	"qp/internal/origins/shared/debstyle"
	"qp/internal/pkgdata"
	"qp/internal/worker"
	"strconv"
	"strings"
	"sync"
)

func parseStatusFile(data []byte, origin string, reasonMap map[string]string) ([]*pkgdata.PkgInfo, error) {
	blocks := bytes.Split(data, []byte("\n\n"))

	inputChan := make(chan map[string]string, len(blocks))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, block := range blocks {
		if len(block) < 1 {
			continue
		}

		fields := debstyle.ParseStatusFields(block)
		if fields[fieldStatus] != "install ok installed" {
			continue
		}

		inputChan <- fields
	}
	close(inputChan)

	resultChan := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(fields map[string]string) (*pkgdata.PkgInfo, error) {
			pkg, err := parseStatusBlock(fields, reasonMap, origin)
			return pkg, err
		},
		0,
		len(blocks),
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(resultChan, errChan)
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

		case fieldEnhances:
			pkg.OptionalFor = append(pkg.OptionalFor, debstyle.ParseRelations(value)...)

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
