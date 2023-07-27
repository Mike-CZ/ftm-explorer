package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// kCoState is the name of the state collection.
	kCoState = "state"

	// kFiStatePk is the name of the primary key of the state collection.
	kFiStatePk = "_id"

	// kPkStateTrxCount is the name of the transaction count key.
	kPkStateTrxCount = "trx_count"
)

// TrxCount returns the number of transactions in the blockchain.
func (db *MongoDb) TrxCount(ctx context.Context) (uint64, error) {
	var result struct {
		TrxCount int64 `bson:"trx_count"`
	}

	err := db.stateCollection().FindOne(ctx, bson.M{kFiStatePk: kPkStateTrxCount}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return 0, nil
	} else if err != nil {
		db.log.Errorf("error getting trx count: %v", err)
		return 0, err
	}

	return uint64(result.TrxCount), nil
}

// IncrementTrxCount increments the number of transactions in the blockchain.
func (db *MongoDb) IncrementTrxCount(ctx context.Context, incrementBy uint) error {
	_, err := db.stateCollection().UpdateOne(
		ctx,
		bson.M{kFiStatePk: kPkStateTrxCount},
		bson.D{
			{"$inc", bson.D{{kPkStateTrxCount, int64(incrementBy)}}},
		},
		options.Update().SetUpsert(true),
	)

	return err
}

// blockCollection returns the state collection.
func (db *MongoDb) stateCollection() *mongo.Collection {
	return db.db.Collection(kCoState)
}
