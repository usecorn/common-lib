package kernels

import (
	"math/big"

	"github.com/pkg/errors"
)

type EarnRequest struct {
	UserAddr   string `json:"userAddr"`
	Source     string `json:"source"`
	SubSource  string `json:"subSource"`
	SourceUser string `json:"-"`
	StartBlock int64  `json:"startBlock"`
	StartTime  int64  `json:"startTime"`
	EarnRate   string `json:"earnRate"`
}

func (e EarnRequest) ReferralBonuses(referralChain []string, tierEarnRates map[int]float64) ([]EarnRequest, error) {
	var out []EarnRequest

	for i := range referralChain {
		req := EarnRequest{
			UserAddr:   referralChain[i],
			Source:     e.Source,
			SubSource:  e.SubSource,
			SourceUser: e.GetSourceUser(),
			StartBlock: e.StartBlock,
			StartTime:  e.StartTime,
		}
		earnRate, ok := big.NewFloat(0).SetString(e.EarnRate)
		if !ok {
			return nil, errors.New("invalid earn rate")
		}
		earnRate.Mul(earnRate, big.NewFloat(tierEarnRates[i]))

		req.EarnRate = earnRate.String()
		out = append(out, req)
	}
	return out, nil
}

func (e EarnRequest) GetSourceUser() string {
	if e.SourceUser == "" {
		return e.UserAddr
	}
	return e.SourceUser
}

func (e EarnRequest) Validate() error {
	if e.StartTime == 0 { // Start time is always required
		return ErrMissingStart
	}

	floatVal, ok := big.NewFloat(0).SetString(e.EarnRate)
	if !ok {
		return ErrInvalidEarnRate
	}

	if floatVal.Sign() < 0 {
		return ErrNegativeRate
	}

	if floatVal.IsInf() {
		return ErrEarnInf
	}

	if len(e.UserAddr) != 42 {
		return ErrInvalidUserAddr
	}
	if len(e.Source) == 0 {
		return ErrEmptySource
	}
	if len(e.SubSource) == 0 {
		return ErrEmptySubSource

	}

	return nil
}

func (e EarnRequest) IsPerBlock() bool {
	return e.StartBlock != 0 // If startBlock is set, then it's per block
}
