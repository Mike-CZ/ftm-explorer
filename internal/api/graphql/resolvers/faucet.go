package resolvers

import (
	"context"
	"fmt"
	"ftm-explorer/internal/auth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// RequestTokens initiates the request for tokens from faucet.
func (rs *RootResolver) RequestTokens(ctx context.Context, args struct {
	Symbol *string
}) (string, error) {
	// get the ip address from the context
	ip, err := auth.GetIpOrErr(ctx)
	if err != nil {
		return "", err
	}

	// run the faucet request as singleflight group to prevent multiple requests from the same ip
	key := fmt.Sprintf("request_tokens_%s", ip)
	phrase, err, _ := rs.sfg.Do(key, func() (interface{}, error) {
		return rs.faucet.RequestTokens(ip, args.Symbol)
	})
	if err != nil {
		return "", err
	}

	return phrase.(string), nil
}

// ClaimTokens claims the tokens from faucet. It requires the user to sign the challenge.
func (rs *RootResolver) ClaimTokens(ctx context.Context, args struct {
	Address      common.Address
	Challenge    string
	Signature    string
	Erc20Address *common.Address
}) (bool, error) {
	// get the ip address from the context
	ip, err := auth.GetIpOrErr(ctx)
	if err != nil {
		return false, err
	}

	// run the faucet request as singleflight group to prevent multiple claims from the same ip
	key := fmt.Sprintf("claim_tokens_%s", ip)
	_, err, _ = rs.sfg.Do(key, func() (interface{}, error) {
		// decode signature
		decodedSignature, err := hexutil.Decode(args.Signature)
		if err != nil {
			return "", fmt.Errorf("signature hex decoding failed; %s", err)
		}
		// verify signature
		_, err = auth.VerifySignature(args.Challenge, args.Address, decodedSignature)
		if err != nil {
			return "", fmt.Errorf("signature verification failed; %s", err)
		}
		// claim tokens
		err = rs.faucet.ClaimTokens(ip, args.Challenge, args.Address, args.Erc20Address)
		return "", err
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
