package flatpak

import (
	"fmt"
	"path/filepath"
	"qp/internal/origins/shared"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

func fetchPackages(origin string, installDirs []string) ([]*pkgdata.PkgInfo, error) {
	estimatedCount := estimatePackageCount(installDirs)

	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	pkgRefChan := discoverPackages(installDirs, estimatedCount, errChan, &errGroup)

	stage1Out := worker.RunWorkers(
		pkgRefChan,
		errChan,
		&errGroup,
		func(pkgRef *PkgRef) (*PkgRef, error) {
			return parseMetadata(pkgRef)
		},
		0,
		estimatedCount,
	)

	stage2Out := worker.RunWorkers(
		stage1Out,
		errChan,
		&errGroup,
		func(pkgRef *PkgRef) (*PkgRef, error) {
			parseMetainfo(pkgRef)
			if pkgRef.Pkg.Title == "" {
				parseDesktopFile(pkgRef)
			}

			if err := applyTimestamps(pkgRef); err != nil {
				return nil, err
			}

			return pkgRef, nil
		},
		0,
		estimatedCount,
	)

	stage3Out := worker.RunWorkers(
		stage2Out,
		errChan,
		&errGroup,
		func(pkgRef *PkgRef) (*PkgRef, error) {
			deployPath := filepath.Join(pkgRef.CommitDir, "deploy")
			version, err := extractVersion(deployPath)
			if err != nil {
				fmt.Println(err)
			}

			pkgRef.Pkg.Version = version

			return pkgRef, nil
		},
		0,
		estimatedCount,
	)

	stage4Out := worker.RunWorkers(
		stage3Out,
		errChan,
		&errGroup,
		func(pkgRef *PkgRef) (*pkgdata.PkgInfo, error) {
			pkg := pkgRef.Pkg
			size, err := shared.GetInstallSize(pkgRef.CommitDir)
			if err != nil {
				return nil, err
			}

			pkg.PkgType = pkgRef.Type
			pkg.Env = fmt.Sprintf("%s (%s)", pkgRef.Scope, pkgRef.Branch)
			pkg.Size = size
			pkg.Origin = origin

			return pkg, nil
		},
		0,
		estimatedCount,
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(stage4Out, errChan)
}
