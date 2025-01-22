package app

import (
	"crypto/rand"
	"regexp"

	"github.com/cockroachdb/errors"
)

var (
	ReferralCodeExp     = regexp.MustCompile(`^[3-9|a-h|j-k|m|n|p|r-t|x|y]{4}-[3-9|a-h|j-k|m|n|p|r-t|x|y]{4}$`)
	RootReferralCodeExp = regexp.MustCompile(`^z[3-9|a-h|j-k|m|n|p|r-t|x|y]{3}-[3-9|a-h|j-k|m|n|p|r-t|x|y]{4}$`)
	KOLReferralCodeExp  = regexp.MustCompile(`^i[3-9|a-h|j-k|m|n|p|r-t|x|y]{3}-[3-9|a-h|j-k|m|n|p|r-t|x|y]{4}$`)
)

// Note: Excludes ilqsvwz012
// Special: z, 0, i
const referralChars = "3456789abcdefghjkmnprtxy"

func NewReferralCode() (string, error) {

	// Referral codes
	out := ""

	rawCode := make([]byte, 8)
	n, err := rand.Read(rawCode)
	if n != 8 {
		return "", errors.New("failed to generate referral code")
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to generate referral code")
	}

	for i, b := range rawCode {
		if i == 4 {
			out += "-"
		}
		out += string(referralChars[b%uint8(len(referralChars))])
	}

	return out, nil
}

// IsValidReferralCode checks if a referral/root referral code is valid according to a regex pattern
func IsValidReferralCode(code string) bool {
	if IsRootReferralCode(code) {
		return RootReferralCodeExp.MatchString(code)
	}
	if IsKOLReferralCode(code) {
		return KOLReferralCodeExp.MatchString(code)
	}
	return ReferralCodeExp.MatchString(code)
}

func IsRootReferralCode(code string) bool {
	return len(code) == 9 && code[0] == 'z'
}

func IsKOLReferralCode(code string) bool {
	return len(code) == 9 && code[0] == 'i'
}

// NewRootReferralCode creates a new root referral code.
// root referral codes are similiar to referral codes, but the first character is
// always z. This means it will never validate as normal referral code.
// They are never created by a user and nobody gets a referral bonus from them.
func NewRootReferralCode() (string, error) {
	// First we generate a referral code
	code, err := NewReferralCode()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate root referral code")
	}
	code = "z" + code[1:]
	return code, nil
}

// NewKOLReferralCode creates a new kol referral code.
// This code is identical to a normal referral code, but the first character is
// always i. This means it will never validate as normal referral code.
func NewKOLReferralCode() (string, error) {
	// First we generate a referral code
	code, err := NewReferralCode()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate KOL referral code")
	}
	code = "i" + code[1:]
	return code, nil
}
