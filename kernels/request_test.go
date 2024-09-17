package kernels

import (
	"testing"

	"github.com/cornbase/common-lib/testutils"
	"github.com/stretchr/testify/require"
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
