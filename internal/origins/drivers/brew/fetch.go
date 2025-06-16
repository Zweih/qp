package brew

import (
	"qp/internal/pkgdata"
	"qp/internal/worker"
	"sync"
)

func fetchPackages(origin string, prefix string) ([]*pkgdata.PkgInfo, error) {
	outChan := make(chan *pkgdata.PkgInfo)
	errChan := make(chan error, worker.DefaultBufferSize)

	var errGroup sync.WaitGroup
	var setupGroup sync.WaitGroup
	setupGroup.Add(2)

	go func() {
		defer setupGroup.Done()
		fetchFormulae(origin, prefix, outChan, errChan, &errGroup)
	}()

	go func() {
		defer setupGroup.Done()
		fetchCasks(origin, prefix, outChan, errChan, &errGroup)
	}()

	go func() {
		setupGroup.Wait()
		errGroup.Wait()
		close(outChan)
		close(errChan)
	}()

	return worker.CollectOutput(outChan, errChan)
}
