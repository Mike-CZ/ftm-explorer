package db

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// kCoAccounts is the name of the accounts collection.
	kCoAccounts = "accounts"

	// kFiAccountAddress is the name of the account address field.
	kFiAccountAddress = "_id"

	// kFiAccountLastSeen is the name of the account last seen field.
	kFiAccountLastSeen = "lastSeen"
)

// AddAccounts adds accounts to the database.
func (db *MongoDb) AddAccounts(ctx context.Context, accs []common.Address, stamp int64) error {
	for _, acc := range accs {
		filter := bson.D{{kFiAccountAddress, acc}}
		update := bson.D{
			{"$set", bson.D{
				{kFiAccountLastSeen, stamp},
			}},
		}
		opts := options.Update().SetUpsert(true)

		if _, err := db.accountCollection().UpdateOne(ctx, filter, update, opts); err != nil {
			db.log.Criticalf("error updating account %s: %v", acc.Hex(), err)
			return err
		}
	}

	return nil
}

// NumberOfAccoutns returns the number of accounts in the database.
func (db *MongoDb) NumberOfAccoutns(ctx context.Context) (uint64, error) {
	count, err := db.accountCollection().CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

// accountCollection returns the accounts collection.
func (db *MongoDb) accountCollection() *mongo.Collection {
	return db.db.Collection(kCoAccounts)
}
