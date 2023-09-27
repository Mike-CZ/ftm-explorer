package db

import (
	"context"
	"fmt"
	"ftm-explorer/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// kCoState is the name of the state collection.
	kCoTimeToFinality = "time_to_finality"

	// kFiTokensRequestIp is the name of the ip address field.
	kFiTimeToFinalityTimestamp = "timestamp"

	// kFiTokensRequestPhrase is the name of the phrase field.
	kFiTimeToFinalityValue = "value"
)

// AddTimeToFinality adds the given time to finality.
func (db *MongoDb) AddTimeToFinality(ctx context.Context, tr *types.Ttf) error {
	if tr == nil {
		return fmt.Errorf("can not add empty tokens request")
	}

	// try to do the insert
	if _, err := db.timeToFinalityCollection().InsertOne(ctx, tr); err != nil {
		db.log.Critical(err)
		return err
	}

	db.log.Debugf("time to finality added. timestamp: %d, value: %f", tr.Timestamp, tr.Value)

	return nil
}

// TtfAvgAggByTimestamp returns average aggregation of time to finality in given time range.
func (db *MongoDb) TtfAvgAggByTimestamp(ctx context.Context, endTime uint64, resolution uint, ticks uint) ([]types.FloatTick, error) {
	type aggregationResult struct {
		Key    int64   `bson:"_id"`
		Result float64 `bson:"aggregation"`
	}

	// define the MongoDB Pipeline for the aggregation
	pipeline := mongo.Pipeline{
		// This stage filters the documents to only pass those with a 'kFiTimeToFinalityTimestamp' value
		// greater than 'endTime - (ticks*resolution)' and less than or equal to 'endTime'.
		{{"$match", bson.D{
			{kFiTimeToFinalityTimestamp, bson.D{
				{"$gt", endTime - uint64(ticks*resolution)}, {"$lte", endTime}},
			},
		}}},
		// This stage groups the documents by a generated '_id' field.
		{{"$group", bson.D{
			// Subtract the `endTime` value from the calculated group distance between `endTime` and `kFiTimeToFinalityTimestamp`.
			{"_id", bson.D{{"$subtract", []interface{}{
				endTime,
				// Multiply the result of division by 'resolution'. The result will be a group of the distance between
				// 'endTime' and 'kFiTimeToFinalityTimestamp'.
				bson.D{{"$multiply", []interface{}{
					bson.D{
						// Convert the result of division into integer.
						{"$toInt", bson.D{
							// Divide the difference between 'endTime' and 'kFiTimeToFinalityTimestamp' by 'resolution'.
							// This will create a value that is the number of 'resolution' intervals between
							// 'endTime' and 'kFiTimeToFinalityTimestamp'.
							{"$divide", []interface{}{
								// Subtract the 'kFiTimeToFinalityTimestamp' value from 'endTime'.
								bson.D{{"$subtract", []interface{}{endTime, fmt.Sprintf("$%s", kFiTimeToFinalityTimestamp)}}},
								resolution,
							}},
						}},
					},
					resolution,
				}}},
			}}}},
			// Avg the values of the 'aggField' field for each group.
			{"aggregation", bson.D{{"$avg", fmt.Sprintf("$%s", kFiTimeToFinalityValue)}}},
		}}},
	}

	// execute the aggregation
	cursor, err := db.timeToFinalityCollection().Aggregate(ctx, pipeline)
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
	resMap := make(map[uint64]float64)
	for _, result := range results {
		resMap[uint64(result.Key)] = result.Result
	}

	// prepare the result
	ticksResult := make([]types.FloatTick, ticks)
	for i, ts := uint(0), endTime-uint64(resolution*(ticks-1)); i < ticks; i, ts = i+1, ts+uint64(resolution) {
		ticksResult[i] = types.FloatTick{
			Time: ts,
		}
		// check if there is some data for the entry
		if val, ok := resMap[ts]; ok {
			ticksResult[i].Value = float64(int(val*100)) / 100
		}
	}

	return ticksResult, nil
}

// timeToFinalityCollection returns the time to finality collection.
func (db *MongoDb) timeToFinalityCollection() *mongo.Collection {
	return db.db.Collection(kCoTimeToFinality)
}
