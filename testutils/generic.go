package testutils

import (
	"encoding/hex"
)

// GenRandHex generates a random hex string
func GenRandHex(n int) string {
	addrBytes := make([]byte, n/2)
	random.Read(addrBytes)
	return hex.EncodeToString(addrBytes)
}

// GenMany generates n items using the given function
func GenMany[T any](n int, fn func() T) []T {
	out := make([]T, n)
	for i := 0; i < n; i++ {
		out[i] = fn()
	}
	return out
}
