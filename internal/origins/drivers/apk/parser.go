package apk

import (
	"bytes"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"qp/internal/worker"
	"strconv"
	"strings"
	"sync"
)

func parseInstalledFile(data []byte, origin string) ([]*pkgdata.PkgInfo, error) {
	blocks := bytes.Split(data, []byte("\n\n"))

	inputChan := make(chan map[string]string, len(blocks))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, block := range blocks {
		if len(block) < 1 {
			continue
		}

		fields := parseFields(block)
		if len(fields) > 0 {
			inputChan <- fields
		}
	}

	close(inputChan)

	resultChan := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(fields map[string]string) (*pkgdata.PkgInfo, error) {
			return parsePackageBlock(fields, origin)
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

func parseFields(block []byte) map[string]string {
	fields := make(map[string]string)

	for line := range bytes.SplitSeq(block, []byte("\n")) {
		if len(line) < 3 || line[1] != ':' {
			continue
		}

		key := string(line[0])
		value := strings.TrimSpace(string(line[2:]))

		fields[key] = value
	}

	return fields
}

func parsePackageBlock(fields map[string]string, origin string) (*pkgdata.PkgInfo, error) {
	pkg := &pkgdata.PkgInfo{
		Origin: origin,
		Reason: consts.ReasonExplicit, // TODO: default, may need inference
	}

	for key, value := range fields {
		switch key {
		case fieldBuildTime:
			if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
				pkg.BuildTimestamp = timestamp
			}
		case fieldInstalledSize:
			if size, err := strconv.ParseInt(value, 10, 64); err == nil {
				pkg.Size = size
			}
		case fieldPackage:
			pkg.Name = value
		case fieldVersion:
			pkg.Version = value
		case fieldArch:
			pkg.Arch = value
		case fieldDescription:
			pkg.Description = value
		case fieldUrl:
			pkg.Url = value
		case fieldLicense:
			pkg.License = value
		case fieldOrigin:
			pkg.PkgBase = value
		case fieldMaintainer:
			pkg.Packager = value
		case fieldDepends:
			pkg.Depends = parseRelations(value)
		case fieldProvides:
			// TODO: provides needs specific parsing
			pkg.Provides = parseRelations(value)
		case fieldReplaces:
			pkg.Replaces = parseRelations(value)
		case fieldRepoTag:
			if value != "" {
				pkg.Groups = []string{value}
			}
		}
	}

	// TODO: set update timestamp to file modification time

	return pkg, nil
}

func parseRelations(value string) []pkgdata.Relation {
	if value == "" {
		return nil
	}

	parts := strings.Fields(value)
	relations := make([]pkgdata.Relation, 0, len(parts))

	for _, part := range parts {
		rel := pkgdata.Relation{Depth: 1}

		if idx := strings.Index(part, "="); idx != -1 {
			rel.Name = part[:idx]
			rel.Version = part[idx+1:]
			rel.Operator = pkgdata.OpEqual
		} else {
			rel.Name = part
		}

		relations = append(relations, rel)
	}

	return relations
}
