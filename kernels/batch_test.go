package kernels

import (
	"math/big"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/usecorn/common-lib/testutils"
	"github.com/usecorn/common-lib/validate"
)

func Test_EarnRequestFullBatch_WithReferralBonuses(t *testing.T) {
	users := []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}
	req := EarnRequestFullBatch{
		UserAddrs:   []string{users[0], users[1]},
		Sources:     []string{"source", "source"},
		SubSources:  []string{"subSource", "subSource"},
		SourceUsers: []string{users[0], testutils.GenRandEVMAddr()},
		StartBlocks: nil,
		StartTimes:  []int64{1000, 2000},
		EarnRates:   []string{"1000", "2000"},
	}

	referralChains := [][]string{
		{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}, // This will not be ignored
		{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}} // Note we expect this to be ignored
	tierEarnRates := map[int]*big.Rat{0: big.NewRat(1, 2), 1: big.NewRat(1, 4)}

	result, err := req.WithReferralBonuses(referralChains, tierEarnRates)
	require.NoError(t, err)
	require.Len(t, result.UserAddrs, 4)

	for i, source := range result.Sources {
		if i == 1 {
			require.Equal(t, req.Sources[i], source)
		} else {
			require.Equal(t, req.Sources[0], source)
		}

	}
	for i, subSource := range result.SubSources {
		if i == 1 {
			require.Equal(t, req.SubSources[i], subSource)
		} else {
			require.Equal(t, req.SubSources[0], subSource)
		}

	}
	for i, sourceUser := range result.SourceUsers {
		if i == 1 {
			require.Equal(t, req.SourceUsers[i], sourceUser)
		} else {
			require.Equal(t, req.SourceUsers[0], sourceUser)
		}

	}
	for i, startTime := range result.StartTimes {
		if i == 1 {
			require.Equal(t, req.StartTimes[i], startTime)
		} else {
			require.Equal(t, req.StartTimes[0], startTime)
		}

	}
	require.Nil(t, result.StartBlocks)

	require.Len(t, result.EarnRates, 4)
	require.Equal(t, "1000", result.EarnRates[0])
	require.Equal(t, "500", result.EarnRates[2])
	require.Equal(t, "250", result.EarnRates[3])
}

func Test_EarnRequestBatch_IsPerBlock(t *testing.T) {
	tests := []struct {
		name     string
		batch    EarnRequestBatch
		expected bool
	}{
		{
			name: "with start block",
			batch: EarnRequestBatch{
				StartBlock: 100,
			},
			expected: true,
		},
		{
			name: "without start block",
			batch: EarnRequestBatch{
				StartBlock: 0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.batch.IsPerBlock()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_EarnRequestBatch_Size(t *testing.T) {
	tests := []struct {
		name     string
		batch    EarnRequestBatch
		expected int
	}{
		{
			name: "empty batch",
			batch: EarnRequestBatch{
				UserAddrs: []string{},
			},
			expected: 0,
		},
		{
			name: "batch with users",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.batch.Size()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_EarnRequestBatch_Validate(t *testing.T) {
	tests := []struct {
		name        string
		batch       EarnRequestBatch
		expectedErr error
	}{
		{
			name: "valid batch",
			batch: EarnRequestBatch{
				UserAddrs:  []string{testutils.GenRandEVMAddr()},
				EarnRates:  []string{"1.5"},
				Source:     "test",
				SubSource:  "unit",
				StartTime:  1000,
				StartBlock: 100,
			},
			expectedErr: nil,
		},
		{
			name: "missing start time",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"1.5"},
				Source:    "test",
				SubSource: "unit",
			},
			expectedErr: ErrMissingStart,
		},
		{
			name: "mismatched lengths",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"1.5", "2.0"},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
			},
			expectedErr: errors.New("userAddrs and earnRates must be the same length"),
		},
		{
			name: "invalid ethereum address",
			batch: EarnRequestBatch{
				UserAddrs: []string{"invalid-address"},
				EarnRates: []string{"1.5"},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
			},
			expectedErr: validate.ErrInvalidEthAddr,
		},
		{
			name: "invalid earn rate",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"invalid"},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
			},
			expectedErr: ErrInvalidEarnRate,
		},
		{
			name: "negative earn rate",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"-1.5"},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
			},
			expectedErr: ErrNegativeRate,
		},
		{
			name: "infinite earn rate",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"+inf"},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
			},
			expectedErr: ErrEarnInf,
		},
		{
			name: "negative start block",
			batch: EarnRequestBatch{
				UserAddrs:  []string{testutils.GenRandEVMAddr()},
				EarnRates:  []string{"1.5"},
				Source:     "test",
				SubSource:  "unit",
				StartTime:  1000,
				StartBlock: -1,
			},
			expectedErr: ErrNonPostiveStartBlock,
		},
		{
			name: "invalid start time",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"1.5"},
				Source:    "test",
				SubSource: "unit",
				StartTime: -1,
			},
			expectedErr: ErrNonPostiveStartTime,
		},
		{
			name: "empty source",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"1.5"},
				SubSource: "unit",
				StartTime: 1000,
			},
			expectedErr: ErrEmptySource,
		},
		{
			name: "empty subsource",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				EarnRates: []string{"1.5"},
				Source:    "test",
				StartTime: 1000,
			},
			expectedErr: ErrEmptySubSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.batch.Validate()
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.expectedErr.Error()), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_EarnRequestBatch_Clone(t *testing.T) {
	original := EarnRequestBatch{
		UserAddrs:   []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()},
		Source:      "test",
		SubSource:   "unit",
		SourceUsers: []string{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()},
		StartBlock:  100,
		StartTime:   1000,
		EarnRates:   []string{"1.5", "2.0"},
	}

	cloned := original.Clone()

	// Verify all fields are equal
	assert.Equal(t, original.UserAddrs, cloned.UserAddrs)
	assert.Equal(t, original.Source, cloned.Source)
	assert.Equal(t, original.SubSource, cloned.SubSource)
	assert.Equal(t, original.SourceUsers, cloned.SourceUsers)
	assert.Equal(t, original.StartBlock, cloned.StartBlock)
	assert.Equal(t, original.StartTime, cloned.StartTime)
	assert.Equal(t, original.EarnRates, cloned.EarnRates)

	// Verify it's a deep copy by modifying the clone
	cloned.UserAddrs[0] = testutils.GenRandEVMAddr()
	assert.NotEqual(t, original.UserAddrs[0], cloned.UserAddrs[0])
}

func TestEarnRequestBatch_WithReferralBonuses(t *testing.T) {
	tests := []struct {
		name           string
		batch          EarnRequestBatch
		referralChains [][]string
		tierEarnRates  map[int]*big.Rat
		expectedSize   int
		expectError    bool
	}{
		{
			name: "valid referral chain",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
				EarnRates: []string{"1.5"},
			},
			referralChains: [][]string{{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}},
			tierEarnRates: map[int]*big.Rat{
				0: big.NewRat(1, 2), // 50%
				1: big.NewRat(1, 4), // 25%
			},
			expectedSize: 3, // Original user + 2 referrals
			expectError:  false,
		},
		{
			name: "invalid earn rate",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
				EarnRates: []string{"invalid"},
			},
			referralChains: [][]string{{testutils.GenRandEVMAddr()}},
			tierEarnRates: map[int]*big.Rat{
				0: big.NewRat(1, 2),
			},
			expectError: true,
		},
		{
			name: "empty referral chain",
			batch: EarnRequestBatch{
				UserAddrs: []string{testutils.GenRandEVMAddr()},
				Source:    "test",
				SubSource: "unit",
				StartTime: 1000,
				EarnRates: []string{"1.5"},
			},
			referralChains: [][]string{{}},
			tierEarnRates: map[int]*big.Rat{
				0: big.NewRat(1, 2),
			},
			expectedSize: 1, // Only original user
			expectError:  false,
		},
		{
			name: "with source users",
			batch: EarnRequestBatch{
				UserAddrs:   []string{testutils.GenRandEVMAddr()},
				SourceUsers: []string{testutils.GenRandEVMAddr()},
				Source:      "test",
				SubSource:   "unit",
				StartTime:   1000,
				EarnRates:   []string{"1.5"},
			},
			referralChains: [][]string{{testutils.GenRandEVMAddr()}},
			tierEarnRates: map[int]*big.Rat{
				0: big.NewRat(1, 2),
			},
			expectedSize: 1, // No referral bonus due to different source user
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.batch.WithReferralBonuses(tt.referralChains, tt.tierEarnRates)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedSize, len(result.UserAddrs))
			assert.Equal(t, tt.expectedSize, len(result.EarnRates))
			assert.Equal(t, tt.expectedSize, len(result.SourceUsers))
		})
	}
}
