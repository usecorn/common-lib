package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EncodeAsInt24(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{-168180, "fd6f0c"},
		{-157200, "fd99f0"},
	}

	for _, test := range tests {
		result, err := EncodeAsInt24(test.value)
		require.NoError(t, err)
		require.Equal(t, test.expected, hex.EncodeToString(result))
	}
}
