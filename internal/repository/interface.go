package repository

//go:generate mockgen -source=interface.go -destination=repository_mock.go -package=repository -mock_names=IRepository=MockRepository

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
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
	UpdateLatestObservedBlock(*types.Block) error

	// GetNewHeadersChannel returns a channel that will receive the latest headers from blockchain.
	GetNewHeadersChannel() <-chan *eth.Header

	// GetTransactionByHash returns the transaction identified by hash.
	GetTransactionByHash(common.Hash) (*types.Transaction, error)

	// GetNumberOfValidators returns the number of validators.
	GetNumberOfValidators() (uint64, error)

	// GetTrxCountAggByTimestamp returns aggregation of transactions in given time range.
	GetTrxCountAggByTimestamp(types.AggResolution, uint, *uint64) ([]types.HexUintTick, error)

	// GetGasUsedAggByTimestamp returns aggregation of gas used in given time range.
	GetGasUsedAggByTimestamp(types.AggResolution, uint, *uint64) ([]types.HexUintTick, error)

	// GetNumberOfAccounts returns the number of accounts in the blockchain.
	GetNumberOfAccounts() uint64

	// SetNumberOfAccounts sets the number of accounts in the blockchain.
	SetNumberOfAccounts(uint64)

	// GetDiskSizePer100MTxs returns the disk size per 100M transactions.
	GetDiskSizePer100MTxs() uint64

	// SetDiskSizePer100MTxs sets the disk size per 100M transactions.
	SetDiskSizePer100MTxs(uint64)

	// GetTrxCount returns the number of transactions in the blockchain.
	GetTrxCount() (uint64, error)

	// IncrementTrxCount increments the number of transactions in the blockchain.
	IncrementTrxCount(uint) error

	// FetchNumberOfAccounts returns the number of accounts in the blockchain.
	// This method will fetch data from remote host.
	FetchNumberOfAccounts() (uint64, error)

	// FetchDiskSizePer100MTxs returns the disk size per 100M transactions.
	// This method will fetch data from remote host.
	FetchDiskSizePer100MTxs() (uint64, error)

	// FetchTimeToFinality returns the time to finality in the blockchain.
	// This method will fetch data from remote host.
	FetchTimeToFinality() (float64, error)
}
