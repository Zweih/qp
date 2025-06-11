package flatpak

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"
	"sync"
)

func fetchPackages(origin string, installDirs []string) ([]*pkgdata.PkgInfo, error) {
	pkgRefs, err := discoverPackages(installDirs)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	if len(pkgRefs) < 1 {
		return []*pkgdata.PkgInfo{}, nil
	}

	inputChan := make(chan *PkgRef, len(pkgRefs))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, pkgRef := range pkgRefs {
		inputChan <- pkgRef
	}
	close(inputChan)

	stage1Out := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(pkgRef *PkgRef) (*PkgRef, error) {
			return parseMetainfo(pkgRef)
		},
		0,
		len(pkgRefs),
	)

	stage2Out := worker.RunWorkers(
		stage1Out,
		errChan,
		&errGroup,
		func(pkgRef *PkgRef) (*pkgdata.PkgInfo, error) {
			pkg := pkgRef.Pkg
			fmt.Println(pkgRef.InstallDir)
			size, err := shared.GetInstallSize(pkgRef.InstallDir)
			if err != nil {
				return nil, err
			}

			pkg.Size = size
			pkg.Origin = origin

			return pkg, nil
		},
		0,
		len(pkgRefs),
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(stage2Out, errChan)
}

func parseMetainfo(pkgRef *PkgRef) (*PkgRef, error) {
	file, err := os.Open(pkgRef.MetadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s metadta file %s: %w", pkgRef.Name, pkgRef.MetadataPath, err)
	}

	defer file.Close()

	pkg := &pkgdata.PkgInfo{
		Name:   pkgRef.Name,
		Arch:   pkgRef.Arch,
		Reason: consts.ReasonExplicit,
	}

	scanner := bufio.NewScanner(file)

	var currentSection string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix("#", line) {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}

		if strings.Contains(line, "=") {
			key, value := parseKeyValue(line)
			applyMetadataField(pkg, currentSection, key, value)
			continue
		}
	}

	pkgRef.Pkg = pkg
	return pkgRef, nil
}

func parseKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func applyMetadataField(
	pkg *pkgdata.PkgInfo,
	section string,
	key string,
	value string,
) {
	switch section {
	case sectionApplication:
		applyApplicationField(pkg, key, value)
	default:
		// TODO: extension
	}
}

func applyApplicationField(pkg *pkgdata.PkgInfo, key string, value string) {
	switch key {
	case fieldRuntime:
		rel, err := parseRuntime(value)
		if err == nil {
			pkg.Depends = append(pkg.Depends, rel)
		}
	}
}

func parseRuntime(runtimeDir string) (pkgdata.Relation, error) {
	parts := filepath.SplitList(runtimeDir)
	if len(parts) < 3 {
		return pkgdata.Relation{}, fmt.Errorf("malformed runtime value: %s", runtimeDir)
	}

	name := parts[0]
	version := parts[len(parts)-1]

	return pkgdata.Relation{Name: name, Version: version}, nil
}
