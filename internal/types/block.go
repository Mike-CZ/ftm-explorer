package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Block represents basic information provided by the API about block inside Opera blockchain.
type Block struct {
	// Number represents the block number.
	Number hexutil.Uint64 `json:"number"`

	// Epoch represents the block epoch.
	Epoch hexutil.Uint64 `json:"epoch"`

	// Hash represents hash of the block.
	Hash common.Hash `json:"hash"`

	// ParentHash represents hash of the parent block.
	ParentHash common.Hash `json:"parentHash"`

	// GasUsed represents the actual total used gas by all transactions in this block.
	GasUsed hexutil.Uint64 `json:"gasUsed"`

	// GasLimit represents the maximum gas allowed in this block.
	GasLimit hexutil.Uint64 `json:"gasLimit"`

	// Timestamp represents the unix timestamp for when the block was collated.
	Timestamp hexutil.Uint64 `json:"timestamp"`

	// Transactions represents array of 32 bytes hashes of transactions included in the block.
	Transactions []common.Hash `json:"transactions"`
}
