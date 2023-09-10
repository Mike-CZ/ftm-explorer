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

// GetLatestTokensRequest returns the latest tokens request for the given ip address.
func (r *Repository) GetLatestTokensRequest(ip string) (*types.TokensRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.LatestTokensRequest(ctx, ip)
}
