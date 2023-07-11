package rpc

import "ftm-explorer/internal/types"

// Rpc represents the interface for the RPC client.
type Rpc interface {
	// BlockByNumber returns the block identified by number.
	BlockByNumber(number uint64) (*types.Block, error)
}
