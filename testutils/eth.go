package testutils

import (
	"encoding/hex"
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	// #nosec G404
	random = rand.New(rand.NewSource(time.Now().Unix()))
}

// GenRandEVMAddr generates a random Ethereum address
func GenRandEVMAddr() string {
	addrBytes := make([]byte, 20)
	random.Read(addrBytes)
	return "0x" + hex.EncodeToString(addrBytes)
}

func GenRandEVMHash() string {
	addrBytes := make([]byte, 32)
	random.Read(addrBytes)
	return "0x" + hex.EncodeToString(addrBytes)
}
