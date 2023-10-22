package resolvers

import (
	"ftm-explorer/internal/types"
	"ftm-explorer/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Transaction struct {
	types.Transaction
	rs *RootResolver
}

// Transaction resolves blockchain transaction by transaction hash.
func (rs *RootResolver) Transaction(args *struct{ Hash common.Hash }) (*Transaction, error) {
	trx, err := rs.repository.GetTransactionByHash(args.Hash)
	if err != nil {
		rs.log.Warningf("Failed to get transaction by hash [%s]; %v", args.Hash.Hex(), err)
		return nil, err
	}

	return &Transaction{Transaction: *trx, rs: rs}, nil
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
	return utils.ParseTrxType(&trx.Transaction)
}

// Block resolves transaction block.
func (trx *Transaction) Block() (*Block, error) {
	if trx.BlockNumber == nil {
		return nil, nil
	}
	block, err := trx.rs.repository.GetBlockByNumber(uint64(*trx.BlockNumber))
	if err != nil {
		trx.rs.log.Warningf("Failed to get block by hash [%s]; %v", trx.BlockHash.Hex(), err)
		return nil, err
	}
	return &Block{Block: *block, rs: trx.rs}, nil
}
