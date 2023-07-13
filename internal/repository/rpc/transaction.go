package rpc

import (
	"context"
	"fmt"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
)

// TransactionByHash returns the transaction identified by hash.
func (rpc *OperaRpc) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	var trx types.Transaction

	// get the block by number
	err := rpc.ftm.CallContext(ctx, &trx, "ftm_getTransactionByHash", hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %v", err)
	}

	// get transaction receipt if transaction is not pending
	if trx.BlockNumber != nil {
		var rec struct {
			CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
			GasUsed           hexutil.Uint64  `json:"gasUsed"`
			ContractAddress   *common.Address `json:"contractAddress,omitempty"`
			Status            hexutil.Uint64  `json:"status"`
			Logs              []eth.Log       `json:"logs"`
		}

		// call for the transaction receipt data
		err := rpc.ftm.Call(&rec, "ftm_getTransactionReceipt", hash)
		if err != nil {
			return nil, err
		}

		// set data
		trx.CumulativeGasUsed = &rec.CumulativeGasUsed
		trx.GasUsed = &rec.GasUsed
		trx.ContractAddress = rec.ContractAddress
		trx.Status = &rec.Status
		trx.Logs = rec.Logs
	}

	return &trx, nil
}
