package db

import (
	"context"
	"fmt"
	"ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// kCoBlocks is the name of the block collection.
	kCoBlocks = "block"

	// kFiBlockNumber is the name of the block number field. It is also the primary key.
	kFiBlockNumber = "_id"

	// kFiBlockTxCount is the name of the block transaction count field.
	kFiBlockTxCount = "txCount"

	// kFiBlockGasUsed is the name of the block gas used field.
	kFiBlockGasUsed = "gasUsed"

	// kFiBlockTimestamp is the name of the block timestamp field.
	kFiBlockTimestamp = "timestamp"
)

// AddBlock adds a block to the database.
func (db *MongoDb) AddBlock(block *types.Block) error {
	if block == nil {
		return fmt.Errorf("can not add empty block")
	}

	// try to do the insert
	if _, err := db.blockCollection().InsertOne(context.Background(), &db_types.Block{
		Number:    int64(block.Number),
		TxCount:   int32(len(block.Transactions)),
		GasUsed:   int64(block.GasUsed),
		Timestamp: int64(block.Timestamp),
	}); err != nil {
		db.log.Critical(err)
		return err
	}

	db.log.Debugf("block %d added to database", int64(block.Number))

	return nil
}

// GetBlock returns a block from the database.
func (db *MongoDb) GetBlock(number int64) (*db_types.Block, error) {
	// try to get the block
	var block db_types.Block
	if err := db.blockCollection().FindOne(context.Background(), bson.D{{Key: kFiBlockNumber, Value: number}}).Decode(&block); err != nil {
		db.log.Critical(err)
		return nil, err
	}

	return &block, nil
}

// blockCollection returns the block collection.
func (db *MongoDb) blockCollection() *mongo.Collection {
	return db.db.Collection(kCoBlocks)
}

// initBlockCollection initializes the block collection.
func (db *MongoDb) initBlockCollection() {
	// prepare index models
	ix := make([]mongo.IndexModel, 0)

	// index the timestamp
	ix = append(ix, mongo.IndexModel{Keys: bson.D{{Key: kFiBlockTimestamp, Value: 1}}})

	// create indexes
	ctx, cancel := context.WithTimeout(context.Background(), kMongoDefaultTimeout)
	defer cancel()
	if _, err := db.blockCollection().Indexes().CreateMany(ctx, ix); err != nil {
		db.log.Panicf("can not create indexes for block collection; %V", err)
	}

	db.log.Debugf("transactions collection initialized")
}
