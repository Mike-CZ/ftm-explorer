package repository

import (
	"ftm-explorer/internal/buffer"
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/types"
)

const kBlkBufferSize = 100

type Repository struct {
	rpc       rpc.Rpc
	blkBuffer *buffer.RingBuffer[uint64, *types.Block]
}

func NewRepository(rpc rpc.Rpc) *Repository {
	return &Repository{
		rpc:       rpc,
		blkBuffer: buffer.NewRingBuffer[uint64, *types.Block](kBlkBufferSize),
	}
}
