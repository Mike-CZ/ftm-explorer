package resolvers

import (
	"ftm-explorer/internal/faucet"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/maze"
	"ftm-explorer/internal/repository"

	"golang.org/x/sync/singleflight"
)

// RootResolver is GraphQL resolver of root namespace.
type RootResolver struct {
	repository repository.IRepository
	log        logger.ILogger
	faucet     faucet.IFaucet
	maze       maze.IMaze

	// singleflight is used to prevent multiple concurrent requests for the same data.
	sfg singleflight.Group

	// isPersisted is true if the explorer is running in persisted mode.
	// this might get removed in the future.
	isPersisted bool
}

// NewResolver creates a new root resolver.
func NewResolver(repository repository.IRepository, log logger.ILogger, faucet faucet.IFaucet, maze maze.IMaze, isPersisted bool) *RootResolver {
	return &RootResolver{
		repository:  repository,
		log:         log.ModuleLogger("resolver"),
		faucet:      faucet,
		maze:        maze,
		isPersisted: isPersisted,
	}
}
