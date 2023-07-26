package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
)

// iService represents a service run by the Manager.
type iService interface {
	// start executes the service
	start()
	// close terminates the service
	close()
	// name provides a name of the service
	name() string
}

// service implements general base for services implementing svc interface.
type service struct {
	repo repository.IRepository
	log  logger.ILogger
	mgr  *Manager
}
