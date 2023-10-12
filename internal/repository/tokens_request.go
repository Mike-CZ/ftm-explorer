package repository

import (
	"context"
	"ftm-explorer/internal/types"
)

// AddTokensRequest adds a new tokens request to the database.
func (r *Repository) AddTokensRequest(tr *types.TokensRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.AddTokensRequest(ctx, tr)
}

// UpdateTokensRequest updates the given tokens request.
func (r *Repository) UpdateTokensRequest(tr *types.TokensRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.UpdateTokensRequest(ctx, tr)
}

// GetLatestUnclaimedTokensRequest returns the latest unclaimed  tokens request for the given ip address.
func (r *Repository) GetLatestUnclaimedTokensRequest(ip string) (*types.TokensRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.LatestUnclaimedTokensRequest(ctx, ip)
}

// GetLatestClaimedTokensRequests returns the latest claimed tokens requests for the given ip address.
func (r *Repository) GetLatestClaimedTokensRequests(ip string, from uint64) ([]types.TokensRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.LatestClaimedTokensRequests(ctx, ip, from)
}
