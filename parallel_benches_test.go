package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
	"sync"
	"testing"
)

const (
	numTestCases = 120
	numChunks    = 1000
)

func hashTestCases(digests [][]byte, chunks [][]byte, threads int) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(digests))
	for i, digest := range digests {
		wg.Add(1)
		go func(d []byte, c []byte) {
			defer wg.Done()
			if err := HashParallel(d, c, threads); err != nil {
				errs <- err
			}
		}(digest, chunks[i])
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func splitAndHash(digests [][]byte, chunks [][]byte, tests int, threads int) error {
	for i := 0; i < len(digests); i += tests {
		end := i + tests
		if end > len(digests) {
			end = len(digests)
		}
		if err := hashTestCases(digests[i:end], chunks[i:end], threads); err != nil {
			return err
		}
	}
	return nil
}

func BenchmarkHash1D(b *testing.B) {
	var testCasesChunks [][]byte
	var testCasesDigests [][]byte
	for i := 0; i < numTestCases; i++ {
		chunks := make([]byte, numChunks<<6)
		_, err := rand.Read(chunks)
		if err != nil {
			b.Fatal(err)
		}
		testCasesChunks = append(testCasesChunks, chunks)
		digests := make([]byte, numChunks<<5)
		testCasesDigests = append(testCasesDigests, digests)
	}
	b.ResetTimer()

	for i := 1; i < numTestCases; i++ {
		if 120%i != 0 {
			continue
		}
		b.Run(fmt.Sprintf("Sending %d tests at a time", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				if err := splitAndHash(testCasesDigests, testCasesChunks, i, 0); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkHash2D(b *testing.B) {
	var testCasesChunks [][]byte
	var testCasesDigests [][]byte
	for i := 0; i < numTestCases; i++ {
		chunks := make([]byte, numChunks<<6)
		_, err := rand.Read(chunks)
		if err != nil {
			b.Fatal(err)
		}
		testCasesChunks = append(testCasesChunks, chunks)
		digests := make([]byte, numChunks<<5)
		testCasesDigests = append(testCasesDigests, digests)
	}
	numThreads := runtime.GOMAXPROCS(0)
	b.ResetTimer()

	for i := 1; i < numTestCases; i++ {
		if 120%i != 0 {
			continue
		}
		for t := 0; t < numThreads; t++ {
			b.Run(fmt.Sprintf("tests: %d threads: %d", i, t), func(b *testing.B) {
				for j := 0; j < b.N; j++ {
					if err := splitAndHash(testCasesDigests, testCasesChunks, i, t); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
