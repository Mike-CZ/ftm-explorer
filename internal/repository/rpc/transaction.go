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

	// get transaction receipt
	var rec struct {
		GasUsed hexutil.Uint64 `json:"gasUsed"`
		Logs    []eth.Log      `json:"logs"`
	}
	err = rpc.ftm.Call(&rec, "ftm_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	}

	// assign receipt data
	trx.GasUsed = rec.GasUsed
	trx.Logs = rec.Logs

	return &trx, nil
}
