package db

//go:generate mockgen -source=interface.go -destination=database_mock.go -package=db -mock_names=IDatabase=MockDatabase

import (
	"context"
	"ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
)

// IDatabase is the interface for database operations.
type IDatabase interface {
	// TrxCountAggByTimestamp returns aggregation of transactions in given time range.
	TrxCountAggByTimestamp(context.Context, uint64, uint, uint) ([]types.HexUintTick, error)

	// GasUsedAggByTimestamp returns aggregation of gas used in given time range.
	GasUsedAggByTimestamp(context.Context, uint64, uint, uint) ([]types.HexUintTick, error)

	// AddBlock adds a block to the database.
	AddBlock(context.Context, *types.Block) error

	// Block returns a block from the database.
	Block(context.Context, uint64) (*db_types.Block, error)

	// TrxCount returns the number of transactions in the blockchain.
	TrxCount(context.Context) (uint64, error)

	// IncrementTrxCount increments the number of transactions in the blockchain.
	IncrementTrxCount(context.Context, uint) error

	// AddTokensRequest adds a new tokens request to the database.
	AddTokensRequest(context.Context, *types.TokensRequest) error

	// UpdateTokensRequest updates the given tokens request.
	UpdateTokensRequest(context.Context, *types.TokensRequest) error

	// LatestTokensRequest returns the latest tokens request for the given ip address.
	LatestTokensRequest(context.Context, string) (*types.TokensRequest, error)

	// AddTimeToFinality adds the given time to finality.
	AddTimeToFinality(context.Context, *types.Ttf) error

	// TtfAvgAggByTimestamp returns average aggregation of time to finality in given time range.
	TtfAvgAggByTimestamp(context.Context, uint64, uint, uint) ([]types.FloatTick, error)

	// Close terminates the database connection.
	Close()
}
