package db

import (
	"context"
	"ftm-explorer/internal/repository/db/types"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// kCoTransactions is the name of the transactions collection.
	kCoTransactions = "transaction"

	// kFiTransactionAddresses is the name of the transaction addresses field.
	kFiTransactionAddresses = "addresses"

	// kFiTransactionHash is the name of the transaction hash field.
	kFiTransactionHash = "hash"

	// kFiTransactionTimestamp is the name of the transaction timestamp field.
	kFiTransactionTimestamp = "timestamp"
)

// AddTransactions adds transactions to the database.
func (db *MongoDb) AddTransactions(ctx context.Context, txs []db_types.Transaction) error {
	interfaceTxs := make([]interface{}, len(txs))
	for i, tx := range txs {
		interfaceTxs[i] = tx
	}

	// try to do the insert
	if _, err := db.transactionCollection().InsertMany(ctx, interfaceTxs); err != nil {
		db.log.Critical(err)
		return err
	}

	return nil
}

// TransactionsWhereAddress returns transactions where the given address is involved.
func (db *MongoDb) TransactionsWhereAddress(ctx context.Context, addr common.Address) ([]db_types.Transaction, error) {
	var transactions []db_types.Transaction

	// Perform the query
	cur, err := db.transactionCollection().Find(ctx, bson.M{kFiTransactionAddresses: addr})
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := cur.Close(ctx); e != nil {
			db.log.Critical(e)
		}
	}()

	// Decode the results
	for cur.Next(ctx) {
		var tx db_types.Transaction
		if err := cur.Decode(&tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	// Check for errors from iterating over rows.
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ShrinkTransactions shrinks the transactions collection. It will persist the given number of transactions.
// It will delete the oldest transactions.
func (db *MongoDb) ShrinkTransactions(ctx context.Context, count int64) error {
	// get the number of transactions
	numOfTrx, err := db.transactionCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	// if there are less transactions than the given count, do nothing
	if numOfTrx <= count {
		return nil
	}

	// Find the timestamp of the Xth most recent record.
	opts := options.FindOne().SetSort(bson.D{{kFiTransactionTimestamp, -1}}).SetSkip(count - 1)
	var result db_types.Transaction
	if err := db.transactionCollection().FindOne(ctx, bson.D{}, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle no document found logically if needed
			return nil
		}
		return err
	}
	cutoffTimestamp := result.Timestamp

	// Delete all records older than the found timestamp.
	deleteFilter := bson.M{kFiTransactionTimestamp: bson.M{"$lt": cutoffTimestamp}}
	_, err = db.transactionCollection().DeleteMany(ctx, deleteFilter)

	return err
}

// transactionCollection returns the transaction collection.
func (db *MongoDb) transactionCollection() *mongo.Collection {
	return db.db.Collection(kCoTransactions)
}

// initTransactionCollection initializes the transaction collection.
func (db *MongoDb) initTransactionCollection() {
	// prepare index models
	ix := make([]mongo.IndexModel, 0)

	// index the timestamp
	ix = append(ix, mongo.IndexModel{Keys: bson.D{{Key: kFiTransactionTimestamp, Value: -1}}})

	// index the addresses
	ix = append(ix, mongo.IndexModel{Keys: bson.D{{Key: kFiTransactionAddresses, Value: 1}}})

	// create indexes
	ctx, cancel := context.WithTimeout(context.Background(), kMongoDefaultTimeout)
	defer cancel()
	if _, err := db.transactionCollection().Indexes().CreateMany(ctx, ix); err != nil {
		db.log.Panicf("can not create indexes for transaction collection; %v", err)
	}

	db.log.Debugf("transactions collection initialized")
}
