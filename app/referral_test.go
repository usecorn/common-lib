package app

import "testing"

func Test_NewReferralCode(t *testing.T) {
	code, err := NewReferralCode()
	if err != nil {
		t.Errorf("error generating referral code: %v", err)
	}
	if len(code) != 9 {
		t.Errorf("expected referral code to be 9 characters, got %d", len(code))
	}
}
