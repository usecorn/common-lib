package abi

import (
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func AddressToBytes32(addr common.Address) []byte {
	out := make([]byte, 32)
	copy(out[12:], addr.Bytes())
	return out
}
