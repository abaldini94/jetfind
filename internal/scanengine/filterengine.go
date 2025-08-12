package scanengine

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

func FilterEngine(pathBuffer []ScanFilteredResult, scanFilter ScanFilter) []ScanFilteredResult {
	if len(pathBuffer) == 0 {
		return make([]ScanFilteredResult, 0)
	}

	const minPathsForParallel = 100
	if len(pathBuffer) < minPathsForParallel {
		filteredResults := make([]ScanFilteredResult, 0, len(pathBuffer)/4)
		for _, path := range pathBuffer {
			if p, filtered := scanFilter.Apply(path.Path); filtered {
				filteredResults = append(filteredResults, p)
			}
		}
		return filteredResults
	}

	var wg sync.WaitGroup

	filteredResults := make([]ScanFilteredResult, 0)
	filteredResultsChan := make(chan ScanFilteredResult, runtime.NumCPU())

	pathsPerWorker := len(pathBuffer)

	numWorkers := runtime.NumCPU()
	if len(pathBuffer) > numWorkers*10 {
		pathsPerWorker = int(math.Floor(float64(len(pathBuffer)) / float64(numWorkers)))
	}

	for i := range numWorkers {
		start := i * pathsPerWorker
		end := start + pathsPerWorker

		if start > len(pathBuffer) {
			break
		}

		if end > len(pathBuffer) {
			end = len(pathBuffer)
		}

		wg.Add(1)
		go func(paths []ScanFilteredResult) {
			defer wg.Done()
			for _, path := range paths {
				p, filtered := scanFilter.Apply(path.Path)
				if filtered {
					filteredResultsChan <- p
				}
			}

		}(pathBuffer[start:end])
	}

	go func() {
		wg.Wait()
		close(filteredResultsChan)
	}()

	for filteredPath := range filteredResultsChan {
		filteredResults = append(filteredResults, filteredPath)
	}

	sort.Slice(filteredResults, func(i, j int) bool {
		return filteredResults[i].Score > filteredResults[j].Score
	})

	return filteredResults
}
