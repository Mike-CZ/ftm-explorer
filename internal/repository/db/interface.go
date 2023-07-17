package db

import (
	"ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
)

// IDatabase is the interface for database operations.
type IDatabase interface {
	// AddBlock adds a block to the database.
	AddBlock(block *types.Block) error

	// GetBlock returns a block from the database.
	GetBlock(number int64) (*db_types.Block, error)

	// Close terminates the database connection.
	Close()
}
