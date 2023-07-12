package repository

import (
	"context"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

// GetTransactionByHash returns the transaction identified by hash.
// It will always fetch the transaction from the RPC.
func (r *Repository) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()
	return r.rpc.TransactionByHash(ctx, hash)
}
