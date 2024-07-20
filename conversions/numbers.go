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
	// the default precision is 53 bits, which is not enough to accurately handle 256-bit integers
	return big.NewFloat(0).SetPrec(SolidityIntSize)
}

func NumericToRat(n pgtype.Numeric) (*big.Rat, error) {
	rat := big.NewRat(1, 1)
	err := n.AssignTo(rat)
	if err != nil {
		return nil, err
	}
	return rat, nil
}

func NumericToFloat(n pgtype.Numeric) (*big.Float, error) {
	rat, err := NumericToRat(n)
	if err != nil {
		return nil, err
	}
	return NewLargeFloat().SetRat(rat), nil
}

func NumericToInt(n pgtype.Numeric) (*big.Int, error) {
	rat, err := NumericToRat(n)
	if err != nil {
		return nil, err
	}
	if rat.Denom().Sign() == 0 {
		return nil, errors.New("cannot convert infinite precision number to int")
	}
	return big.NewInt(0).Div(rat.Num(), rat.Denom()), nil
}

func MustNumericToInt(n pgtype.Numeric) *big.Int {
	out, err := NumericToInt(n)
	if err != nil {
		panic(err)
	}
	return out
}

func FloatToNumeric(f *big.Float) (pgtype.Numeric, error) {
	n := &pgtype.Numeric{}
	err := n.Set(f.String())
	return *n, err
}

func MustFloatToNumeric(f *big.Float) pgtype.Numeric {
	n, err := FloatToNumeric(f)
	if err != nil {
		panic(err)
	}
	return n
}

func RatToNumeric(r *big.Rat) (pgtype.Numeric, error) {
	return FloatToNumeric(big.NewFloat(0).SetRat(r))
}

func IntToNumeric(i *big.Int) (pgtype.Numeric, error) {
	n := &pgtype.Numeric{}
	err := n.Set(i.String())
	return *n, err
}

func MustIntToNumeric(i *big.Int) pgtype.Numeric {
	n, err := IntToNumeric(i)
	if err != nil {
		panic(err)
	}
	return n
}
