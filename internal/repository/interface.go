package repository

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
)

type IRepository interface {
	// GetBlockByNumber returns the block identified by number.
	GetBlockByNumber(uint64) (*types.Block, error)

	// GetLatestObservedBlocks returns the number of latest observed blocks.
	GetLatestObservedBlocks(uint) []*types.Block

	// GetLatestObservedBlock returns the latest observed block.
	GetLatestObservedBlock() *types.Block

	// UpdateLatestObservedBlock updates the latest observed block.
	UpdateLatestObservedBlock(*types.Block)

	// GetNewHeadersChannel returns a channel that will receive the latest headers from blockchain.
	GetNewHeadersChannel() <-chan *eth.Header

	// GetTransactionByHash returns the transaction identified by hash.
	GetTransactionByHash(common.Hash) (*types.Transaction, error)

	// GetTrxCountAggByTimestamp returns aggregation of transactions in given time range.
	GetTrxCountAggByTimestamp(types.AggResolution, uint, *uint64) ([]types.Tick[hexutil.Uint64], error)

	// GetGasUsedAggByTimestamp returns aggregation of gas used in given time range.
	GetGasUsedAggByTimestamp(types.AggResolution, uint, *uint64) ([]types.Tick[hexutil.Uint64], error)
}
