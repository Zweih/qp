package pipx

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

type PkgMeta struct {
	Pkg          *pkgdata.PkgInfo
	DirPath      string
	SitePkgsPath string
}

func fetchPackages(venvRoot string, origin string) ([]*pkgdata.PkgInfo, error) {
	dirs, err := os.ReadDir(venvRoot)
	if err != nil {
		return []*pkgdata.PkgInfo{}, fmt.Errorf("failed to read pipx venv root: %w", err)
	}

	inputChan := make(chan os.DirEntry, len(dirs))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		inputChan <- dir
	}
	close(inputChan)

	stage1 := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(dir os.DirEntry) (*PkgMeta, error) {
			dirPath := filepath.Join(venvRoot, dir.Name())
			libPath := filepath.Join(dirPath, "lib")
			versionRoot, err := findVersionPath(libPath)
			if err != nil {
				return nil, fmt.Errorf("couldn't locate versioned root for %s: %v", dir.Name(), err)
			}

			sitePkgsPath := filepath.Join(versionRoot, "site-packages")
			metadataPath, err := findDistPath(sitePkgsPath, dir.Name())
			if err != nil {
				return nil, fmt.Errorf("couldn't locate metadata file for %s: %v", dir.Name(), err)
			}

			pkg, err := parseMetadataFile(metadataPath)
			if err != nil {
				return nil, fmt.Errorf("metadata parsing failed for %s: %v", metadataPath, err)
			}

			metaJson, err := os.Stat(filepath.Join(dirPath, "pipx_metadata.json"))
			if err != nil {
				return nil, fmt.Errorf("couldn't find pipx_metadata.json in %s: %w", dirPath, err)
			}

			pkg.UpdateTimestamp = metaJson.ModTime().Unix()
			pkg.Origin = origin
			pkg.Reason = consts.ReasonExplicit

			creationTime, reliable, err := shared.GetCreationTime(dirPath)
			if err == nil && reliable {
				pkg.InstallTimestamp = creationTime
			}

			return &PkgMeta{
				Pkg:          pkg,
				DirPath:      dirPath,
				SitePkgsPath: sitePkgsPath,
			}, nil
		},
		0,
		len(dirs),
	)

	stage2 := worker.RunWorkers(
		stage1,
		errChan,
		&errGroup,
		func(pMeta *PkgMeta) (*PkgMeta, error) {
			size, err := shared.GetInstallSize(pMeta.DirPath)
			if err != nil {
				return nil, err
			}

			pMeta.Pkg.Size = size
			return pMeta, nil
		},
		0,
		len(dirs),
	)

	stage3 := worker.RunWorkers(
		stage2,
		errChan,
		&errGroup,
		func(pMeta *PkgMeta) (*pkgdata.PkgInfo, error) {
			arch, err := inferArchitecture(pMeta.SitePkgsPath)
			if err != nil {
				return nil, err
			}

			pMeta.Pkg.Arch = arch
			return pMeta.Pkg, nil
		},
		0,
		len(dirs),
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(stage3, errChan)
}
