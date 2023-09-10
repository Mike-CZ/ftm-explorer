package faucet

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

// IFaucet represents a faucet interface. It provides access to the
// faucet functionality. It is used to request and claim tokens.
type IFaucet interface {
	// RequestTokens requests tokens for the given ip address.
	RequestTokens(ip string) (*types.TokensRequest, error)

	// ClaimTokens claims tokens for the given phrase and receiver address.
	ClaimTokens(ip string, phrase string, receiver common.Address) error
}
