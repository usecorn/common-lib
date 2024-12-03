package kernels

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usecorn/common-lib/testutils"
)

func Test_EarnRequestFullBatch_WithReferralBonuses(t *testing.T) {
	req := EarnRequestFullBatch{
		UserAddrs:   []string{testutils.GenRandEVMAddr()},
		Sources:     []string{"source"},
		SubSources:  []string{"subSource"},
		SourceUsers: nil,
		StartBlocks: nil,
		StartTimes:  []int64{1000},
		EarnRates:   []string{"1000"},
	}

	referralChains := [][]string{{testutils.GenRandEVMAddr(), testutils.GenRandEVMAddr()}}
	tierEarnRates := map[int]*big.Rat{0: big.NewRat(1, 2), 1: big.NewRat(1, 4)}

	result, err := req.WithReferralBonuses(referralChains, tierEarnRates)
	require.NoError(t, err)
	require.Len(t, result.UserAddrs, 3)

	for _, source := range result.Sources {
		require.Equal(t, req.Sources[0], source)
	}
	for _, subSource := range result.SubSources {
		require.Equal(t, req.SubSources[0], subSource)
	}
	for _, sourceUser := range result.SourceUsers {
		require.Equal(t, req.UserAddrs[0], sourceUser)
	}
	for _, startTime := range result.StartTimes {
		require.Equal(t, req.StartTimes[0], startTime)
	}
	require.Nil(t, result.StartBlocks)

	require.Equal(t, "1000", result.EarnRates[0])
	require.Equal(t, "500", result.EarnRates[1])
	require.Equal(t, "250", result.EarnRates[2])
}
