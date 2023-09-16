package rpc

import (
	"context"
	"ftm-explorer/internal/config"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	client "github.com/ethereum/go-ethereum/rpc"
)

// OperaRpc is a rpc client for Fantom Opera
type OperaRpc struct {
	ftm *client.Client
	// received blocks proxy
	wg       sync.WaitGroup
	sigClose chan struct{}
	headers  chan *types.Header
	// sfc contract address
	sfcAddress common.Address
	// closed flag
	closed bool
}

// NewOperaRpc returns a new rpc client for Fantom Opera
func NewOperaRpc(cfg *config.Rpc) (*OperaRpc, error) {
	ftm, err := client.Dial(cfg.OperaRpcUrl)
	if err != nil {
		return nil, err
	}
	return &OperaRpc{
		ftm:        ftm,
		sfcAddress: common.HexToAddress(cfg.SfcAddress),
	}, nil
}

// PendingNonceAt returns the nonce of the account at the given block.
func (rpc *OperaRpc) PendingNonceAt(ctx context.Context, address common.Address) (uint64, error) {
	return ethclient.NewClient(rpc.ftm).PendingNonceAt(ctx, address)
}

// SuggestGasPrice suggests a gas price.
func (rpc *OperaRpc) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return ethclient.NewClient(rpc.ftm).SuggestGasPrice(ctx)
}

// NetworkID returns the network ID.
func (rpc *OperaRpc) NetworkID(ctx context.Context) (*big.Int, error) {
	return ethclient.NewClient(rpc.ftm).NetworkID(ctx)
}

// Close closes the RPC client.
func (rpc *OperaRpc) Close() {
	if rpc.closed {
		return
	}

	if rpc.headers != nil {
		rpc.sigClose <- struct{}{}
		rpc.wg.Wait()
		close(rpc.headers)
		close(rpc.sigClose)
	}

	if rpc.ftm != nil {
		rpc.ftm.Close()
	}

	rpc.closed = true
}
