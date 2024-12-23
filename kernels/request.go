package kernels

import (
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/usecorn/common-lib/conversions"
	"github.com/usecorn/common-lib/validate"
)

type EarnRequest struct {
	UserAddr   string `json:"userAddr"`
	Source     string `json:"source"`
	SubSource  string `json:"subSource"`
	SourceUser string `json:"sourceUser"`
	StartBlock int64  `json:"startBlock"`
	StartTime  int64  `json:"startTime"`
	EarnRate   string `json:"earnRate"`
}

func (er EarnRequest) Clone() EarnRequest {
	var out EarnRequest
	err := copier.Copy(&out, &er)
	if err != nil {
		panic(err)
	}
	return out
}

func (e EarnRequest) ReferralBonuses(referralChain []string, tierEarnRates map[int]*big.Rat) ([]EarnRequest, error) {
	var out []EarnRequest

	if e.SourceUser == "" || e.SourceUser == e.UserAddr || len(referralChain) == 0 {
		return nil, nil
	}

	for i := range referralChain {
		req := EarnRequest{
			UserAddr:   referralChain[i],
			Source:     e.Source,
			SubSource:  e.SubSource,
			SourceUser: e.GetSourceUser(),
			StartBlock: e.StartBlock,
			StartTime:  e.StartTime,
		}
		earnRate, ok := conversions.NewLargeFloat().SetString(e.EarnRate)
		if !ok {
			return nil, errors.New("invalid earn rate")
		}

		earnRate.Mul(earnRate, conversions.NewLargeFloat().SetRat(tierEarnRates[i]))

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

type GrantRequest struct {
	UUID            uuid.UUID `json:"uuid"`
	UserAddr        string    `json:"userAddr"`
	Amount          string    `json:"amount"`
	Source          string    `json:"source"`
	SubSource       string    `json:"subSource"`
	SourceUser      string    `json:"-"`
	Category        string    `json:"category"`
	GrantTime       int64     `json:"grantTime"`
	ExcludeReferral bool      `json:"excludeReferral"`
}

func (gr GrantRequest) GetSourceUser() string {
	if gr.SourceUser == "" {
		return strings.ToLower(gr.UserAddr)
	}
	return strings.ToLower(gr.SourceUser)
}

func (gr GrantRequest) ReferralBonuses(referralChain []string, tierEarnRates map[int]*big.Rat) []GrantRequest {
	var out []GrantRequest

	parsedAmount, ok := big.NewRat(1, 1).SetString(gr.Amount)
	if !ok {
		// This should never happen
		panic("invalid amount: " + gr.Amount)
	}
	if parsedAmount.Sign() < 1 { // Referral penalties definitely shouldn't exist
		return nil
	}

	for i := range referralChain {
		multiplier := big.NewRat(1, 1).Set(tierEarnRates[i])
		req := GrantRequest{
			UUID:       uuid.NewSHA1(gr.UUID, []byte{byte(i >> 24 & 0xFF), byte(i >> 16 & 0xFF), byte(i >> 8 & 0xFF), byte(i & 0xFF)}),
			UserAddr:   referralChain[i],
			Amount:     multiplier.Mul(parsedAmount, multiplier).FloatString(20), // Points are accurate to 20 decimal places
			Source:     gr.Source,
			SourceUser: gr.GetSourceUser(),
		}
		out = append(out, req)
	}
	return out
}

func (g GrantRequest) Validate() error {
	if len(g.UserAddr) != 42 {
		return ErrInvalidUserAddr
	}

	if !validate.EthAddrExp.MatchString(g.UserAddr) {
		return ErrInvalidUserAddr
	}

	if len(g.Source) == 0 {
		return ErrEmptySource
	}
	if len(g.Category) == 0 {
		return ErrEmptyCategory
	}
	if g.GrantTime == 0 {
		return ErrMissingGrantTime
	}
	parsedAmount, ok := big.NewRat(1, 1).SetString(g.Amount)
	if !ok {
		return errors.Errorf("invalid amount: %s", g.Amount)
	}

	if parsedAmount.Sign() <= 0 {
		return ErrNegativeAmount
	}
	return nil

}
