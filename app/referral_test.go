package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewReferralCode(t *testing.T) {
	code, err := NewReferralCode()
	if err != nil {
		t.Errorf("error generating referral code: %v", err)
	}
	if len(code) != 9 {
		t.Errorf("expected referral code to be 9 characters, got %d", len(code))
	}
}

func Test_IsValidReferralCode(t *testing.T) {
	validCodes := []string{
		"ipm7-ffe3",
		"zyr7-cbn8",
	}

	for _, code := range validCodes {
		require.Truef(t, IsValidReferralCode(code), "expected %s to be valid", code)
	}
}
