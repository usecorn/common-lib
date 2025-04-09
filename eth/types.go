package eth

import (
	"context"
	"database/sql"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	EthereumNetwork = "ethereum"
	CornMainnet     = "corn-mainnet"
)

const (
	EthereumChainID    = 1
	CornMainnetChainID = 21000000
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

func CalcSortIndex(blockNumber uint64, logIndex, txIndex uint) uint64 {
	// Test assumptions
	if blockNumber > 0xFFFFFFFF {
		panic("Block number is too large")
	}
	if txIndex > 0xFFFF {
		panic("Tx index is too large")
	}
	if logIndex > 0xFFFF {
		panic("Log index is too large")
	}
	// End test assumptions

	// Block number has the highest priority
	var out uint64 = 0xFFFFFFFF & blockNumber
	out = out << 16
	// Next is tx index
	out = out | (0xFFFF & uint64(txIndex))
	out = out << 16
	// Finally we have log index
	out = out | (0xFFFF & uint64(logIndex))
	return out
}

func CalcSortIndexFromLog(raw *types.Log) uint64 {
	return CalcSortIndex(raw.BlockNumber, raw.Index, raw.TxIndex)
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
	TXIndex     uint
}

func (et ERC20Transfer) SortIndex() uint64 {
	return CalcSortIndex(et.BlockNumber, et.LogIndex, et.TXIndex)
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

func SortTransfers(transfers []ERC20Transfer) {
	sort.Slice(transfers, func(i, j int) bool { // We need to handle these events in order.
		return transfers[i].SortIndex() < transfers[j].SortIndex()
	})
}
