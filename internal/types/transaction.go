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
	BlockHash *common.Hash `json:"blockHash"`

	// BlockNumber represents number of the block where this transaction was in.
	BlockNumber *hexutil.Uint64 `json:"blockNumber"`

	// From represents address of the sender.
	From common.Address `json:"from"`

	// To represents the address of the receiver. Nil when it's a contract creation transaction.
	To *common.Address `json:"to"`

	// ContractAddress represents the address of contract created, if a contract creation transaction, otherwise nil.
	ContractAddress *common.Address `json:"contract"`

	// Nonce represents the number of transactions made by the sender prior to this one.
	Nonce hexutil.Uint64 `json:"nonce"`

	// Gas represents gas provided by the sender.
	Gas hexutil.Uint64 `json:"gas"`

	// GasUsed represents the amount of gas used by this specific transaction alone.
	GasUsed *hexutil.Uint64 `json:"gasUsed"`

	// CumulativeGasUsed represents the total amount of gas used when this transaction was executed in the block.
	CumulativeGasUsed *hexutil.Uint64 `json:"cumulativeGasUsed"`

	// GasPrice represents gas price provided by the sender in Wei.
	GasPrice hexutil.Big `json:"gasPrice"`

	// Value represents value transferred in Wei.
	Value hexutil.Big `json:"value"`

	// Input represents the data send along with the transaction.
	InputData hexutil.Bytes `json:"input"`

	// TrxIndex represents integer of the transaction's index position in the block. nil when it's pending.
	TrxIndex *hexutil.Uint64 `json:"transactionIndex"`

	// Status represents transaction status; value is either 1 (success) or 0 (failure)
	Status *hexutil.Uint64 `json:"status"`

	// Logs represents a list of log records created along with the transaction
	Logs []types.Log `json:"logs"`
}
