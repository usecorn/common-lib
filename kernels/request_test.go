package kernels

import (
	"fmt"
	"testing"

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
	rq := EarnRequest{
		UserAddr:   testutils.GenRandEVMAddr(),
		Source:     "source",
		StartBlock: 1,
		StartTime:  0,
		EarnRate:   "100",
	}

	tierRates := map[int]float64{
		0: 0.1,
		1: 0.2,
		2: 0.3,
		3: 0.4,
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

	require.EqualValues(t, fmt.Sprintf("%d", int(100*tierRates[0])), bonuses[0].EarnRate)
	require.EqualValues(t, fmt.Sprintf("%d", int(100*tierRates[1])), bonuses[1].EarnRate)

}
