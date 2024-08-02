package validate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	BtcAddrRegex        = `^(bc1|[13]|tb1|[2mn])[a-zA-HJ-NP-Z0-9]{25,64}$`
	BtcTestnetAddrRegex = `^(tb1|[2mn])[a-zA-HJ-NP-Z0-9]{25,64}$`
	BtcMainnetRegex     = `^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,64}$`
)

var (
	btcAddrExp        = regexp.MustCompile(BtcAddrRegex)
	btcTestnetAddrExp = regexp.MustCompile(BtcTestnetAddrRegex)
	btcMainnetExp     = regexp.MustCompile(BtcMainnetRegex)
)

// IsTapRoot checks if a BTC address is a taproot address
func IsTapRoot(address string) bool {
	return strings.HasPrefix(strings.ToLower(address), "bc1p") || strings.HasPrefix(strings.ToLower(address), "tb1p")
}

// IsBitcoinTestnet checks if a BTC address is a testnet address
func IsBitcoinTestnet(address string) bool {
	return btcTestnetAddrExp.MatchString(strings.ToLower(address))
}

// IsBitcoinMainnet checks if a BTC address is a mainnet address
func IsBitcoinMainnet(address string) bool {
	return btcMainnetExp.MatchString(strings.ToLower(address))
}

// GetValidBtcAddr returns a valid Bitcoin address or an error if the address is invalid.
func GetValidBtcAddr(addr string) (string, error) {
	if !btcAddrExp.MatchString(addr) {
		return "", errors.New("invalid bitcoin address")
	}
	return addr, nil
}
