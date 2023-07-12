package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Block represents basic information provided by the API about block inside Opera blockchain.
type Block struct {
	// Number represents the block number.
	Number uint64
	// Epoch represents the block epoch.
	Epoch uint64
	// Hash represents hash of the block. nil when its pending block.
	Hash common.Hash
	// GasUsed represents the actual total used gas by all transactions in this block.
	GasUsed uint64
	// TimeStamp represents the timestamp for when the block was collated.
	TimeStamp time.Time
	// Txs represents array of 32 bytes hashes of transactions included in the block.
	Txs []common.Hash
}
