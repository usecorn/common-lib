package validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cornbase/common-lib/pkg/testutils"
)

func Test_GetValidEthAddr(t *testing.T) {

	validAddrWithCaps := strings.ToUpper(testutils.GenRandEVMAddr())

	parsed, err := GetValidEthAddr(validAddrWithCaps)
	require.NoError(t, err)
	require.Equal(t, strings.ToLower(validAddrWithCaps), parsed)

	validAddrLower := testutils.GenRandEVMAddr()

	parsed, err = GetValidEthAddr(validAddrLower)
	require.NoError(t, err)
	require.Equal(t, validAddrLower, parsed)

	validButNoPrefix := testutils.GenRandEVMAddr()[2:]

	parsed, err = GetValidEthAddr(validButNoPrefix)
	require.NoError(t, err)
	require.Equal(t, strings.ToLower("0x"+validButNoPrefix), parsed)

	invalidAddr := "INSERT ME INTO THE CORN, WE MUST JOIN THEM"
	_, err = GetValidEthAddr(invalidAddr)
	require.Error(t, err)
}
