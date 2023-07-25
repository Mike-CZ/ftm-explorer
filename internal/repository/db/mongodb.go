package db

import (
	"context"
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// kMongoDefaultTimeout is the default timeout for MongoDb database operations.
const kMongoDefaultTimeout = 5 * time.Second

// MongoDb represents a MongoDB database.
type MongoDb struct {
	client *mongo.Client
	db     *mongo.Database
	log    logger.ILogger
}

// NewMongoDb creates new MongoDb instance.
func NewMongoDb(cfg *config.MongoDb, log logger.ILogger) (*MongoDb, error) {
	log.Debugf("connecting mongodb at %s:%d/%s", cfg.Host, cfg.Port, cfg.Db)

	// open the database connection
	con, err := connectDb(cfg)
	if err != nil {
		log.Criticalf("can not contact the database; %s", err.Error())
		return nil, err
	}

	// log the event
	log.Notice("mongodb connection established")

	// return the bridge
	db := &MongoDb{
		client: con,
		db:     con.Database(cfg.Db),
		log:    log,
	}

	// initialize the collections
	db.initBlockCollection()

	return db, nil
}

func (db *MongoDb) Close() {
	if db.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), kMongoDefaultTimeout)
		defer cancel()

		err := db.client.Disconnect(ctx)
		if err != nil {
			db.log.Errorf("error on closing database connection; %v", err)
		}

		db.log.Info("database connection is closed")
	}
}

// connectDb opens MongoDb database connection
func connectDb(cfg *config.MongoDb) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kMongoDefaultTimeout)
	defer cancel()

	// prepare the connection string for authentication
	ucs := ""
	if cfg.User != nil && cfg.Password != nil {
		ucs = fmt.Sprintf("%s:%s@", *cfg.User, *cfg.Password)
	}

	// create new MongoDb client
	cs := fmt.Sprintf("mongodb://%s%s:%d/%s", ucs, cfg.Host, cfg.Port, cfg.Db)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cs))
	if err != nil {
		return nil, err
	}

	// validate the connection was indeed established
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
