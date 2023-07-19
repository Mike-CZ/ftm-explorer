package db

import (
	"context"
	"ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// IDatabase is the interface for database operations.
type IDatabase interface {
	// TrxCountAggByTimestamp returns aggregation of transactions in given time range.
	TrxCountAggByTimestamp(context.Context, uint64, uint, uint) ([]types.Tick[hexutil.Uint64], error)

	// GasUsedAggByTimestamp returns aggregation of gas used in given time range.
	GasUsedAggByTimestamp(context.Context, uint64, uint, uint) ([]types.Tick[hexutil.Uint64], error)

	// AddBlock adds a block to the database.
	AddBlock(context.Context, *types.Block) error

	// Block returns a block from the database.
	Block(context.Context, uint64) (*db_types.Block, error)

	// Close terminates the database connection.
	Close()
}
