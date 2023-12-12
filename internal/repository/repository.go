package repository

import (
	"ftm-explorer/internal/buffer"
	"ftm-explorer/internal/repository/db"
	"ftm-explorer/internal/repository/meta_fetcher"
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/types"
	"time"
)

// kRpcTimeout represents the timeout for RPC calls.
const kRpcTimeout = 5 * time.Second

// kDbTimeout represents the timeout for DB calls.
const kDbTimeout = 5 * time.Second

// Repository represents the repository.
// It contains the RPC client and a buffer for blocks.
// The buffer is used to store the latest observed blocks.
type Repository struct {
	rpc         rpc.IRpc
	db          db.IDatabase
	metaFetcher meta_fetcher.IMetaFetcher
	blkBuffer   *buffer.BlocksBuffer

	// numberOfAccounts is the number of accounts in the blockchain.
	numberOfAccounts uint64
	// diskSizePer100MTxs is the disk size per 100M transactions.
	diskSizePer100MTxs uint64
	// diskSizePer100MTxs is the disk size pruned per 100M transactions.
	diskSizePrunedPer100MTxs uint64

	// txCountPer10Secs is the number of transactions per 10 seconds.
	txCountPer10Secs []types.HexUintTick
	// gasUsedPer10Secs is the gas used per 10 seconds.
	gasUsedPer10Secs []types.HexUintTick
	// ttfPer10Secs is the time to finality per 10 seconds.
	ttfPer10Secs []types.FloatTick

	// isIdle indicates if the chain is idle.
	isIdle bool

	// isIdleOverride indicates if the chain is idle. This is used to override the isIdle value.
	isIdleOverride bool
}

// NewRepository creates a new repository.
func NewRepository(blkBufferSize uint, rpc rpc.IRpc, db db.IDatabase, mf meta_fetcher.IMetaFetcher) *Repository {
	return &Repository{
		rpc:              rpc,
		db:               db,
		metaFetcher:      mf,
		blkBuffer:        buffer.NewBlocksBuffer(blkBufferSize),
		numberOfAccounts: 0,
	}
}
