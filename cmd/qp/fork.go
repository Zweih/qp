package main

import (
	"fmt"
	"os"
	"os/exec"
	"qp/api/driver"
	"qp/internal/config"
	"qp/internal/origins"
	"qp/internal/pipeline/phase"
	"qp/internal/storage"
	"sync"
)

func forkCacheWorker(originName string) error {
	cmd := exec.Command(os.Args[0], "--internal-cache-worker="+originName)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start background cache process: %w", err)
	}

	return nil
}

// TODO: add debouncing
func rebuildCache(originName string) error {
	cacheBasePath, err := storage.GetCachePath()
	if err != nil {
		return fmt.Errorf("failed to set up cache dir: %v", err)
	}

	drivers := origins.AvailableDrivers()
	if len(drivers) == 0 {
		return fmt.Errorf("no supported package origins detected")
	}

	var wg sync.WaitGroup
	cfg := &config.Config{
		NoCache:    false,
		RegenCache: false,
	}

	for _, d := range drivers {
		if originName != "all" && originName != d.Name() {
			continue
		}

		wg.Add(1)
		go func(targetDriver driver.Driver) {
			defer wg.Done()
			pipeline := phase.NewPipeline(targetDriver, cfg, false, cacheBasePath)

			if storage.IsLockFileExists(pipeline.CacheRoot) {
				return
			}

			if err := storage.CreateLockFile(pipeline.CacheRoot); err != nil {
				return
			}

			defer storage.RemoveLockFile(pipeline.CacheRoot)

			_, err := pipeline.RunCacheOnly()
			if err != nil {
				return
			}
		}(d)
	}

	wg.Wait()
	return nil
}
