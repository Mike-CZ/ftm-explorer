package resolvers

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
)

// RootResolver is GraphQL resolver of root namespace.
type RootResolver struct {
	repository *repository.Repository
	log        logger.Logger
}

// NewResolver creates a new root resolver.
func NewResolver(repository *repository.Repository, log logger.Logger) *RootResolver {
	return &RootResolver{
		repository: repository,
		log:        log,
	}
}
