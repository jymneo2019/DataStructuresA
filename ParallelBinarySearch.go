package main

import (
	"sort"
	"sync"
	"sync/atomic"
)

// parallelBinarySearch performs binary search using multiple goroutines
func parallelBinarySearch(arr []int, target int, numGoroutines int) int {
	if len(arr) == 0 {
		return -1
	}

	// Ensure array is sorted
	sort.Ints(arr)

	// Calculate chunk size for each goroutine
	chunkSize := len(arr) / numGoroutines
	resultChan := make(chan int, numGoroutines)

	var wg sync.WaitGroup
	var found int32 = 0 // Atomic flag to signal early termination

	// Launch goroutines for parallel search
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Calculate start and end indices for this goroutine
			start := goroutineID * chunkSize
			end := start + chunkSize - 1
			if goroutineID == numGoroutines-1 {
				end = len(arr) - 1 // Last goroutine takes remaining elements
			}

			// Perform binary search in this segment
			for start <= end && atomic.LoadInt32(&found) == 0 {
				mid := start + (end-start)/2
				if arr[mid] == target {
					if atomic.CompareAndSwapInt32(&found, 0, 1) {
						resultChan <- mid // Send result if this is the first to find
					}
					return
				} else if arr[mid] < target {
					start = mid + 1
				} else {
					end = mid - 1
				}
			}
		}(i)
	}

	// Close result channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	if result, ok := <-resultChan; ok {
		return result
	}

	return -1
}
