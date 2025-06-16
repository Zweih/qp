package npm

import (
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/storage"
	"qp/internal/worker"
	"sync"
)

type NpmDriver struct {
	modulesDirs []string
}

func (d *NpmDriver) Name() string {
	return consts.OriginNpm
}

func (d *NpmDriver) Detect() bool {
	modulesDirs, err := getGlobalModulesDirs()
	if err != nil {
		return false
	}

	for _, modulesDir := range modulesDirs {
		if _, err := os.Stat(modulesDir); err != nil {
			continue
		}

		d.modulesDirs = append(d.modulesDirs, modulesDir)
	}

	return true
}

func (d *NpmDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	outChan := make(chan *pkgdata.PkgInfo)
	errChan := make(chan error, worker.DefaultBufferSize)

	var errGroup sync.WaitGroup
	var setupGroup sync.WaitGroup

	for _, modulesDir := range d.modulesDirs {
		setupGroup.Add(1)
		go func(dir string) {
			defer setupGroup.Done()
			fetchPackages(d.Name(), dir, outChan, errChan, &errGroup)
		}(modulesDir)
	}

	go func() {
		setupGroup.Wait()
		errGroup.Wait()
		close(outChan)
		close(errChan)
	}()

	return worker.CollectOutput(outChan, errChan)
}

func (d *NpmDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *NpmDriver) LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(cacheRoot)
}

func (d *NpmDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *NpmDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	var err error
	for _, modulesDir := range d.modulesDirs {
		isStale, err := shared.BfsStale(modulesDir, cacheMtime, 2)
		if isStale {
			return true, err
		}
	}

	return false, err
}
