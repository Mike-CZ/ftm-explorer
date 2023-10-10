package repository

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// AddAccounts adds accounts to the database.
func (r *Repository) AddAccounts(accs []common.Address, stamp int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.AddAccounts(ctx, accs, stamp)
}

// GetNumberOfAccountsInDb returns the number of accounts in the database.
func (r *Repository) GetNumberOfAccountsInDb() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.NumberOfAccoutns(ctx)
}
