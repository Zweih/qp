package brew

import (
	"path/filepath"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

type installedPkg struct {
	Name        string
	Version     string
	ReceiptPath string
	VersionPath string
	ModTime     int64
}

func fetchFormulae(
	origin string,
	prefix string,
	outChan chan<- *pkgdata.PkgInfo,
	errChan chan<- error,
	errGroup *sync.WaitGroup,
) {
	binRoot := filepath.Join(prefix, binSubPath)
	cellarRoot := filepath.Join(prefix, cellarSubPath)
	installedPkgs, err := getInstalledPkgs(cellarRoot, binRoot)
	if err != nil {
		errChan <- err
		return
	}

	if len(installedPkgs) < 1 {
		return
	}

	wanted := make(map[string]struct{}, len(installedPkgs))
	for _, iPkg := range installedPkgs {
		wanted[iPkg.Name] = struct{}{}
	}

	var formulaMeta map[string]*FormulaMetadata
	var metaErr error
	var metaWg sync.WaitGroup

	metaWg.Add(1)
	go func() {
		defer metaWg.Done()
		formulaMeta, metaErr = loadMetadata(formulaCachePath, getFormulaKey, wanted)
	}()

	inputChan := make(chan installedPkg, len(installedPkgs))
	for _, iPkg := range installedPkgs {
		inputChan <- iPkg
	}

	close(inputChan)

	stage1Out := worker.RunWorkers(
		inputChan,
		errChan,
		errGroup,
		func(iPkg installedPkg) (*pkgdata.PkgInfo, error) {
			return parseFormulaReceipt(iPkg.ReceiptPath, iPkg.Version)
		},
		0,
		len(installedPkgs),
	)

	stage2Out := worker.RunWorkers(
		stage1Out,
		errChan,
		errGroup,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			versionPath := filepath.Join(prefix, cellarSubPath, pkg.Name, pkg.Version)
			if size, err := getInstallSize(versionPath); err == nil {
				pkg.Size = size
			}

			return pkg, nil
		},
		0,
		len(installedPkgs),
	)

	metaWg.Wait()
	if metaErr != nil {
		errChan <- metaErr
		return
	}

	stage3Out := worker.RunWorkers(
		stage2Out,
		errChan,
		errGroup,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			if meta, ok := formulaMeta[pkg.Name]; ok {
				mergeFormulaMetadata(pkg, meta)
			}

			pkg.Origin = origin
			return pkg, nil
		},
		0,
		len(installedPkgs),
	)

	for pkg := range stage3Out {
		outChan <- pkg
	}
}
