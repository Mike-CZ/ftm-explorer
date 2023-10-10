package rpc

//go:generate mockgen -source=rpc.go -destination=rpc_mock.go -package=rpc -mock_names=IRpc=MockRpc

import (
	"context"
	"ftm-explorer/internal/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
)

// IRpc represents the interface for the RPC client.
type IRpc interface {
	// BlockByNumber returns the block identified by number.
	BlockByNumber(context.Context, uint64) (*types.Block, error)
	// TransactionByHash returns the transaction identified by hash.
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, error)
	// ObservedHeadProxy provides a channel fed with new headers.
	ObservedHeadProxy() <-chan *eth.Header
	// NumberOfValidators returns the number of validators.
	NumberOfValidators(context.Context) (uint64, error)
	// SendSignedTransaction sends the signed transaction.
	SendSignedTransaction(context.Context, *eth.Transaction) error
	// PendingNonceAt returns the nonce of the account at the given block.
	PendingNonceAt(context.Context, common.Address) (uint64, error)
	// SuggestGasPrice suggests a gas price.
	SuggestGasPrice(context.Context) (*big.Int, error)
	// NetworkID returns the network ID.
	NetworkID(context.Context) (*big.Int, error)
	// AccountBalance returns the balance of the account.
	AccountBalance(context.Context, common.Address) (*hexutil.Big, error)
	// Close closes the RPC client.
	Close()
}
