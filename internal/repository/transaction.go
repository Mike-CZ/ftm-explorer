package repository

import (
	"context"
	db_types "ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
	"time"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
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

// SendSignedTransaction sends the signed transaction.
func (r *Repository) SendSignedTransaction(tx *eth.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.SendSignedTransaction(ctx, tx)
}

// AddTransactions adds transactions to the database.
func (r *Repository) AddTransactions(txs []db_types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.AddTransactions(ctx, txs)
}

// GetLastTransactionsWhereAddress returns the last transactions for the given address.
func (r *Repository) GetLastTransactionsWhereAddress(addr common.Address, count uint) ([]db_types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.LastTransactionsWhereAddress(ctx, addr, count)
}

// ShrinkTransactions shrinks the transactions collection. It will persist the given number of transactions.
// It will delete the oldest transactions.
func (r *Repository) ShrinkTransactions(count int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return r.db.ShrinkTransactions(ctx, count)
}
