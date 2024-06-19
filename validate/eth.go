package validate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	EthAddrRegex = `^0x[0-9|a-f|A-F]{40}$`
	BtcAddrRegex = `^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,59}$`
)

var (
	ethAddrExp        = regexp.MustCompile(EthAddrRegex)
	ErrInvalidEthAddr = errors.New("invalid ethereum address")
)

func GetValidEthAddr(addr string) (string, error) {
	out := strings.ToLower(addr)
	if len(out) == 40 {
		out = "0x" + out // Add the 0x prefix if it's missing
	}
	if len(out) != 42 { // 0x + 40 characters
		return "", ErrInvalidEthAddr
	}
	if !ethAddrExp.MatchString(out) {
		return "", ErrInvalidEthAddr
	}
	return out, nil
}
