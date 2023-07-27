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

// GetTrxCount returns the number of transactions in the blockchain.
func (r *Repository) GetTrxCount() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	count, err := r.db.TrxCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// IncrementTrxCount increments the number of transactions in the blockchain.
func (r *Repository) IncrementTrxCount(incrementBy uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.IncrementTrxCount(ctx, incrementBy)
}
