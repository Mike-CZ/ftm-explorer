package resolvers

import (
	"ftm-explorer/internal/faucet"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"

	"golang.org/x/sync/singleflight"
)

// RootResolver is GraphQL resolver of root namespace.
type RootResolver struct {
	repository repository.IRepository
	log        logger.ILogger
	faucet     faucet.IFaucet

	// singleflight is used to prevent multiple concurrent requests for the same data.
	sfg singleflight.Group
}

// NewResolver creates a new root resolver.
func NewResolver(repository repository.IRepository, log logger.ILogger, faucet faucet.IFaucet) *RootResolver {
	return &RootResolver{
		repository: repository,
		log:        log.ModuleLogger("resolver"),
		faucet:     faucet,
	}
}
