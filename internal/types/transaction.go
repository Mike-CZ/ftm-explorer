package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	// Hash represents 32 bytes hash of the transaction.
	Hash common.Hash
	// BlockHash represents hash of the block where this transaction was in.
	BlockHash common.Hash
	// BlockNumber represents number of the block where this transaction was in
	BlockNumber uint64
	// Timestamp represents the timestamp of the transaction.
	Timestamp time.Time
	// From represents address of the sender.
	From common.Address
	// To represents the address of the receiver.
	To common.Address
	// GasUsed represents the amount of gas used by this specific transaction alone.
	GasUsed uint64
	// GasPrice represents gas price provided by the sender in Wei.
	GasPrice hexutil.Big
	// RewardToClaim represents the amount of reward to claim in Wei.
	RewardToClaim hexutil.Big
	// Logs represents a list of log records created along with the transaction
	Logs []types.Log
}
