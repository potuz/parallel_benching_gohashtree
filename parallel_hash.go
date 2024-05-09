package main

import (
	"runtime"

	"github.com/prysmaticlabs/gohashtree"
)

func HashParallel(digests []byte, chunks []byte, threads int) error {
	if threads == 0 {
		threads = runtime.GOMAXPROCS(0)
	}
	numChunks := len(chunks) >> 6
	if threads < 2 || numChunks < 2 {
		return gohashtree.HashByteSlice(digests, chunks)
	}
	halfChunks := ((numChunks + 1) >> 1) << 6
	halfDigests := halfChunks >> 1
	go HashParallel(digests[:halfDigests], chunks[:halfChunks], threads/2)
	return HashParallel(digests[halfDigests:], chunks[halfChunks:], threads/2)
}
