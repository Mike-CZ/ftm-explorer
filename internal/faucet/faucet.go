package faucet

import (
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// challengePrefix is template of message to be signed using Metamask.
	challengePrefix = "Sign following challenge to obtain tokens:\n"
	// challengePrefixLen is the length of the challenge prefix.
	challengePrefixLen = len(challengePrefix)
)

// Faucet represents a faucet instance. It provides access to the
// faucet functionality. It is used to request and claim tokens.
type Faucet struct {
	pg     IFaucetPhraseGenerator
	wallet IFaucetWallet
	repo   repository.IRepository
	cfg    *config.Faucet
}

// NewFaucet creates a new faucet instance.
func NewFaucet(cfg *config.Faucet, pg IFaucetPhraseGenerator, w IFaucetWallet, repo repository.IRepository) *Faucet {
	return &Faucet{
		pg:     pg,
		wallet: w,
		repo:   repo,
		cfg:    cfg,
	}
}

// RequestTokens requests tokens for the given ip address and phrase. Returns
// the challenge to be signed by the user.
func (f *Faucet) RequestTokens(ipAddress string) (string, error) {
	// check if the ip address is already in the database
	tr, err := f.repo.GetLatestTokensRequest(ipAddress)
	if err != nil {
		return "", fmt.Errorf("error getting latest tokens request: %v", err)
	}

	// if claim already exists, check its status
	if tr != nil {
		// check if the request was not already fulfilled
		// if it was not, return the request
		if tr.ClaimedAt == nil {
			return challengePrefix + tr.Phrase, nil
		}
		// otherwise check if the claim limit was reached
		diff := uint(time.Now().Sub(time.Unix(*tr.ClaimedAt, 0)).Seconds())
		if diff < f.cfg.ClaimLimitSeconds {
			return "", fmt.Errorf("must wait %d seconds before claiming again", f.cfg.ClaimLimitSeconds-diff)
		}
	}

	phrase, err := f.pg.GeneratePhrase()
	if err != nil {
		return "", fmt.Errorf("error generating phrase: %v", err)
	}

	// create a new request
	tr = &types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
	}
	err = f.repo.AddTokensRequest(tr)
	if err != nil {
		return "", err
	}
	return challengePrefix + tr.Phrase, nil
}

// ClaimTokens claims tokens for the given phrase and receiver address.
func (f *Faucet) ClaimTokens(ip string, phrase string, receiver common.Address) error {
	// check the phrase
	if len(phrase) < challengePrefixLen || strings.Compare(phrase[:challengePrefixLen], challengePrefix) != 0 {
		return fmt.Errorf("invalid phrase")
	}

	// remove the challenge prefix
	phrase = phrase[challengePrefixLen:]

	// get the latest request for the given ip address
	tr, err := f.repo.GetLatestTokensRequest(ip)
	if err != nil {
		return err
	}
	if tr == nil {
		return fmt.Errorf("no request found for ip address %s", ip)
	}

	// check the phrase matches
	if tr.Phrase != phrase {
		return fmt.Errorf("invalid phrase")
	}

	// check if the request was already fulfilled
	if tr.ClaimedAt != nil {
		return fmt.Errorf("tokens already claimed")
	}

	// update the request
	tr.Receiver = &receiver
	now := time.Now().Unix()
	tr.ClaimedAt = &now
	err = f.repo.UpdateTokensRequest(tr)
	if err != nil {
		return fmt.Errorf("error updating tokens request: %v", err)
	}

	// send wei to the receiver
	if err = f.wallet.SendWeiToAddress(getTokensAmountInWei(float64(f.cfg.ClaimTokensAmount)), receiver); err != nil {
		// if we got error, we need to set back the request to the previous state
		tr.Receiver = nil
		tr.ClaimedAt = nil
		_ = f.repo.UpdateTokensRequest(tr)
		return fmt.Errorf("error sending wei to address: %v", err)
	}

	return nil
}

// getTokensAmountInWei converts the given amount of tokens to wei.
func getTokensAmountInWei(amount float64) *big.Int {
	a := new(big.Rat).SetFloat64(amount)
	e := new(big.Rat).SetFloat64(params.Ether)
	product := new(big.Rat).Mul(a, e)
	result := new(big.Int).Div(product.Num(), product.Denom())
	return result
}
