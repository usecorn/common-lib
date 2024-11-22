package validate

import (
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
)

const (
	EthAddrRegex = `^0x[0-9|a-f|A-F]{40}$`
)

var (
	EthAddrExp        = regexp.MustCompile(EthAddrRegex)
	ErrInvalidEthAddr = errors.New("invalid ethereum address")
)

// GetValidEthAddr returns a valid Ethereum address or an error if the address is invalid.
// will add the 0x prefix if it's missing, and lowercase the input.
func GetValidEthAddr(addr string) (string, error) {
	out := strings.ToLower(addr)
	if len(out) == 40 {
		out = "0x" + out // Add the 0x prefix if it's missing
	}
	if len(out) != 42 { // 0x + 40 characters
		return "", ErrInvalidEthAddr
	}
	if !EthAddrExp.MatchString(out) {
		return "", ErrInvalidEthAddr
	}
	return out, nil
}
