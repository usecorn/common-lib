package eth

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/usecorn/common-lib/app"
	"github.com/usecorn/common-lib/server/config"
)

// CurrentSafeBlockHead returns the current block number minus the lag blocks. Additionally with retry if the RPC call fails.
func CurrentSafeBlockHead(ctx context.Context, conf *config.Chain, ethClient EthClient) (blockNum uint64, err error) {

	for range conf.RPCMaxRetries {
		blockNum, err = ethClient.BlockNumber(ctx)
		if err == nil {
			return blockNum - conf.LagBlocks, nil
		}
		err = app.SleepContext(ctx, conf.RPCRetryDelay)
		if err != nil {
			return
		}
	}

	return 0, err
}

// GetHeaderByNumberRetry retrieves the block header by number with retries, and will also attempt to redial the client if it fails.
// There are times when the connection itself is dropped, and redialing is necessary to get a new connection.
func GetHeaderByNumberRetry(ctx context.Context, network string, conf *config.Chain, client *ethclient.Client, blockNo int64) (*types.Header, error) {

	for i := range conf.RPCMaxRetries {
		header, err := client.HeaderByNumber(ctx, big.NewInt(blockNo))
		if err == nil {
			return header, nil
		}
		if i == conf.RPCMaxRetries-1 {
			return nil, errors.Wrap(err, "failed to get header by number")
		}

		err = app.SleepContext(ctx, conf.RPCRetryDelay)
		if err != nil {
			return nil, errors.Wrap(err, "sleep interrupted")
		}
		switch network {
		case EthereumNetwork:
			client, err = ethclient.DialContext(ctx, conf.RPCURL)
		case CornMainnet:
			client, err = ethclient.DialContext(ctx, conf.CornRPCURL)
		default:
			return nil, errors.Errorf("unknown network %s", network)
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to redial")
		}
	}
	return nil, errors.Errorf("failed to get block %d", blockNo)
}

func HexToByte32(hexString string) ([32]byte, error) {
	b, err := hex.DecodeString(strings.TrimLeft(hexString, "0x"))
	if err != nil {
		panic(err)
	}
	var out [32]byte
	if len(b) > 32 {
		return out, errors.Errorf("input hex string is too long: %d", len(b))
	}
	if len(b) < 32 {
		b = append(b, make([]byte, 32-len(b))...)
	}
	copy(out[:], b)
	return out, nil
}
