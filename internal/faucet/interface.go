package faucet

//go:generate mockgen -source=interface.go -destination=faucet_mock.go -package=faucet -mock_names=IFaucet=MockFaucet,IFaucetPhraseGenerator=MockFaucetPhraseGenerator,IFaucetWallet=MockFaucetWallet

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// IFaucet represents a faucet interface. It provides access to the
// faucet functionality. It is used to request and claim tokens.
type IFaucet interface {
	// RequestTokens requests tokens for the given ip address.
	RequestTokens(string) (string, error)

	// ClaimTokens claims tokens for the given phrase and receiver address.
	ClaimTokens(ip string, phrase string, receiver common.Address, erc20 *common.Address) error
}

// IFaucetPhraseGenerator represents a faucet phrase generator interface.
// It is used to generate a phrase for the faucet.
type IFaucetPhraseGenerator interface {
	// GeneratePhrase generates a phrase for the faucet.
	GeneratePhrase() (string, error)
}

// IFaucetWallet represents a faucet wallet interface.
// It is used to send wei to the given address.
type IFaucetWallet interface {
	// SendWeiToAddress sends wei to the given address.
	SendWeiToAddress(amount *big.Int, receiver common.Address) error

	// MintErc20TokensToAddress sends erc20 tokens to the given address.
	MintErc20TokensToAddress(common.Address, common.Address, *big.Int) error
}
