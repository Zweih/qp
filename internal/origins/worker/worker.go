package worker

import (
	"errors"
	"runtime"
	"sync"
)

const DefaultBufferSize = 64

func RunWorkers[I any, O any](
	inputChan <-chan I,
	workerFunc func(I) (O, error),
	numWorkers int, // pass 0 unless testing or intentionally limiting
	bufferSize int,
) (<-chan O, <-chan error) {
	outputChan := make(chan O, bufferSize)
	errChan := make(chan error, DefaultBufferSize)

	if bufferSize == 0 {
		close(outputChan)
		close(errChan)
		return outputChan, errChan
	}

	if numWorkers <= 0 {
		// fun fact: NumCPU() does account for hyperthreading
		numWorkers = getWorkerCount(runtime.NumCPU())
	}

	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for item := range inputChan {
				pkg, err := workerFunc(item)
				if err != nil {
					errChan <- err
					continue
				}

				outputChan <- pkg
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputChan)
		close(errChan)
	}()

	return outputChan, errChan
}

func CollectOutput[O any](resultChan <-chan O, errChan <-chan error) ([]O, error) {
	var results []O
	var errs []error

	for result := range resultChan {
		results = append(results, result)
	}

	for err := range errChan {
		errs = append(errs, err)
	}

	var combinedErr error
	if len(errs) > 0 {
		combinedErr = errors.Join(errs...)
	}

	return results, combinedErr
}

func MergeErrors(errChans ...<-chan error) <-chan error {
	out := make(chan error, DefaultBufferSize)
	var wg sync.WaitGroup

	for _, errChan := range errChans {
		wg.Add(1)
		go func(eChan <-chan error) {
			defer wg.Done()

			for err := range eChan {
				out <- err
			}
		}(errChan)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func getWorkerCount(numCPUs int) int {
	numWorkers := numCPUs * 2
	if numCPUs <= 2 {
		// let's keep it simple for embedded devices
		numWorkers = numCPUs
	}

	return min(numWorkers, 12) // avoid overthreading on high-core systems
}
