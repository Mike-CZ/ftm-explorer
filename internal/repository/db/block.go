package db

import (
	"context"
	"fmt"
	"ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// kCoBlocks is the name of the block collection.
	kCoBlocks = "block"

	// kFiBlockNumber is the name of the block number field. It is also the primary key.
	kFiBlockNumber = "_id"

	// kFiBlockTxCount is the name of the block transaction count field.
	kFiBlockTxCount = "txsCount"

	// kFiBlockGasUsed is the name of the block gas used field.
	kFiBlockGasUsed = "gasUsed"

	// kFiBlockTimestamp is the name of the block timestamp field.
	kFiBlockTimestamp = "timestamp"
)

// TrxCountAggByTimestamp returns aggregation of transactions in given time range.
func (db *MongoDb) TrxCountAggByTimestamp(ctx context.Context, endTime uint64, resolution uint, ticks uint) ([]types.HexUintTick, error) {
	resMap, err := db.getBlkAggByTimestamp(ctx, endTime, resolution, ticks, kFiBlockTxCount)
	if err != nil {
		db.log.Critical(err)
		return nil, err
	}

	// prepare the result
	ticksResult := make([]types.HexUintTick, ticks)
	for i, ts := uint(0), endTime-uint64(resolution*(ticks-1)); i < ticks; i, ts = i+1, ts+uint64(resolution) {
		ticksResult[i] = types.HexUintTick{
			Time: ts,
		}
		// check if there is some data for the entry
		if val, ok := resMap[ts]; ok {
			ticksResult[i].Value = hexutil.Uint64(val)
		}
	}

	return ticksResult, nil
}

// GasUsedAggByTimestamp returns aggregation of gas used in given time range.
func (db *MongoDb) GasUsedAggByTimestamp(ctx context.Context, endTime uint64, resolution uint, ticks uint) ([]types.HexUintTick, error) {
	resMap, err := db.getBlkAggByTimestamp(ctx, endTime, resolution, ticks, kFiBlockGasUsed)
	if err != nil {
		db.log.Critical(err)
		return nil, err
	}

	// prepare the result
	ticksResult := make([]types.HexUintTick, ticks)
	for i, ts := uint(0), endTime-uint64(resolution*(ticks-1)); i < ticks; i, ts = i+1, ts+uint64(resolution) {
		ticksResult[i] = types.HexUintTick{
			Time: ts,
		}
		// check if there is some data for the entry
		if val, ok := resMap[ts]; ok {
			ticksResult[i].Value = hexutil.Uint64(val)
		}
	}

	return ticksResult, nil
}

// AddBlock adds a block to the database.
func (db *MongoDb) AddBlock(ctx context.Context, block *types.Block) error {
	if block == nil {
		return fmt.Errorf("can not add empty block")
	}

	// try to do the insert
	if _, err := db.blockCollection().InsertOne(ctx, &db_types.Block{
		Number:    int64(block.Number),
		TxsCount:  int32(len(block.Transactions)),
		GasUsed:   int64(block.GasUsed),
		Timestamp: int64(block.Timestamp),
	}); err != nil {
		db.log.Critical(err)
		return err
	}

	db.log.Debugf("block %d added to database", int64(block.Number))

	return nil
}

// Block returns the block with the given number.
func (db *MongoDb) Block(ctx context.Context, number uint64) (*db_types.Block, error) {
	// try to get the block
	var block db_types.Block
	if err := db.blockCollection().FindOne(ctx, bson.D{{Key: kFiBlockNumber, Value: number}}).Decode(&block); err != nil {
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
		db.log.Panicf("can not create indexes for block collection; %v", err)
	}

	db.log.Debugf("blocks collection initialized")
}

// getBlkAggByTimestamp returns the aggregation result for the given time range.
// The aggregation result is a map where the key is the timestamp and the value is the result.
// The aggregation is calculated since the `endTime - (ticks*resolution)` to `endTime`.
// The key is always the last timestamp of the aggregation period.
// params:
//
//	ctx        - the context
//	endTime    - the end time
//	ticks      - the number of ticks to return in backward direction
//	resolution - the resolution in seconds (1 minute = 60, 1 hour = 3600, 1 day = 86400)
//	aggField   - the field to aggregate
func (db *MongoDb) getBlkAggByTimestamp(ctx context.Context, endTime uint64, resolution uint, ticks uint, aggField string) (map[uint64]int64, error) {
	type aggregationResult struct {
		Key    int64 `bson:"_id"`
		Result int64 `bson:"aggregation"`
	}

	// define the MongoDB Pipeline for the aggregation
	pipeline := mongo.Pipeline{
		// This stage filters the documents to only pass those with a 'kFiBlockTimestamp' value
		// greater than 'endTime - (ticks*resolution)' and less than or equal to 'endTime'.
		{{"$match", bson.D{
			{kFiBlockTimestamp, bson.D{
				{"$gt", endTime - uint64(ticks*resolution)}, {"$lte", endTime}},
			},
		}}},
		// This stage groups the documents by a generated '_id' field.
		{{"$group", bson.D{
			// Subtract the `endTime` value from the calculated group distance between `endTime` and `kFiBlockTimestamp`.
			{"_id", bson.D{{"$subtract", []interface{}{
				endTime,
				// Multiply the result of division by 'resolution'. The result will be a group of the distance between
				// 'endTime' and 'kFiBlockTimestamp'.
				bson.D{{"$multiply", []interface{}{
					bson.D{
						// Convert the result of division into integer.
						{"$toInt", bson.D{
							// Divide the difference between 'endTime' and 'kFiBlockTimestamp' by 'resolution'.
							// This will create a value that is the number of 'resolution' intervals between
							// 'endTime' and 'kFiBlockTimestamp'.
							{"$divide", []interface{}{
								// Subtract the 'kFiBlockTimestamp' value from 'endTime'.
								bson.D{{"$subtract", []interface{}{endTime, fmt.Sprintf("$%s", kFiBlockTimestamp)}}},
								resolution,
							}},
						}},
					},
					resolution,
				}}},
			}}}},
			// Sum the values of the 'aggField' field for each group.
			{"aggregation", bson.D{{"$sum", fmt.Sprintf("$%s", aggField)}}},
		}}},
	}

	// execute the aggregation
	cursor, err := db.blockCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	// get the results
	var results []aggregationResult
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	// convert the results to the map
	res := make(map[uint64]int64)
	for _, result := range results {
		res[uint64(result.Key)] = result.Result
	}

	return res, nil
}
