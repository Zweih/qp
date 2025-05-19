package brew

import (
	"io/fs"
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
}

func fetchFormulae(
	origin string,
	prefix string,
) ([]*pkgdata.PkgInfo, error) {
	binRoot := filepath.Join(prefix, binSubPath)
	cellarRoot := filepath.Join(prefix, cellarSubPath)
	installedPkgs, err := getInstalledPkgs(cellarRoot, binRoot)
	if err != nil {
		return nil, err
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

	stage1Out, stage1Err := worker.RunWorkers(
		inputChan,
		func(iPkg installedPkg) (*pkgdata.PkgInfo, error) {
			return parseInstallReceipt(iPkg.ReceiptPath, iPkg.Version)
		},
		0,
		len(installedPkgs),
	)

	stage2Out, stage2Err := worker.RunWorkers(
		stage1Out,
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
		return nil, metaErr
	}

	stage3Out, stage3Err := worker.RunWorkers(
		stage2Out,
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

	allErrs := worker.MergeErrors(stage1Err, stage2Err, stage3Err)
	return worker.CollectOutput(stage3Out, allErrs)
}

func getInstallSize(dir string) (int64, error) {
	var total int64

	err := filepath.WalkDir(dir, func(_ string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !dir.IsDir() {
			info, err := dir.Info()
			if err != nil {
				return err
			}

			total += info.Size()
		}

		return nil
	})

	return total, err
}
