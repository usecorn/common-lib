package validate

import "strings"

// IsTapRoot checks if a BTC address is a taproot address
func IsTapRoot(address string) bool {
	return strings.HasPrefix(strings.ToLower(address), "bc1p")
}
