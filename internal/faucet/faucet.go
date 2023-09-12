package faucet

import (
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Faucet represents a faucet instance. It provides access to the
// faucet functionality. It is used to request and claim tokens.
type Faucet struct {
	repo repository.IRepository
	cfg  *config.Faucet
}

// NewFaucet creates a new faucet instance.
func NewFaucet(repo repository.IRepository, cfg *config.Faucet) *Faucet {
	return &Faucet{
		repo: repo,
		cfg:  cfg,
	}
}

// RequestTokens requests tokens for the given ip address and phrase.
func (f Faucet) RequestTokens(ipAddress string) (*types.TokensRequest, error) {
	// check if the ip address is already in the database
	tr, err := f.repo.GetLatestTokensRequest(ipAddress)
	if err != nil {
		return nil, err
	}

	// if claim already exists, check its status
	if tr != nil {
		// check if the request was not already fulfilled
		// if it was not, return the request
		if tr.ClaimedAt == nil {
			return tr, nil
		}
		// otherwise check if the claim limit was reached
		diff := uint(time.Now().Sub(time.Unix(*tr.ClaimedAt, 0)).Seconds())
		if diff < f.cfg.ClaimLimitSeconds {
			return nil, fmt.Errorf("must wait %d seconds before claiming again", f.cfg.ClaimLimitSeconds-diff)
		}
	}

	// create a new request
	tr = &types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "some phrase",
	}
	err = f.repo.AddTokensRequest(tr)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// ClaimTokens claims tokens for the given phrase and receiver address.
func (f Faucet) ClaimTokens(ip string, phrase string, receiver common.Address) error {
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
		return err
	}

	// TODO: send the tokens

	return nil
}
