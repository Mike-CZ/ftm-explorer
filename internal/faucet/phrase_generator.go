package faucet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tyler-smith/go-bip39"
)

// PhraseGenerator is an implementation of the IFaucetPhraseGenerator interface.
// It is used to generate a phrase for the faucet.
type PhraseGenerator struct{}

// NewPhraseGenerator creates a new PhraseGenerator.
func NewPhraseGenerator() *PhraseGenerator {
	return &PhraseGenerator{}
}

// GeneratePhrase generates a phrase for the faucet.
func (f *PhraseGenerator) GeneratePhrase() (string, error) {
	// generate phrase based on bip-39 standard
	entropy, err := bip39.NewEntropy(256) // 256 bits to get a 24-word mnemonic
	if err != nil {
		return "", fmt.Errorf("error generating entropy: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("error generating mnemonic: %v", err)
	}
	// calculate hash from mnemonic
	hash, err := bip39.MnemonicToByteArray(mnemonic)
	if err != nil {
		return "", fmt.Errorf("error calculating hash from mnemonic: %v", err)
	}
	// return hex encoded hash
	return hexutil.Encode(hash), nil
}
