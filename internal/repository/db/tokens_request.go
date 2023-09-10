package db

import (
	"context"
	"fmt"
	"ftm-explorer/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// kCoState is the name of the state collection.
	kCoTokensRequest = "tokens_request"

	// kFiStatePk is the name of the primary key of the state collection.
	kFiTokensRequestPk = "_id"

	// kFiTokensRequestIp is the name of the ip address field.
	kFiTokensRequestIp = "ip"

	// kFiTokensRequestPhrase is the name of the phrase field.
	kFiTokensRequestPhrase = "phrase"

	// kFiTokensRequestReceiver is the name of the receiver field.
	kFiTokensRequestReceiver = "receiver"

	// kFiTokensRequestClaimedAt is the name of the claimed at field.
	kFiTokensRequestClaimedAt = "claimed_at"
)

// AddTokensRequest adds a new tokens request to the database.
func (db *MongoDb) AddTokensRequest(ctx context.Context, tr *types.TokensRequest) error {
	if tr == nil {
		return fmt.Errorf("can not add empty tokens request")
	}

	// try to do the insert
	tr.Id = primitive.NewObjectID()
	if _, err := db.tokensRequestCollection().InsertOne(ctx, tr); err != nil {
		tr.Id = primitive.NilObjectID
		db.log.Critical(err)
		return err
	}

	db.log.Debugf("tokens request added. ip: %s, phrase: %s", tr.IpAddress, tr.Phrase)

	return nil
}

// UpdateTokensRequest updates the given tokens request.
func (db *MongoDb) UpdateTokensRequest(ctx context.Context, tr *types.TokensRequest) error {
	if tr.Id == primitive.NilObjectID {
		return fmt.Errorf("can not update tokens request with empty id")
	}

	// Filter by the ID of the TokensRequest
	filter := bson.M{kFiTokensRequestPk: tr.Id}

	// Define update operation
	update := bson.M{
		"$set": bson.M{
			kFiTokensRequestIp:        tr.IpAddress,
			kFiTokensRequestPhrase:    tr.Phrase,
			kFiTokensRequestReceiver:  tr.Receiver,
			kFiTokensRequestClaimedAt: tr.ClaimedAt,
		},
	}

	// Apply the update to the database
	updateResult, err := db.tokensRequestCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		db.log.Critical(err)
		return err
	}

	// If no document matched the filter, then return an error
	if updateResult.MatchedCount == 0 {
		return fmt.Errorf("no document found to update with id: %s", tr.Id)
	}

	db.log.Debugf("tokens request updated. id: %s, ip: %s, phrase: %s claimed_at", tr.Id, tr.IpAddress, tr.Phrase, tr.ClaimedAt)
	return nil
}

// LatestTokensRequest returns the latest tokens request for the given ip address.
func (db *MongoDb) LatestTokensRequest(ctx context.Context, ipAddress string) (*types.TokensRequest, error) {
	var tr types.TokensRequest

	// try to find the latest tokens request for the given ip address
	opts := options.FindOne().SetSort(bson.D{{kFiTokensRequestPk, -1}})
	if err := db.tokensRequestCollection().FindOne(ctx, bson.M{kFiTokensRequestIp: ipAddress}, opts).Decode(&tr); err != nil {
		if err == mongo.ErrNoDocuments {
			db.log.Debugf("no tokens request found for ip address: %s", ipAddress)
			return nil, nil
		}
		db.log.Criticalf("failed to get tokens request for ip address: %s. err: %v", ipAddress, err)
		return nil, err
	}

	return &tr, nil
}

// blockCollection returns the tokens request collection.
func (db *MongoDb) tokensRequestCollection() *mongo.Collection {
	return db.db.Collection(kCoTokensRequest)
}
