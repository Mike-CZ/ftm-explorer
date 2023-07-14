package resolvers

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

// Transaction resolves blockchain transaction by transaction hash.
func (rs *RootResolver) Transaction(args *struct{ Hash common.Hash }) (*types.Transaction, error) {
	trx, err := rs.repository.GetTransactionByHash(args.Hash)
	if err != nil {
		rs.log.Warningf("Failed to get transaction by hash [%s]; %v", args.Hash.Hex(), err)
		return nil, err
	}
	return trx, nil
}
