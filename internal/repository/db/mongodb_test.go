package db

import (
	"context"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/types"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Test connection to MongoDB
func TestMongoDb_Connection(t *testing.T) {
	_ = startMongoDb(t)
}

// Test adding and getting a block from MongoDB
func TestMongoDb_AddAndGetBlock(t *testing.T) {
	db := startMongoDb(t)
	defer func() {
		db.Close()
	}()

	// define block to add
	block := types.Block{
		Number:  1,
		GasUsed: 699_999,
		Transactions: []common.Hash{
			common.HexToHash("0x1"),
			common.HexToHash("0x2"),
		},
	}

	// add block
	if err := db.AddBlock(&block); err != nil {
		t.Fatalf("failed to add block: %v", err)
	}

	// get block
	returnedBlock, err := db.GetBlock(int64(block.Number))
	if err != nil {
		t.Fatalf("failed to get block: %v", err)
	}

	// compare blocks
	if returnedBlock.Number != int64(block.Number) {
		t.Fatalf("expected block number %d, got %d", int64(block.Number), returnedBlock.Number)
	}
	if returnedBlock.GasUsed != int64(block.GasUsed) {
		t.Fatalf("expected block gas used %d, got %d", int64(block.GasUsed), returnedBlock.GasUsed)
	}
	if returnedBlock.TxCount != int32(len(block.Transactions)) {
		t.Fatalf("expected block transactions length %d, got %d", len(block.Transactions), returnedBlock.TxCount)
	}
	if returnedBlock.Timestamp != int64(block.Timestamp) {
		t.Fatalf("expected block timestamp %d, got %d", int64(block.Timestamp), returnedBlock.Timestamp)
	}
}

// startMongoDb starts MongoDB in a Docker container and returns the MongoDb instance.
func startMongoDb(t *testing.T) *MongoDb {
	t.Helper()

	dbName := "test_db"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// start MongoDB in a Docker container
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:6.0.8",
			ExposedPorts: []string{"27017/tcp"},
			Env:          map[string]string{"MONGO_INITDB_DATABASE": dbName},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	})
	p, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	// create new MongoDb instance
	args := struct {
		cfg *config.MongoDb
		log logger.ILogger
	}{
		cfg: &config.MongoDb{
			Host:     "localhost",
			Port:     p.Int(),
			Database: dbName,
		},
		log: logger.NewMockLogger(),
	}
	db, err := NewMongoDb(args.cfg, args.log)
	if err != nil {
		t.Fatalf("NewMongoDb() error = %v", err)
	}

	return db
}
