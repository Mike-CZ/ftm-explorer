package db

import (
	"context"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	db_types "ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	// define block to add
	block := types.Block{
		Number:    1,
		GasUsed:   699_999,
		Timestamp: 1_689_601_270,
		Transactions: []common.Hash{
			common.HexToHash("0x1"),
			common.HexToHash("0x2"),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add block
	if err := db.AddBlock(ctx, &block); err != nil {
		t.Fatalf("failed to add block: %v", err)
	}

	// get block
	returnedBlock, err := db.Block(ctx, uint64(block.Number))
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
	if returnedBlock.TxsCount != int32(len(block.Transactions)) {
		t.Fatalf("expected block transactions length %d, got %d", len(block.Transactions), returnedBlock.TxsCount)
	}
	if returnedBlock.Timestamp != int64(block.Timestamp) {
		t.Fatalf("expected block timestamp %d, got %d", int64(block.Timestamp), returnedBlock.Timestamp)
	}
}

// Test getting transactions per day from MongoDB
func TestMongoDb_GetTransactionsPerDay(t *testing.T) {
	db := startMongoDb(t)

	// start on 21st of February 2000 at 5:00 UTC
	ts := time.Date(2000, 2, 15, 5, 0, 0, 0, time.UTC)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add 10.000 blocks
	// we add 1 block per 6 hours, so we have 4 blocks per day
	for i := 0; i < 10_000; i++ {
		// add 1 transaction per second
		block := types.Block{
			Number:    hexutil.Uint64(i),
			GasUsed:   699_999,
			Timestamp: hexutil.Uint64(ts.Unix()),
		}
		// add i % 10 + 1 transactions per block, so we have 1 to 10 transactions per block
		for j := 0; j < i%10+1; j++ {
			block.Transactions = append(block.Transactions, common.HexToHash("0x1"))
		}

		// add block
		if err := db.AddBlock(ctx, &block); err != nil {
			t.Fatalf("failed to add block: %v", err)
		}

		// add 6 hours only when there will be next iteration
		if i < 9_999 {
			ts = ts.Add(time.Hour * 6)
		}
	}

	// we added in total 10.000 blocks, so we expect 2.500 days, because we have 4 blocks per day.
	// +2 is because we want to test boundaries
	expectedDays := uint(2_500) + 2

	// shift start by 24h where we should get 0 transactions
	start := ts.Add(time.Hour * 24)
	tpd, err := db.TrxCountAggByTimestamp(ctx, uint64(start.Unix()), types.AggResolutionDay.ToDuration(), expectedDays)
	if err != nil {
		t.Fatalf("failed to get transactions per day: %v", err)
	}

	// check returned number of days
	if len(tpd) != int(expectedDays) {
		t.Fatalf("expected %d days, got %d", expectedDays, len(tpd))
	}

	// check the boundaries
	if tpd[0].Value != 0 {
		t.Fatalf("expected 0 entries, got %d", tpd[0].Value)
	}
	if tpd[len(tpd)-1].Value != 0 {
		t.Fatalf("expected 0 entries, got %d", tpd[len(tpd)-1].Value)
	}

	// check transactions per day
	for i, r := range tpd[1 : len(tpd)-1] {
		day := i * 4
		expected := (day%10 + 1) + ((day+1)%10 + 1) + ((day+2)%10 + 1) + ((day+3)%10 + 1)
		if r.Value != hexutil.Uint64(expected) {
			t.Fatalf("expected %d transactions, got %d", expected, r.Value)
		}
	}
}

// Test getting gas used per day from MongoDB. The main focus is on big numbers.
func TestMongoDb_GetGasUsedPerDay(t *testing.T) {
	db := startMongoDb(t)

	// start on 21st of February 2000 at 5:00 UTC
	ts := time.Date(2000, 2, 15, 5, 0, 0, 0, time.UTC)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add 10.000 blocks
	// we add 1 block per 6 hours, so we have 4 blocks per day
	// we will add more gas than usual, so we don't have to add too many blocks
	for i := 0; i < 10_000; i++ {
		block := types.Block{
			Number: hexutil.Uint64(i),
			// add some really huge number
			// let's say 1 tx costs 100_000_000 gas, we want 1 block per second with ~1000 transactions
			// that means 100_000_000 * 1000 * 60 * 60 * 24 = 8_640_000_000_000_000
			// and since we have 4 blocks per day, we need to divide it by 4 = 2_160_000_000_000_000
			GasUsed:   2_160_000_000_000_000,
			Timestamp: hexutil.Uint64(ts.Unix()),
			Transactions: []common.Hash{
				common.HexToHash("0x1"),
			},
		}

		// add block
		if err := db.AddBlock(ctx, &block); err != nil {
			t.Fatalf("failed to add block: %v", err)
		}

		// add 6 hours only when there will be next iteration
		if i < 9_999 {
			ts = ts.Add(time.Hour * 6)
		}
	}

	// we added in total 10.000 blocks, so we expect 2.500 days, because we have 4 blocks per day.
	// +2 is because we want to test boundaries
	expectedDays := uint(2_500) + 2

	// shift start by 24h where we should get 0 gas used
	start := ts.Add(time.Hour * 24)
	gpd, err := db.GasUsedAggByTimestamp(ctx, uint64(start.Unix()), types.AggResolutionDay.ToDuration(), expectedDays)
	if err != nil {
		t.Fatalf("failed to get gas used per day: %v", err)
	}

	// check returned number of days
	if len(gpd) != int(expectedDays) {
		t.Fatalf("expected %d days, got %d", expectedDays, len(gpd))
	}

	// check the boundaries
	if gpd[0].Value != 0 {
		t.Fatalf("expected 0 entries, got %d", gpd[0].Value)
	}
	if gpd[len(gpd)-1].Value != 0 {
		t.Fatalf("expected 0 entries, got %d", gpd[len(gpd)-1].Value)
	}

	// check gas used per day
	expected := uint64(2_160_000_000_000_000 * 4)
	for _, r := range gpd[1 : len(gpd)-1] {
		if r.Value != hexutil.Uint64(expected) {
			t.Fatalf("expected %d gas used, got %d", expected, r.Value)
		}
	}
}

// Test incrementing and getting trx count from MongoDB.
func TestMongoDb_IncrementAndGetTrxCount(t *testing.T) {
	db := startMongoDb(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// try getting trx count without previous increment
	trxCount, err := db.TrxCount(ctx)
	if err != nil {
		t.Fatalf("failed to get trx count: %v", err)
	}
	if trxCount != 0 {
		t.Fatalf("expected 0 trx count, got %d", trxCount)
	}

	// increment trx count
	if err := db.IncrementTrxCount(ctx, 5); err != nil {
		t.Fatalf("failed to increment trx count: %v", err)
	}

	// try getting trx count after increment
	trxCount, err = db.TrxCount(ctx)
	if err != nil {
		t.Fatalf("failed to get trx count: %v", err)
	}
	if trxCount != 5 {
		t.Fatalf("expected 5 trx count, got %d", trxCount)
	}

	// increment trx count again
	if err := db.IncrementTrxCount(ctx, 10); err != nil {
		t.Fatalf("failed to increment trx count: %v", err)
	}

	// try getting trx count after increment
	trxCount, err = db.TrxCount(ctx)
	if err != nil {
		t.Fatalf("failed to get trx count: %v", err)
	}
	if trxCount != 15 {
		t.Fatalf("expected 15 trx count, got %d", trxCount)
	}
}

// Test adding tokens request.
func TestMongoDb_AddAndGetLatestTokensRequest(t *testing.T) {
	db := startMongoDb(t)

	// define tokens request to add
	tr := types.TokensRequest{
		IpAddress: "192.168.0.1",
		Phrase:    "test phrase",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add tokens request
	if err := db.AddTokensRequest(ctx, &tr); err != nil {
		t.Fatalf("failed to add tokens request: %v", err)
	}

	// get latest tokens request
	latestTr, err := db.LatestTokensRequest(ctx, tr.IpAddress)
	if err != nil {
		t.Fatalf("failed to get latest tokens request: %v", err)
	}
	if latestTr.Phrase != tr.Phrase {
		t.Fatalf("expected phrase %s, got %s", tr.Phrase, latestTr.Phrase)
	}

	// try adding another tokens request
	tr.Phrase = "test phrase 2"
	if err := db.AddTokensRequest(ctx, &tr); err != nil {
		t.Fatalf("failed to add tokens request: %v", err)
	}

	// get latest tokens request
	latestTr, err = db.LatestTokensRequest(ctx, tr.IpAddress)
	if err != nil {
		t.Fatalf("failed to get latest tokens request: %v", err)
	}
	if latestTr.Phrase != tr.Phrase {
		t.Fatalf("expected phrase %s, got %s", tr.Phrase, latestTr.Phrase)
	}
}

// Test updating tokens request.
func TestMongoDb_UpdateTokensRequest(t *testing.T) {
	db := startMongoDb(t)

	// define tokens request to add
	tr := types.TokensRequest{
		IpAddress: "192.168.0.1",
		Phrase:    "test phrase",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add tokens request
	if err := db.AddTokensRequest(ctx, &tr); err != nil {
		t.Fatalf("failed to add tokens request: %v", err)
	}

	// get latest tokens request
	latestTr, err := db.LatestTokensRequest(ctx, tr.IpAddress)
	if err != nil {
		t.Fatalf("failed to get latest tokens request: %v", err)
	}
	if latestTr.Phrase != tr.Phrase {
		t.Fatalf("expected phrase %s, got %s", tr.Phrase, latestTr.Phrase)
	}

	if latestTr.ClaimedAt != nil {
		t.Fatalf("expected nil, got %v", latestTr.ClaimedAt)
	}

	// set claimed at
	claimedAt := time.Now().Unix()
	latestTr.ClaimedAt = &claimedAt

	// update tokens request
	if err := db.UpdateTokensRequest(ctx, latestTr); err != nil {
		t.Fatalf("failed to update tokens request: %v", err)
	}

	// get latest tokens request
	latestTr, err = db.LatestTokensRequest(ctx, tr.IpAddress)
	if err != nil {
		t.Fatalf("failed to get latest tokens request: %v", err)
	}

	if latestTr.ClaimedAt == nil {
		t.Fatalf("expected %d, got nil", claimedAt)
	}

	if *latestTr.ClaimedAt != claimedAt {
		t.Fatalf("expected %d, got %d", claimedAt, latestTr.ClaimedAt)
	}
}

// Test adding time to finality.
func TestMongoDb_AddTimeToFinality(t *testing.T) {
	db := startMongoDb(t)

	// define time to finality to add
	ttf := types.Ttf{
		Timestamp: 1_689_601_270,
		Value:     1.234,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add time to finality
	if err := db.AddTimeToFinality(ctx, &ttf); err != nil {
		t.Fatalf("failed to add time to finality: %v", err)
	}
}

// Test getting time to finality per day from MongoDB.
func TestMongoDb_GetTimeToFinalityAggregation(t *testing.T) {
	db := startMongoDb(t)

	ttfs := []types.Ttf{
		{Timestamp: 1_689_601_271, Value: 10},
		{Timestamp: 1_689_601_278, Value: 5},
		{Timestamp: 1_689_601_282, Value: 1},
		{Timestamp: 1_689_601_290, Value: 2},
	}

	// add ttfs into db
	for _, ttf := range ttfs {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

		if err := db.AddTimeToFinality(ctx, &ttf); err != nil {
			t.Fatalf("failed to add time to finality: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// get time to finality aggregation
	ttfAgg, err := db.TtfAvgAggByTimestamp(ctx, 1_689_601_290, types.AggResolutionSeconds.ToDuration(), 3)
	if err != nil {
		t.Fatalf("failed to get time to finality aggregation: %v", err)
	}

	// check returned number of ticks
	// we expect 2 ticks, because agg `seconds` is 10 seconds and our range is 71 - 90 (in ts)
	if len(ttfAgg) != 3 {
		t.Fatalf("expected 3 ticks, got %d", len(ttfAgg))
	}

	// check time to finality aggregation
	if ttfAgg[1].Value != 7.5 {
		t.Fatalf("expected 7.5, got %f", ttfAgg[0].Value)
	}
	if ttfAgg[2].Value != 1.5 {
		t.Fatalf("expected 1.5, got %f", ttfAgg[1].Value)
	}
}

// Test transactions can be added to MongoDB.
func TestMongoDb_AddAndGetTransactions(t *testing.T) {
	db := startMongoDb(t)

	// define transactions to add
	txs := []db_types.Transaction{
		{
			Addresses: []common.Address{
				common.HexToAddress("0x1"),
				common.HexToAddress("0x2"),
			},
			Hash:      common.Hash{0x12},
			Timestamp: 1_689_601_270,
		},
		{
			Addresses: []common.Address{
				common.HexToAddress("0x3"),
				common.HexToAddress("0x4"),
				common.HexToAddress("0x5"),
			},
			Hash:      common.Hash{0x34},
			Timestamp: 1_689_601_271,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// add transactions
	if err := db.AddTransactions(ctx, txs); err != nil {
		t.Fatalf("failed to add transactions: %v", err)
	}

	// get transactions
	returnedTxs, err := db.TransactionsWhereAddress(ctx, common.HexToAddress("0x1"))
	if err != nil {
		t.Fatalf("failed to get transactions: %v", err)
	}

	// compare transactions
	if len(returnedTxs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(returnedTxs))
	}
	if returnedTxs[0].Hash != txs[0].Hash {
		t.Fatalf("expected hash %s, got %s", txs[0].Hash.Hex(), returnedTxs[0].Hash.Hex())
	}
	if returnedTxs[0].Timestamp != txs[0].Timestamp {
		t.Fatalf("expected timestamp %d, got %d", txs[0].Timestamp, returnedTxs[0].Timestamp)
	}

	// get transactions for another address that is not in the database
	returnedTxs, err = db.TransactionsWhereAddress(ctx, common.HexToAddress("0x6"))
	if err != nil {
		t.Fatalf("failed to get transactions: %v", err)
	}

	// compare transactions
	if len(returnedTxs) != 0 {
		t.Fatalf("expected 0 transaction, got %d", len(returnedTxs))
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
			Host: "localhost",
			Port: p.Int(),
			Db:   dbName,
		},
		log: logger.NewMockLogger(),
	}
	db, err := NewMongoDb(args.cfg, args.log)
	if err != nil {
		t.Fatalf("NewMongoDb() error = %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	return db
}
