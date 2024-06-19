package validate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var ethAddrExp = regexp.MustCompile(`^0x[0-9|a-f|A-F]{40}$`)

func GetValidEthAddr(addr string) (string, error) {
	out := strings.ToLower(addr)
	if len(out) == 40 {
		out = "0x" + out // Add the 0x prefix if it's missing
	}
	if len(out) != 42 { // 0x + 40 characters
		return "", errors.New("invalid ethereum address")
	}
	if !ethAddrExp.MatchString(out) {
		return "", errors.New("invalid ethereum address")
	}
	return out, nil
}
