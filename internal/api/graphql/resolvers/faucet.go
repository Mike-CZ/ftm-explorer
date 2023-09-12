package resolvers

import (
	"context"
	"fmt"
	"ftm-explorer/internal/auth"
	"ftm-explorer/internal/types"
)

// RequestTokens initiates the request for tokens from faucet.
func (rs *RootResolver) RequestTokens(ctx context.Context) (string, error) {
	// get the ip address from the context
	ip, err := auth.GetIpOrErr(ctx)
	if err != nil {
		return "", err
	}

	// run the faucet request as singleflight group to prevent multiple requests from the same ip
	key := fmt.Sprintf("request_tokens_%s", ip)
	tr, err, _ := rs.sfg.Do(key, func() (interface{}, error) {
		return rs.faucet.RequestTokens(ip)
	})
	if err != nil {
		return "", err
	}

	return tr.(*types.TokensRequest).Phrase, nil
}
