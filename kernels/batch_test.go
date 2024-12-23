package kernels

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usecorn/common-lib/testutils"
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
