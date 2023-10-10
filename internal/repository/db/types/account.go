package db_types

import "github.com/ethereum/go-ethereum/common"

type Account struct {
	Address  common.Address `bson:"_id"`
	LastSeen int64          `bson:"lastSeen"`
}
