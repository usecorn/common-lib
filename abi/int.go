package abi

import "github.com/pkg/errors"

func EncodeAsInt24(value int64) ([]byte, error) {
	if value > 1<<24-1 || value < -(1<<24) {
		return nil, errors.Errorf("value out of range for int24: %d", value)
	}
	if value < 0 {
		return EncodeAsInt24(1<<24 + value)
	}
	return []byte{byte((value >> 16) & 0xFF), byte((value >> 8) & 0xFF), byte(value & 0xFF)}, nil
}
