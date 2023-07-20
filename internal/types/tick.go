package types

import "github.com/ethereum/go-ethereum/common/hexutil"

// Tick is a timestamp/value pair.
type Tick[T any] struct {
	Time  uint64
	Value T
}

// HexUintTick is a timestamp/hexutil.Uint64 pair.
type HexUintTick = Tick[hexutil.Uint64]
