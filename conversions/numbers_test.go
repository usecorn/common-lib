package conversions

import (
	"math/big"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/require"
)

func Test_NumericToRat(t *testing.T) {
	n := pgtype.Numeric{}
	err := n.Set("1.5")
	require.NoError(t, err)
	rat, err := NumericToRat(n)
	require.NoError(t, err)
	require.EqualValues(t, big.NewRat(3, 2), rat)
}

func Test_NumericToFloat(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		n := pgtype.Numeric{}
		err := n.Set("1.5")
		require.NoError(t, err)
		f, err := NumericToFloat(n)
		require.NoError(t, err)
		res, _ := f.Float64()
		require.EqualValues(t, 1.5, res)
	})
	t.Run("larger numbers", func(t *testing.T) {
		n := pgtype.Numeric{}
		err := n.Set("105111000000000000000001")
		require.NoError(t, err)
		f, err := NumericToFloat(n)
		require.NoError(t, err)
		require.True(t, f.IsInt())
		i, _ := f.Int(nil)
		require.EqualValues(t, "105111000000000000000001", i.String())
	})
}

func Test_NumericToInt(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		n := pgtype.Numeric{}
		err := n.Set("1")
		require.NoError(t, err)
		i, err := NumericToInt(n)
		require.NoError(t, err)
		require.EqualValues(t, big.NewInt(1), i)
	})

	t.Run("large number test", func(t *testing.T) {
		n := pgtype.Numeric{}
		err := n.Set("703000000000000000000000")
		require.NoError(t, err)
		i, err := NumericToInt(n)
		require.NoError(t, err)
		require.Equal(t, "703000000000000000000000", i.String())
	})
}

func Test_FloatToNumeric(t *testing.T) {
	f := big.NewFloat(1.5)
	n, err := FloatToNumeric(f)
	require.NoError(t, err)

	f2, err := NumericToFloat(n)
	require.NoError(t, err)
	require.EqualValues(t, f, f2)
}

func Test_RatToNumeric(t *testing.T) {
	r := big.NewRat(3, 2)
	n, err := RatToNumeric(r)
	require.NoError(t, err)

	r2, err := NumericToRat(n)
	require.NoError(t, err)
	require.EqualValues(t, r, r2)
}

func Test_IntToNumeric(t *testing.T) {
	i := big.NewInt(1)
	n, err := IntToNumeric(i)
	require.NoError(t, err)

	i2, err := NumericToInt(n)
	require.NoError(t, err)
	require.EqualValues(t, i, i2)

}
