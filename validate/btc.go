package validate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	BtcAddrRegex = `^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,64}$`
)

var (
	btcAddrExp = regexp.MustCompile(BtcAddrRegex)
)

// IsTapRoot checks if a BTC address is a taproot address
func IsTapRoot(address string) bool {
	return strings.HasPrefix(strings.ToLower(address), "bc1p")
}

// GetValidBtcAddr returns a valid Bitcoin address or an error if the address is invalid.
func GetValidBtcAddr(addr string) (string, error) {
	if !btcAddrExp.MatchString(addr) {
		return "", errors.New("invalid bitcoin address")
	}
	return addr, nil
}
