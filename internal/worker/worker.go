package worker

import (
	"errors"
	"runtime"
	"sync"
)

const DefaultBufferSize = 64

var ErrSkip = errors.New("skip this item")

func RunWorkers[I any, O any](
	inputChan <-chan I,
	errChan chan<- error,
	errGroup *sync.WaitGroup,
	workerFunc func(I) (O, error),
	numWorkers int, // pass 0 unless testing or intentionally limiting
	bufferSize int,
) <-chan O {
	outputChan := make(chan O, bufferSize)

	if bufferSize == 0 {
		close(outputChan)
		return outputChan
	}

	if numWorkers <= 0 {
		// fun fact: NumCPU() does account for hyperthreading
		numWorkers = getWorkerCount(runtime.NumCPU())
	}

	var stageGroup sync.WaitGroup

	for range numWorkers {
		stageGroup.Add(1)
		errGroup.Add(1)

		go func() {
			defer stageGroup.Done()
			defer errGroup.Done()

			for item := range inputChan {
				result, err := workerFunc(item)
				if err != nil {
					if errors.Is(err, ErrSkip) {
						continue // skip silently
					}

					errChan <- err
					continue
				}

				outputChan <- result
			}
		}()
	}

	go func() {
		stageGroup.Wait()
		close(outputChan)
	}()

	return outputChan
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

func getWorkerCount(numCPUs int) int {
	numWorkers := numCPUs * 2
	if numCPUs <= 2 {
		// let's keep it simple for embedded devices
		numWorkers = numCPUs
	}

	return min(numWorkers, 12) // avoid overthreading on high-core systems
}
