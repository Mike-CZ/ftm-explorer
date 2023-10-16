package faucet

import (
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// kFaucetChallengePrefix is template of message to be signed using Metamask.
	kFaucetChallengePrefix = "Sign following challenge to obtain tokens:\n"
	// kFaucetChallengePrefixLen is the length of the challenge prefix.
	kFaucetChallengePrefixLen = len(kFaucetChallengePrefix)
)

// FaucetErc20 represents a faucet erc20 token.
type FaucetErc20 struct {
	address common.Address
	wallet  IFaucetWallet
}

// Faucet represents a faucet instance. It provides access to the
// faucet functionality. It is used to request and claim tokens.
type Faucet struct {
	pg              IFaucetPhraseGenerator
	wallet          IFaucetWallet
	repo            repository.IRepository
	cfg             *config.Faucet
	erc20s          map[common.Address]FaucetErc20
	erc20MintAmount *big.Int
}

// NewFaucet creates a new faucet instance.
func NewFaucet(cfg *config.Faucet, pg IFaucetPhraseGenerator, w IFaucetWallet, erc20s []FaucetErc20, repo repository.IRepository) (*Faucet, error) {
	mintAmount, err := hexutil.DecodeBig(cfg.Erc20MintAmountHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding faucet erc20 amount: %v", err)
	}
	f := &Faucet{
		pg:              pg,
		wallet:          w,
		repo:            repo,
		cfg:             cfg,
		erc20s:          make(map[common.Address]FaucetErc20),
		erc20MintAmount: mintAmount,
	}
	for _, erc20 := range erc20s {
		f.erc20s[erc20.address] = erc20
	}
	return f, nil
}

// NewFaucetErc20s creates a new faucet erc20 tokens.
func NewFaucetErc20s(cfg *config.Faucet, repo repository.IRepository, log logger.ILogger) ([]FaucetErc20, error) {
	erc20s := make([]FaucetErc20, len(cfg.Erc20s))
	for i, erc20 := range cfg.Erc20s {
		address := common.HexToAddress(erc20.Address)
		wallet, err := NewWallet(repo, log, erc20.MinterPk)
		if err != nil {
			return nil, fmt.Errorf("error creating faucet wallet: %v", err)
		}
		erc20s[i] = FaucetErc20{
			address: address,
			wallet:  wallet,
		}
	}
	return erc20s, nil
}

// RequestTokens requests tokens for the given ip address and phrase. Returns
// the challenge to be signed by the user.
func (f *Faucet) RequestTokens(ipAddress string) (string, error) {
	// check if the ip address is already in the database
	tr, err := f.repo.GetLatestUnclaimedTokensRequest(ipAddress)
	if err != nil {
		return "", fmt.Errorf("error getting latest tokens request: %v", err)
	}

	// if there is unclaimed request, return it
	if tr != nil {
		return kFaucetChallengePrefix + tr.Phrase, nil
	}

	// otherwise get all requests for the given ip address in the last 24 hours
	requests, err := f.repo.GetLatestClaimedTokensRequests(ipAddress, uint64(time.Now().Unix()-24*60*60))
	if err != nil {
		return "", fmt.Errorf("error getting latest tokens requests: %v", err)
	}

	// check if the number of requests is greater than the allowed number of claims per day
	if len(requests) >= int(f.cfg.ClaimsPerDay) {
		return "", fmt.Errorf("too many requests, you are allowed to claim %d times per day", f.cfg.ClaimsPerDay)
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
	return kFaucetChallengePrefix + tr.Phrase, nil
}

// ClaimTokens claims tokens for the given phrase and receiver address.
func (f *Faucet) ClaimTokens(ip string, phrase string, receiver common.Address, erc20 *common.Address) error {
	// check the phrase
	if len(phrase) < kFaucetChallengePrefixLen || strings.Compare(phrase[:kFaucetChallengePrefixLen], kFaucetChallengePrefix) != 0 {
		return fmt.Errorf("invalid phrase")
	}

	// remove the challenge prefix
	phrase = phrase[kFaucetChallengePrefixLen:]

	// get the latest request for the given ip address
	tr, err := f.repo.GetLatestUnclaimedTokensRequest(ip)
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

	// check if it is known erc20 token when erc20 is not nil
	if erc20 != nil {
		_, ok := f.erc20s[*erc20]
		if !ok {
			return fmt.Errorf("unknown erc20 token")
		}
	}

	// update the request
	tr.Receiver = &receiver
	now := time.Now().Unix()
	tr.ClaimedAt = &now
	err = f.repo.UpdateTokensRequest(tr)
	if err != nil {
		return fmt.Errorf("error updating tokens request: %v", err)
	}

	// send native tokens to the receiver if erc20 is nil
	if erc20 == nil {
		if err = f.wallet.SendWeiToAddress(getTokensAmountInWei(float64(f.cfg.ClaimTokensAmount)), receiver); err != nil {
			// if we got error, we need to set back the request to the previous state
			tr.Receiver = nil
			tr.ClaimedAt = nil
			_ = f.repo.UpdateTokensRequest(tr)
			return fmt.Errorf("error sending wei to address: %v", err)
		}
		return nil
	}

	erc20Faucet := f.erc20s[*erc20]

	// send erc20 tokens to the receiver
	if err = erc20Faucet.wallet.MintErc20TokensToAddress(erc20Faucet.address, receiver, f.erc20MintAmount); err != nil {
		// if we got error, we need to set back the request to the previous state
		tr.Receiver = nil
		tr.ClaimedAt = nil
		_ = f.repo.UpdateTokensRequest(tr)
		return fmt.Errorf("error minting erc20 to address: %v", err)
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
