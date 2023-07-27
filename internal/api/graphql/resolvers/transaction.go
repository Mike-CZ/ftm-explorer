package resolvers

import (
	"ftm-explorer/internal/types"
	"ftm-explorer/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Transaction types.Transaction

// Transaction resolves blockchain transaction by transaction hash.
func (rs *RootResolver) Transaction(args *struct{ Hash common.Hash }) (*Transaction, error) {
	trx, err := rs.repository.GetTransactionByHash(args.Hash)
	if err != nil {
		rs.log.Warningf("Failed to get transaction by hash [%s]; %v", args.Hash.Hex(), err)
		return nil, err
	}

	return (*Transaction)(trx), nil
}

// NumberOfTransactions resolves number of transactions in the blockchain.
func (rs *RootResolver) NumberOfTransactions() (hexutil.Uint64, error) {
	count, err := rs.repository.GetTrxCount()
	if err != nil {
		rs.log.Warningf("Failed to get number of transactions; %v", err)
		return 0, err
	}
	return hexutil.Uint64(count), nil
}

// Type resolves transaction type.
func (trx *Transaction) Type() string {
	return utils.ParseTrxType((*types.Transaction)(trx))
}
