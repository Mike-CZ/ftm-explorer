package rpc

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

// Rpc represents the interface for the RPC client.
type Rpc interface {
	// BlockByNumber returns the block identified by number.
	BlockByNumber(uint64) (*types.Block, error)
	// TransactionByHash returns the transaction identified by hash.
	TransactionByHash(common.Hash) (*types.Transaction, error)
}
