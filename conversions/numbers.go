package conversions

import (
	"errors"
	"math/big"

	"github.com/jackc/pgtype"
)

// SolidityIntSize is the size of a Solidity int in bits.
const SolidityIntSize uint = 256

// NewLargeFloat creates a new big.Float with a mantissa large enough to accurately represent the 256-bit integers used in Solidity.
func NewLargeFloat() *big.Float {
	// The default precision is 53 bits, which is not enough to accurately handle 256-bit integers.
	return big.NewFloat(0).SetPrec(SolidityIntSize)
}

// NumericToRat converts a pgtype.Numeric to a big.Rat.
func NumericToRat(n pgtype.Numeric) (*big.Rat, error) {
	rat := big.NewRat(1, 1)
	if err := n.AssignTo(rat); err != nil {
		return nil, errors.New("failed to convert Numeric to Rat: " + err.Error())
	}
	return rat, nil
}

// NumericToFloat converts a pgtype.Numeric to a big.Float.
func NumericToFloat(n pgtype.Numeric) (*big.Float, error) {
	rat, err := NumericToRat(n)
	if err != nil {
		return nil, err
	}
	return NewLargeFloat().SetRat(rat), nil
}

// NumericToInt converts a pgtype.Numeric to a big.Int, ensuring denominator is non-zero.
func NumericToInt(n pgtype.Numeric) (*big.Int, error) {
	rat, err := NumericToRat(n)
	if err != nil {
		return nil, err
	}
	if rat.Denom().Sign() == 0 {
		return nil, errors.New("cannot convert infinite precision number to int: denominator is zero")
	}
	return big.NewInt(0).Div(rat.Num(), rat.Denom()), nil
}

// MustNumericToInt converts a pgtype.Numeric to a big.Int and panics on error.
func MustNumericToInt(n pgtype.Numeric) *big.Int {
	out, err := NumericToInt(n)
	if err != nil {
		panic("MustNumericToInt failed: " + err.Error())
	}
	return out
}

// FloatToNumeric converts a big.Float to pgtype.Numeric.
func FloatToNumeric(f *big.Float) (pgtype.Numeric, error) {
	n := &pgtype.Numeric{}
	if err := n.Set(f.Text('f', -1)); err != nil {
		return pgtype.Numeric{}, errors.New("failed to convert Float to Numeric: " + err.Error())
	}
	return *n, nil
}

// MustFloatToNumeric converts a big.Float to pgtype.Numeric and panics on error.
func MustFloatToNumeric(f *big.Float) pgtype.Numeric {
	n, err := FloatToNumeric(f)
	if err != nil {
		panic("MustFloatToNumeric failed: " + err.Error())
	}
	return n
}

// RatToNumeric converts a big.Rat to pgtype.Numeric with precision.
func RatToNumeric(r *big.Rat) (pgtype.Numeric, error) {
	n := &pgtype.Numeric{}
	if err := n.Set(r.FloatString(20)); err != nil {
		return pgtype.Numeric{}, errors.New("failed to convert Rat to Numeric: " + err.Error())
	}
	return *n, nil
}

// IntToNumeric converts a big.Int to pgtype.Numeric.
func IntToNumeric(i *big.Int) (pgtype.Numeric, error) {
	n := &pgtype.Numeric{}
	if err := n.Set(i.String()); err != nil {
		return pgtype.Numeric{}, errors.New("failed to convert Int to Numeric: " + err.Error())
	}
	return *n, nil
}

// MustIntToNumeric converts a big.Int to pgtype.Numeric and panics on error.
func MustIntToNumeric(i *big.Int) pgtype.Numeric {
	n, err := IntToNumeric(i)
	if err != nil {
		panic("MustIntToNumeric failed: " + err.Error())
	}
	return n
}
