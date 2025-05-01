package abi

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_AddressToBytes32(t *testing.T) {
	res := AddressToBytes32(common.HexToAddress("0x44f49ff0da2498bCb1D3Dc7C0f999578F67FD8C6"))
	require.Equal(t, 32, len(res))
	require.Equal(t, "00000000000000000000000044f49ff0da2498bcb1d3dc7c0f999578f67fd8c6", fmt.Sprintf("%x", res))
}
