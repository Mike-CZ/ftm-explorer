package repository

import "context"

// GetNumberOfValidators returns the number of validators.
func (r *Repository) GetNumberOfValidators() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()
	return r.rpc.NumberOfValidators(ctx)
}
