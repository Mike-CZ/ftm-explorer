package resolvers

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
)

// RootResolver is GraphQL resolver of root namespace.
type RootResolver struct {
	repository repository.IRepository
	log        logger.ILogger
}

// NewResolver creates a new root resolver.
func NewResolver(repository repository.IRepository, log logger.ILogger) *RootResolver {
	return &RootResolver{
		repository: repository,
		log:        log.ModuleLogger("resolver"),
	}
}
