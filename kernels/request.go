package kernels

import (
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/usecorn/common-lib/validate"
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

func (er EarnRequest) Clone() EarnRequest {
	var out EarnRequest
	err := copier.Copy(&out, &er)
	if err != nil {
		panic(err)
	}
	return out
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

type GrantRequest struct {
	UUID            uuid.UUID `json:"uuid"`
	UserAddr        string    `json:"userAddr"`
	Amount          int64     `json:"amount"`
	Source          string    `json:"source"`
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

func (gr GrantRequest) ReferralBonuses(referralChain []string, tierEarnRates map[int]float64) []GrantRequest {
	var out []GrantRequest
	if gr.Amount <= 0 { // Referral penalties definitely shouldn't exist
		return nil
	}

	for i := range referralChain {
		req := GrantRequest{
			UUID:       uuid.NewSHA1(gr.UUID, []byte{byte(i >> 24 & 0xFF), byte(i >> 16 & 0xFF), byte(i >> 8 & 0xFF), byte(i & 0xFF)}),
			UserAddr:   referralChain[i],
			Amount:     int64(float64(gr.Amount) * tierEarnRates[i]),
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
	if g.Amount <= 0 {
		return ErrNegativeAmount
	}
	return nil

}
