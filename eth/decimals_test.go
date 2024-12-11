package eth

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToDecimalForm(t *testing.T) {
	type args struct {
		balance  *big.Int
		decimals int
	}

	tests := []args{
		{
			balance:  big.NewInt(1000000),
			decimals: 4,
		},
		{
			balance:  big.NewInt(1000000),
			decimals: 6,
		},
		{
			balance:  big.NewInt(1000000),
			decimals: 8,
		},
	}

	expectedResult := []*big.Float{
		big.NewFloat(100),
		big.NewFloat(1),
		big.NewFloat(0.01),
	}

	for i, test := range tests {
		result := ToDecimalForm(test.balance, test.decimals)
		require.EqualValues(t, expectedResult[i].String(), result.String())

	}
}

func Test_FromDecimalForm(t *testing.T) {
	type args struct {
		balance  *big.Float
		decimals int
	}

	tests := []args{
		{
			balance:  big.NewFloat(100),
			decimals: 4,
		},
		{
			balance:  big.NewFloat(1),
			decimals: 6,
		},
		{
			balance:  big.NewFloat(0.01),
			decimals: 8,
		},
	}

	expectedResult := []*big.Int{
		big.NewInt(1000000),
		big.NewInt(1000000),
		big.NewInt(1000000),
	}

	for i, test := range tests {
		result := FromDecimalForm(test.balance, test.decimals)
		require.EqualValues(t, expectedResult[i], result)

	}
}
