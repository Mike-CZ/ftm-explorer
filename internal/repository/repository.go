package repository

import (
	"ftm-explorer/internal/buffer"
	"ftm-explorer/internal/repository/db"
	"ftm-explorer/internal/repository/rpc"
	"time"
)

// kBlkBufferSize represents the size of the block buffer.
const kBlkBufferSize = 10_000

// kRpcTimeout represents the timeout for RPC calls.
const kRpcTimeout = 5 * time.Second

// kDbTimeout represents the timeout for DB calls.
const kDbTimeout = 5 * time.Second

// Repository represents the repository.
// It contains the RPC client and a buffer for blocks.
// The buffer is used to store the latest observed blocks.
type Repository struct {
	rpc       rpc.IRpc
	db        db.IDatabase
	blkBuffer *buffer.BlocksBuffer
}

// NewRepository creates a new repository.
func NewRepository(rpc rpc.IRpc, db db.IDatabase) *Repository {
	return &Repository{
		rpc:       rpc,
		db:        db,
		blkBuffer: buffer.NewBlocksBuffer(kBlkBufferSize),
	}
}
