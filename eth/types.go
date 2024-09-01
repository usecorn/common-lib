package eth

import (
	"context"
	"database/sql"
	"math/big"
	"time"

	gtypes "github.com/ethereum/go-ethereum/core/types"
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*gtypes.Block, error)
}

type ERC20Transfer struct {
	From        string
	To          string
	Value       *big.Int
	TXHash      string
	LogIndex    uint
	Token       string
	BlockNumber uint64
	Timestamp   time.Time
}

func (et ERC20Transfer) IsMint() bool {
	return et.From == "0x0000000000000000000000000000000000000000"
}

func (et ERC20Transfer) IsBurn() bool {
	return et.To == "0x0000000000000000000000000000000000000000"
}

func (et ERC20Transfer) SQLFrom() sql.NullString {
	if et.IsMint() {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: et.From, Valid: true}
}

func (et ERC20Transfer) SQLTo() sql.NullString {
	if et.IsBurn() {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: et.To, Valid: true}

}

type Result[T any] struct {
	Err error
	Val T
}
