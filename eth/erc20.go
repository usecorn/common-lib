package eth

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/usecorn/common-lib/eth/contracts"
)

type E20Cache interface {
	Set(key string, val interface{}) error
	GetString(key string) (out string, err error)
	GetInt64(key string) (int64, error)
	GetUint64(key string) (uint64, error)
}

type erc20 struct {
	log       logrus.Ext1FieldLogger
	erc20     *contracts.ERC20
	token     string
	ethClient *ethclient.Client
	metaDB    E20Cache
	decimals  int
	network   string
}

type ERC20 interface {
	TransferEvents(ctx context.Context, start, end uint64) ([]ERC20Transfer, error)
	BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
	Decimals(ctx context.Context) (int, error)
}

func NewERC20(log logrus.Ext1FieldLogger, metaDB E20Cache, ethClient *ethclient.Client, addr common.Address, network string) (ERC20, error) {
	erc20Contract, err := contracts.NewERC20(addr, ethClient)
	if err != nil {
		return nil, err
	}
	return &erc20{
		log:       log,
		erc20:     erc20Contract,
		token:     strings.ToLower(addr.Hex()),
		metaDB:    metaDB,
		decimals:  -1,
		ethClient: ethClient,
		network:   network,
	}, nil
}

func (et *erc20) BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	return et.erc20.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
}

func (et *erc20) TransferEvents(ctx context.Context, start, end uint64) ([]ERC20Transfer, error) {
	iter, err := et.erc20.FilterTransfer(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: ctx,
	}, nil, nil)

	if err != nil {
		return nil, err
	}
	out := []ERC20Transfer{}

	for iter.Next() {
		out = append(out, ERC20Transfer{
			From:        strings.ToLower(iter.Event.From.Hex()),
			To:          strings.ToLower(iter.Event.To.Hex()),
			Value:       iter.Event.Value,
			TXHash:      strings.ToLower(iter.Event.Raw.TxHash.String()),
			LogIndex:    iter.Event.Raw.Index,
			Token:       et.token,
			BlockNumber: iter.Event.Raw.BlockNumber,
			TXIndex:     iter.Event.Raw.TxIndex,
		})

	}

	return out, et.fillTimestamps(ctx, out)
}

func (et *erc20) fillTimestamps(ctx context.Context, events []ERC20Transfer) error {
	blockTimestamps := make(map[uint64]uint64)
	for _, event := range events {
		blockTimestamps[event.BlockNumber] = 0
	}

	for blockNum := range blockTimestamps {
		block, err := et.ethClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNum))
		if err != nil {
			return errors.Wrap(err, "failed to get block by number")
		}
		blockTimestamps[blockNum] = block.Time()
	}

	for i, event := range events {
		events[i].Timestamp = time.Unix(int64(blockTimestamps[event.BlockNumber]), 0)
	}
	return nil
}

func (et *erc20) Decimals(ctx context.Context) (int, error) {

	if et.decimals != -1 {
		return et.decimals, nil
	}
	var decimalsKey string
	if et.network == EthereumNetwork {
		decimalsKey = "erc20::" + et.token + "::decimals"
	} else {
		decimalsKey = "erc20::" + et.network + "::" + et.token + "::decimals"
	}

	decimals, err := et.metaDB.GetInt64(decimalsKey)
	if err == nil {
		et.decimals = int(decimals)
		return et.decimals, nil
	}

	val, err := et.erc20.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get decimals for %s", et.token)
	}
	err = et.metaDB.Set(decimalsKey, int64(val))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to set decimals for %s", et.token)
	}
	return int(val), nil
}
