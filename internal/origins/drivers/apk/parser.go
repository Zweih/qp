package apk

import (
	"bytes"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/worker"
	"strconv"
	"strings"
	"sync"
)

func parseInstalledFile(
	data []byte,
	origin string,
	reasonMap map[string]bool,
) ([]*pkgdata.PkgInfo, error) {
	blocks := bytes.Split(data, []byte("\n\n"))

	inputChan := make(chan []byte, len(blocks))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, block := range blocks {
		if len(block) < 1 {
			continue
		}

		inputChan <- block
	}

	close(inputChan)

	resultChan := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(block []byte) (*pkgdata.PkgInfo, error) {
			return parsePackageBlock(block, origin, reasonMap)
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

func parsePackageBlock(
	block []byte,
	origin string,
	reasonMap map[string]bool,
) (*pkgdata.PkgInfo, error) {
	pkg := &pkgdata.PkgInfo{
		Origin: origin,
		Reason: consts.ReasonDependency, // default
	}

	var currentDir string
	var packageFiles []string

	for line := range bytes.SplitSeq(block, []byte("\n")) {
		if len(line) < 3 || line[1] != ':' {
			continue
		}

		key := string(line[0])
		value := strings.TrimSpace(string(line[2:]))
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
			pkg.Provides = parseRelations(value)
		case fieldReplaces:
			pkg.Replaces = parseRelations(value)
		case fieldRepoTag:
			if value != "" {
				pkg.Groups = []string{value}
			}
		case fieldDirectory:
			currentDir = value
		case fieldFile:
			packageFiles = append(packageFiles, filepath.Join("/", currentDir, value))
		}
	}

	pkg.UpdateTimestamp = findMostRecentModTime(packageFiles)
	pkg.InstallTimestamp = findOldestCreationTime(packageFiles)

	if _, exists := reasonMap[pkg.Name]; exists {
		pkg.Reason = consts.ReasonExplicit
	}

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

		// TODO: perhaps we should remove the "else" clause, too much nesting
		if idx := strings.Index(part, "="); idx != -1 {
			rel.Name = part[:idx]
			rel.Version = part[idx+1:]
			rel.Operator = pkgdata.OpEqual
		} else {
			rel.Name = part
		}

		rel.Name = stripApkPrefix(rel.Name)
		relations = append(relations, rel)
	}

	return relations
}

func stripApkPrefix(name string) string {
	if i := strings.Index(name, ":"); i != -1 {
		return name[i+1:]
	}

	return name
}

func findMostRecentModTime(files []string) int64 {
	var mostRecent int64 = 0

	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			modTime := info.ModTime().Unix()
			if modTime > mostRecent {
				mostRecent = modTime
			}
		}
	}

	return mostRecent
}

func findOldestCreationTime(files []string) int64 {
	if shared.InDocker() {
		return 0
	}

	var oldest int64 = 0

	for _, file := range files {
		if creationTime, reliable, err := shared.GetCreationTime(file); err == nil && reliable {
			if oldest == 0 || creationTime < oldest {
				oldest = creationTime
			}
		}
	}

	return oldest
}
