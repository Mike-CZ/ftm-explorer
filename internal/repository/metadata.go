package repository

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// GetNumberOfAccounts returns the number of accounts in the blockchain.
func (r *Repository) GetNumberOfAccounts() uint64 {
	return r.numberOfAccounts
}

// SetNumberOfAccounts sets the number of accounts in the blockchain.
func (r *Repository) SetNumberOfAccounts(number uint64) {
	r.numberOfAccounts = number
}

// GetDiskSizePer100MTxs returns the disk size per 100M transactions.
func (r *Repository) GetDiskSizePer100MTxs() uint64 {
	return r.diskSizePer100MTxs
}

// SetDiskSizePer100MTxs sets the disk size per 100M transactions.
func (r *Repository) SetDiskSizePer100MTxs(number uint64) {
	r.diskSizePer100MTxs = number
}

// PendingNonceAt returns the nonce of the account at the given block.
func (r *Repository) PendingNonceAt(address common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.PendingNonceAt(ctx, address)
}

// SuggestGasPrice suggests a gas price.
func (r *Repository) SuggestGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.SuggestGasPrice(ctx)
}

// NetworkID returns the network ID.
func (r *Repository) NetworkID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.NetworkID(ctx)
}
