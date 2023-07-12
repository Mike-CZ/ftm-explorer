package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	// Hash represents 32 bytes hash of the transaction.
	Hash common.Hash `json:"hash"`
	// BlockHash represents hash of the block where this transaction was in.
	BlockHash common.Hash `json:"blockHash"`
	// BlockNumber represents number of the block where this transaction was in.
	BlockNumber hexutil.Uint64 `json:"blockNumber"`
	// From represents address of the sender.
	From common.Address `json:"from"`
	// To represents the address of the receiver. Nil when it's a contract creation transaction.
	To *common.Address `json:"to"`
	// GasUsed represents the amount of gas used by this specific transaction alone.
	GasUsed hexutil.Uint64 `json:"gasUsed"`
	// GasPrice represents gas price provided by the sender in Wei.
	GasPrice hexutil.Big `json:"gasPrice"`
	// Logs represents a list of log records created along with the transaction
	Logs []types.Log `json:"logs"`
}
