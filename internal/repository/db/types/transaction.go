package db_types

import "github.com/ethereum/go-ethereum/common"

// Transaction represents a transaction in the database.
// We only need a few data, so we only define those fields.
type Transaction struct {
	Addresses []common.Address `bson:"addresses"`
	Hash      common.Hash      `bson:"hash"`
	Timestamp int64            `bson:"timestamp"`
}
