package testutils

import (
	"encoding/hex"
)

func GenRandHex(n int) string {
	addrBytes := make([]byte, n/2)
	random.Read(addrBytes)
	return hex.EncodeToString(addrBytes)
}
