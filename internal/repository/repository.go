package repository

import (
	"ftm-explorer/internal/buffer"
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/types"
	"time"
)

// kBlkBufferSize represents the size of the block buffer.
const kBlkBufferSize = 100

// kRpcTimeout represents the timeout for RPC calls.
const kRpcTimeout = 5 * time.Second

// Repository represents the repository.
// It contains the RPC client and a buffer for blocks.
// The buffer is used to store the latest observed blocks.
type Repository struct {
	rpc       rpc.IRpc
	blkBuffer *buffer.RingBuffer[uint64, *types.Block]
}

// NewRepository creates a new repository.
func NewRepository(rpc rpc.IRpc) *Repository {
	return &Repository{
		rpc:       rpc,
		blkBuffer: buffer.NewRingBuffer[uint64, *types.Block](kBlkBufferSize),
	}
}
