package db

import (
	"context"
	"errors"
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

// LatestUnclaimedTokensRequest returns the latest unclaimed tokens request for the given ip address.
func (db *MongoDb) LatestUnclaimedTokensRequest(ctx context.Context, ipAddress string) (*types.TokensRequest, error) {
	var tr types.TokensRequest

	// try to find the latest tokens request for the given ip address
	opts := options.FindOne().SetSort(bson.D{{kFiTokensRequestPk, -1}})
	if err := db.tokensRequestCollection().FindOne(ctx, bson.M{
		kFiTokensRequestIp:        ipAddress,
		kFiTokensRequestClaimedAt: nil,
	}, opts).Decode(&tr); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			db.log.Debugf("no unclaimed tokens request found for ip address: %s", ipAddress)
			return nil, nil
		}
		db.log.Criticalf("failed to get unclaimed tokens request for ip address: %s. err: %v", ipAddress, err)
		return nil, err
	}

	return &tr, nil
}

// LatestClaimedTokensRequests returns the latest claimed tokens requests for the given ip address.
// The from parameter is used to determine the starting point of the query.
func (db *MongoDb) LatestClaimedTokensRequests(ctx context.Context, ipAddress string, from uint64) ([]types.TokensRequest, error) {
	var requests []types.TokensRequest

	opts := options.Find().SetSort(bson.D{{kFiTokensRequestPk, -1}})
	filter := bson.M{
		kFiTokensRequestIp:        ipAddress,
		kFiTokensRequestClaimedAt: bson.M{"$gte": from},
	}
	cursor, err := db.tokensRequestCollection().Find(ctx, filter, opts)
	if err != nil {
		db.log.Criticalf("failed to get claimed tokens requests for ip address: %s. err: %v", ipAddress, err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			db.log.Criticalf("failed to close cursor for claimed tokens requests for ip address: %s. err: %v", ipAddress, err)
		}
	}()

	for cursor.Next(ctx) {
		var request types.TokensRequest
		if err := cursor.Decode(&request); err != nil {
			db.log.Criticalf("failed to decode claimed tokens request for ip address: %s. err: %v", ipAddress, err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := cursor.Err(); err != nil {
		db.log.Criticalf("failed to get claimed tokens requests for ip address: %s. err: %v", ipAddress, err)
		return nil, err // Return nil and the error if an error occurs during iteration
	}

	return requests, nil
}

// blockCollection returns the tokens request collection.
func (db *MongoDb) tokensRequestCollection() *mongo.Collection {
	return db.db.Collection(kCoTokensRequest)
}
