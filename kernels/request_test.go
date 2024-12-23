package kernels

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/usecorn/common-lib/testutils"
)

func Test_EarnRequest_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		req  EarnRequest
		err  error
	}{
		{
			name: "startBlock and startTime both zero",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				SubSource:  "sub",
				StartBlock: 0,
				StartTime:  0,
				EarnRate:   "0.45",
			},
			err: ErrMissingStart,
		},
		{
			name: "startBlock and startTime both non-zero",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				SubSource:  "sub",
				StartBlock: 1,
				StartTime:  1,
				EarnRate:   "0.45",
			},
			err: nil,
		},
		{
			name: "startBlock non-zero",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				SubSource:  "sub",
				StartBlock: 1,
				StartTime:  1,
				EarnRate:   "0.45",
			},
			err: nil,
		},
		{
			name: "startTime non-zero",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				StartBlock: 0,
				SubSource:  "sub",
				StartTime:  333333,
				EarnRate:   "0.45",
			},
			err: nil,
		},
		{
			name: "negative earn rate",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				StartBlock: 0,
				SubSource:  "sub",
				StartTime:  1,
				EarnRate:   "-0.45",
			},
			err: ErrNegativeRate,
		},
		{
			name: "invalid user address",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr() + "f",
				Source:     "source",
				StartBlock: 0,
				StartTime:  1,
				EarnRate:   "0.45",
			},
			err: ErrInvalidUserAddr,
		},
		{
			name: "empty source",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "",
				StartBlock: 0,
				StartTime:  1,
				EarnRate:   "0.45",
			},
			err: ErrEmptySource,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			require.EqualValues(t, tt.err, err)
		})
	}
}

func Test_EarnRequest_IsPerBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		req  EarnRequest
		want bool
	}{
		{
			name: "startBlock non-zero",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				StartBlock: 1,
				StartTime:  0,
				EarnRate:   "0.45",
			},
			want: true,
		},
		{
			name: "startTime non-zero",
			req: EarnRequest{
				UserAddr:   testutils.GenRandEVMAddr(),
				Source:     "source",
				StartBlock: 0,
				StartTime:  1,
				EarnRate:   "0.45",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.IsPerBlock()
			require.EqualValues(t, tt.want, got)
		})
	}
}

func Test_EarnRequest_GetSourceUser(t *testing.T) {
	t.Parallel()
	user1 := testutils.GenRandEVMAddr()
	user2 := testutils.GenRandEVMAddr()
	tests := []struct {
		name   string
		req    EarnRequest
		result string
	}{
		{
			name: "get source user",
			req: EarnRequest{
				UserAddr:   user1,
				Source:     "source",
				StartBlock: 1,
				StartTime:  0,
				EarnRate:   "0.45",
			},
			result: user1,
		},
		{
			name: "get source user",
			req: EarnRequest{
				UserAddr:   user1,
				Source:     "source",
				SourceUser: user2,
				StartBlock: 1,
				StartTime:  0,
				EarnRate:   "0.45",
			},
			result: user2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.result, tt.req.GetSourceUser())
		})
	}
}

func Test_EarnRequest_ReferralBonuses(t *testing.T) {
	t.Parallel()
	tierRates := map[int]*big.Rat{
		0: big.NewRat(1, 10),
		1: big.NewRat(1, 20),
		2: big.NewRat(1, 30),
		3: big.NewRat(1, 40),
	}

	t.Run("happy path", func(t *testing.T) {
		rq := EarnRequest{
			UserAddr:   testutils.GenRandEVMAddr(),
			Source:     "source",
			StartBlock: 1,
			StartTime:  0,
			EarnRate:   "100",
		}

		referralChain := []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}

		bonuses, err := rq.ReferralBonuses(referralChain, tierRates)
		require.NoError(t, err)

		require.Len(t, bonuses, 2)

		for i := range bonuses {
			require.Equal(t, referralChain[i], bonuses[i].UserAddr)
			require.Equal(t, rq.GetSourceUser(), bonuses[i].SourceUser)
			require.Equal(t, rq.StartBlock, bonuses[i].StartBlock)
			require.Equal(t, rq.StartTime, bonuses[i].StartTime)
		}

		require.EqualValues(t, fmt.Sprintf("%d", int(100*1/tierRates[0].Denom().Int64())), bonuses[0].EarnRate)
		require.EqualValues(t, fmt.Sprintf("%d", int(100*1/tierRates[1].Denom().Int64())), bonuses[1].EarnRate)
	})

	t.Run("source user set", func(t *testing.T) {
		req := EarnRequest{
			EarnRate:   "100",
			UserAddr:   testutils.GenRandEVMAddr(),
			SourceUser: testutils.GenRandEVMAddr(),
			Source:     "ohio",
			SubSource:  "corn",
		}

		res, err := req.ReferralBonuses([]string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}, tierRates)
		require.NoError(t, err)
		require.Nil(t, res)
	})

}

func Test_GrantRequest_ReferralBonuses(t *testing.T) {
	tierRates := map[int]*big.Rat{
		0: big.NewRat(1, 10),
		1: big.NewRat(2, 10),
		2: big.NewRat(3, 10),
		3: big.NewRat(4, 10),
	}
	expected := []string{
		"10.00000000000000000000",
		"20.00000000000000000000",
		"30.00000000000000000000",
		"40.00000000000000000000",
	}
	t.Parallel()

	id := uuid.New()
	t.Run("negative amount", func(t *testing.T) {
		req := GrantRequest{
			UUID:     id,
			Amount:   "-100",
			UserAddr: testutils.GenRandEVMAddr(),
			Source:   "ohio",
			Category: "category",
		}

		res := req.ReferralBonuses([]string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}, tierRates)

		require.Nil(t, res)
	})

	t.Run("valid request", func(t *testing.T) {
		req := GrantRequest{
			UUID:     id,
			Amount:   "100",
			UserAddr: testutils.GenRandEVMAddr(),
			Source:   "kansas",
			Category: "category",
		}

		addrs := []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}

		res := req.ReferralBonuses(addrs, tierRates)

		require.Len(t, res, 2)

		for i := range res {
			require.NotEqual(t, req.UUID, res[i].UUID)
			require.Equal(t, addrs[i], res[i].UserAddr)
			require.Equal(t, expected[i], res[i].Amount)
		}
	})

	t.Run("uuids are stable", func(t *testing.T) {
		req := GrantRequest{
			UUID:     id,
			Amount:   "100",
			UserAddr: testutils.GenRandEVMAddr(),
			Source:   "arkansas",
		}

		addrs := []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}

		res := req.ReferralBonuses(addrs, tierRates)

		require.Len(t, res, 2)

		res2 := req.ReferralBonuses(addrs, tierRates)

		for i := range res {
			require.EqualValues(t, res[i].UUID, res2[i].UUID)
		}
	})

}

func Test_GrantRequest_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		req  GrantRequest
		err  error
	}{
		{
			name: "invalid user address",
			req: GrantRequest{
				UUID:      uuid.New(),
				UserAddr:  testutils.GenRandEVMAddr() + "5",
				Amount:    "100",
				GrantTime: 123214251,
				Category:  "category",
			},
			err: ErrInvalidUserAddr,
		},
		{
			name: "valid user address",
			req: GrantRequest{
				UUID:      uuid.New(),
				UserAddr:  testutils.GenRandEVMAddr(),
				Amount:    "100",
				Source:    "wyoming",
				Category:  "category",
				GrantTime: 123214251,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			require.EqualValues(t, tt.err, err)
		})
	}
}
