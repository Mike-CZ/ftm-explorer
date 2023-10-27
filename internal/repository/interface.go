package repository

//go:generate mockgen -source=interface.go -destination=repository_mock.go -package=repository -mock_names=IRepository=MockRepository

import (
	db_types "ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
	"math/big"

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

	// GetDiskSizePrunedPer100MTxs returns the disk size pruned per 100M transactions.
	GetDiskSizePrunedPer100MTxs() uint64

	// SetDiskSizePrunedPer100MTxs sets the disk size pruned per 100M transactions.
	SetDiskSizePrunedPer100MTxs(uint64)

	// GetTxCountPer10Secs returns transactions per 10 seconds.
	GetTxCountPer10Secs() []types.HexUintTick

	// SetTxCountPer10Secs sets transactions per 10 seconds.
	SetTxCountPer10Secs([]types.HexUintTick)

	// GetGasUsedPer10Secs returns gas used per 10 seconds.
	GetGasUsedPer10Secs() []types.HexUintTick

	// SetGasUsedPer10Secs sets gas used per 10 seconds.
	SetGasUsedPer10Secs([]types.HexUintTick)

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

	// FetchDiskSizePrunedPer100MTxs returns the disk size pruned per 100M transactions.
	// This method will fetch data from remote host.
	FetchDiskSizePrunedPer100MTxs() (uint64, error)

	// FetchTimeToFinality returns the time to finality in the blockchain.
	// This method will fetch data from remote host.
	FetchTimeToFinality() (float64, error)

	// GetTimeToFinality returns the time to finality in the blockchain.
	GetTimeToFinality() float64

	// AddTokensRequest adds a new tokens request to the database.
	AddTokensRequest(*types.TokensRequest) error

	// UpdateTokensRequest updates the given tokens request.
	UpdateTokensRequest(*types.TokensRequest) error

	// GetLatestUnclaimedTokensRequest returns the latest tokens request for the given ip address.
	GetLatestUnclaimedTokensRequest(string) (*types.TokensRequest, error)

	// GetLatestClaimedTokensRequests returns the latest claimed tokens requests for the given ip address.
	GetLatestClaimedTokensRequests(string, uint64) ([]types.TokensRequest, error)

	// SendSignedTransaction sends the signed transaction.
	SendSignedTransaction(*eth.Transaction) error

	// PendingNonceAt returns the nonce of the account at the given block.
	PendingNonceAt(common.Address) (uint64, error)

	// SuggestGasPrice suggests a gas price.
	SuggestGasPrice() (*big.Int, error)

	// NetworkID returns the network ID.
	NetworkID() (*big.Int, error)

	// GetTimeToBlock returns the time to block.
	GetTimeToBlock() float64

	// AddTimeToFinality adds the given time to finality.
	AddTimeToFinality(*types.Ttf) error

	// GetTtfAvgAggByTimestamp returns average aggregation of time to finality in given time range.
	GetTtfAvgAggByTimestamp(types.AggResolution, uint, uint64) ([]types.FloatTick, error)

	// GetTimeToFinalityPer10Secs returns time to finality per 10 seconds.
	GetTimeToFinalityPer10Secs() []types.FloatTick

	// SetTimeToFinalityPer10Secs sets time to finality per 10 seconds.
	SetTimeToFinalityPer10Secs([]types.FloatTick)

	// AddTransactions adds transactions to the database.
	AddTransactions([]db_types.Transaction) error

	// GetLastTransactionsWhereAddress returns the last transactions where the given address is involved.
	GetLastTransactionsWhereAddress(common.Address, uint) ([]db_types.Transaction, error)

	// IsIdle returns isIdle.
	IsIdle() bool

	// SetIsIdle sets isIdle.
	SetIsIdle(isIdle bool)

	// ShrinkTransactions shrinks the transactions collection. It will persist the given number of transactions.
	// It will delete the oldest transactions.
	ShrinkTransactions(int64) error

	// ShrinkTtf shrinks the time to finality collection. It will persist the given number of ttfs.
	ShrinkTtf(int64) error

	// AddAccounts adds accounts to the database.
	AddAccounts(accs []common.Address, stamp int64) error

	// GetNumberOfAccountsInDb returns the number of accounts in the database.
	GetNumberOfAccountsInDb() (uint64, error)

	// AccountBalance returns the balance of the account.
	AccountBalance(common.Address) (*hexutil.Big, error)
}
