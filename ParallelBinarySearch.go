package main

import (
	"fmt"
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

func main() {
	// Sample sorted array
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25}

	// Define the target value to search for
	target := 17

	// Number of goroutines to use for parallel search
	numGoroutines := 4

	// Call the parallelBinarySearch function
	result := parallelBinarySearch(arr, target, numGoroutines)

	// Print the result
	if result != -1 {
		fmt.Printf("Target %d found at index %d\n", target, result)
	} else {
		fmt.Printf("Target %d not found in the array\n", target)
	}
}
