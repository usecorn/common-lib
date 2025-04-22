package eth

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HexToByte32(t *testing.T) {
	// Test cases
	tests := []string{
		"0x9039bf8b5c3cd6f2d3f937e8a2e59ef6af0109a0d0f3499e7dbf75be0aef75ec",
		"0x2547ba491a7ff9e8cfcaa3e1c0da739f4fdc1be9fe4a37bfcdf570002153a0de"}

	shortTests := []string{
		"0xda5ddd7270381a7c2717ad10d1c0ecb19e3cdfb2",
	}

	for _, test := range tests {
		result, err := HexToByte32(test)
		require.NoError(t, err)
		require.Len(t, result, 32, "HexToByte32 result length is not 32")
		resultHex := hex.EncodeToString(result[:])
		require.Equal(t, resultHex, test[2:], "HexToByte32 result does not match input")

	}

	for _, test := range shortTests {
		result, err := HexToByte32(test)
		require.NoError(t, err)
		require.Len(t, result, 32, "HexToByte32 result length is not 32")
		resultHex := hex.EncodeToString(result[:])
		require.Equal(t, strings.TrimRight(resultHex, "0"), test[2:], "HexToByte32 result should be zeroed for short input")
	}
}
