package rpc

//go:generate mockgen -source=rpc.go -destination=rpc_mock.go -package=rpc -mock_names=IRpc=MockRpc

import (
	"context"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

// IRpc represents the interface for the RPC client.
type IRpc interface {
	// BlockByNumber returns the block identified by number.
	BlockByNumber(context.Context, uint64) (*types.Block, error)
	// TransactionByHash returns the transaction identified by hash.
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, error)
	// ObservedHeadProxy provides a channel fed with new headers.
	ObservedHeadProxy() <-chan *eth.Header
	// NumberOfValidators returns the number of validators.
	NumberOfValidators(context.Context) (uint64, error)
	// Close closes the RPC client.
	Close()
}
