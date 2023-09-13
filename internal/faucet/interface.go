package faucet

//go:generate mockgen -source=interface.go -destination=faucet_mock.go -package=faucet -mock_names=IFaucet=MockFaucet,IFaucetPhraseGenerator=MockFaucetPhraseGenerator

import (
	"github.com/ethereum/go-ethereum/common"
)

// IFaucet represents a faucet interface. It provides access to the
// faucet functionality. It is used to request and claim tokens.
type IFaucet interface {
	// RequestTokens requests tokens for the given ip address.
	RequestTokens(ip string) (string, error)

	// ClaimTokens claims tokens for the given phrase and receiver address.
	ClaimTokens(ip string, phrase string, receiver common.Address) error
}

// IFaucetPhraseGenerator represents a faucet phrase generator interface.
// It is used to generate a phrase for the faucet.
type IFaucetPhraseGenerator interface {
	// GeneratePhrase generates a phrase for the faucet.
	GeneratePhrase() string
}
