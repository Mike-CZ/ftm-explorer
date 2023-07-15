package rpc

import (
	"ftm-explorer/internal/config"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	client "github.com/ethereum/go-ethereum/rpc"
)

// OperaRpc is a rpc client for Fantom Opera
type OperaRpc struct {
	ftm *client.Client
	// received blocks proxy
	wg       *sync.WaitGroup
	sigClose chan struct{}
	headers  chan *types.Header
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
		ftm:    ftm,
		wg:     new(sync.WaitGroup),
		closed: false,
	}, nil
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
