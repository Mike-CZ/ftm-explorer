package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ftm-explorer/internal/api/graphql/resolvers"
	"ftm-explorer/internal/api/handlers"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
)

// apiTestCase represents a test case for the API server.
type apiTestCase struct {
	testName      string
	requestBody   string
	buildStubs    func(*repository.MockRepository)
	checkResponse func(*testing.T, *http.Response)
}

// Test cases for the API server.
func TestApiServer_Run(t *testing.T) {
	// initialize stubs
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockLogger := logger.NewMockLogger()

	// initialize test server
	handler := handlers.ApiHandler([]string{"*"}, resolvers.NewResolver(mockRepository, mockLogger), mockLogger)
	server := httptest.NewServer(handler)
	defer server.Close()

	// use table-driven testing to test multiple cases
	testCases := []apiTestCase{
		getTransactionTestCase(t),
		getBlockTestCase(t),
		getRecentBlocksTestCase(t),
		getCurrentBlockHeightTestCase(t),
		getBlockTimestampTxsCountAggregationsTestCase(t),
		getBlockTimestampGasUsedAggregationsTestCase(t),
		getNumberOfAccountsTestCase(t),
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// build stubs
			if tc.buildStubs != nil {
				tc.buildStubs(mockRepository)
			}
			// make request
			resp, err := server.Client().Post(server.URL, "application/json", strings.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}
			t.Cleanup(func() {
				_ = resp.Body.Close()
			})
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status OK, got %v", resp.Status)
			}
			// check response
			if tc.checkResponse != nil {
				tc.checkResponse(t, resp)
			}
		})
	}
}

// getTransactionTestCase returns a test case for a transaction not found error.
func getTransactionTestCase(t *testing.T) apiTestCase {
	trx := getTestTransaction(t)
	return apiTestCase{
		testName:    "GetTransaction",
		requestBody: fmt.Sprintf(`{"query": "query { transaction(hash: \"%s\") { hash, blockHash, blockNumber, from, to, contractAddress, nonce, gas, gasUsed, cumulativeGasUsed, gasPrice, value, input, transactionIndex, status }}"}`, trx.Hash.Hex()),
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetTransactionByHash(gomock.Eq(trx.Hash)).Return(&trx, nil)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into transaction
			trxRes := struct {
				Trx types.Transaction `json:"transaction"`
			}{}
			if err := json.Unmarshal(apiRes.Data, &trxRes); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate transaction
			validateTransaction(t, trx, trxRes.Trx)
		},
	}
}

// getTransactionTestCase returns a test case for a transaction not found error.
func getBlockTestCase(t *testing.T) apiTestCase {
	block := getTestBlock(t)
	return apiTestCase{
		testName:    "Block",
		requestBody: fmt.Sprintf(`{"query": "query { block(number: \"%s\") { number, epoch, hash, parentHash, timestamp, gasLimit, gasUsed, transactions, transactionsCount }}"}`, block.Number.String()),
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetBlockByNumber(gomock.Eq(uint64(block.Number))).Return(&block, nil)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into block
			blockRes := struct {
				Block struct {
					types.Block
					// add transactions count, which is not part of the block struct
					TransactionsCount int32 `json:"transactionsCount"`
				} `json:"block"`
			}{}
			if err := json.Unmarshal(apiRes.Data, &blockRes); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate block
			validateBlock(t, block, blockRes.Block.Block)
			// validate transactions count
			if blockRes.Block.TransactionsCount != int32(len(block.Transactions)) {
				t.Errorf("expected transactions count %v, got %v", len(block.Transactions), blockRes.Block.TransactionsCount)
			}
		},
	}
}

// getRecentBlocksTestCase returns a test case for a recent blocks query.
func getRecentBlocksTestCase(_ *testing.T) apiTestCase {
	blocks := []*types.Block{
		{Number: hexutil.Uint64(1)}, {Number: hexutil.Uint64(2)}, {Number: hexutil.Uint64(3)},
		{Number: hexutil.Uint64(4)}, {Number: hexutil.Uint64(5)},
	}
	return apiTestCase{
		testName:    "GetRecentBlocks",
		requestBody: `{"query": "query { recentBlocks(limit: 5) { number }}"}`,
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetLatestObservedBlocks(gomock.Eq(uint(5))).Return(blocks)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into block
			blockRes := struct {
				Blocks []*types.Block `json:"recentBlocks"`
			}{}
			if err := json.Unmarshal(apiRes.Data, &blockRes); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate blocks
			if len(blockRes.Blocks) != len(blocks) {
				t.Errorf("expected blocks count %v, got %v", len(blocks), len(blockRes.Blocks))
			}
			for i, block := range blocks {
				validateBlock(t, *block, *blockRes.Blocks[i])
			}
		},
	}
}

// getRecentBlocksTestCase returns a test case for a recent blocks query.
func getCurrentBlockHeightTestCase(_ *testing.T) apiTestCase {
	var blockHeight uint64 = 100_000
	return apiTestCase{
		testName:    "GetCurrentBlockHeight",
		requestBody: `{"query": "query { currentBlockHeight}"}`,
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetLatestObservedBlock().Return(&types.Block{Number: hexutil.Uint64(blockHeight)})
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into response
			heightRes := struct {
				BlockHeight hexutil.Uint64 `json:"currentBlockHeight"`
			}{}
			if err := json.Unmarshal(apiRes.Data, &heightRes); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate height
			if uint64(heightRes.BlockHeight) != blockHeight {
				t.Errorf("expected block height %v, got %v", blockHeight, heightRes.BlockHeight)
			}
		},
	}
}

// getBlockTimestampTxsCountAggregationsTestCase returns a test case for a block timestamp trx count aggregations query.
func getBlockTimestampTxsCountAggregationsTestCase(_ *testing.T) apiTestCase {
	agg := []types.HexUintTick{
		{Value: hexutil.Uint64(178), Time: 1_690_099_503},
		{Value: hexutil.Uint64(155), Time: 1_690_099_563},
		{Value: hexutil.Uint64(201), Time: 1_690_099_623},
		{Value: hexutil.Uint64(167), Time: 1_690_099_683},
		{Value: hexutil.Uint64(180), Time: 1_690_099_743},
	}
	return apiTestCase{
		testName:    "GetBlockTimestampAggregations",
		requestBody: `{"query": "query { blockTimestampAggregations(subject: TXS_COUNT, resolution: MINUTE, ticks: 5) { timestamp, value }}"}`,
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetTrxCountAggByTimestamp(gomock.Eq(types.AggResolutionMinute), gomock.Eq(uint(5)), gomock.Nil()).Return(agg, nil)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into response
			var response struct {
				Aggregations []struct {
					Timestamp int32          `json:"timestamp"`
					Value     hexutil.Uint64 `json:"value"`
				} `json:"blockTimestampAggregations"`
			}
			if err := json.Unmarshal(apiRes.Data, &response); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate aggregations
			if len(response.Aggregations) != len(agg) {
				t.Errorf("expected aggregations count %v, got %v", len(agg), len(response.Aggregations))
			}
			for i, tick := range agg {
				if response.Aggregations[i].Timestamp != int32(tick.Time) {
					t.Errorf("expected timestamp %v, got %v", tick.Time, response.Aggregations[i].Timestamp)
				}
				if response.Aggregations[i].Value != tick.Value {
					t.Errorf("expected value %v, got %v", tick.Value, response.Aggregations[i].Value)
				}
			}
		},
	}
}

// getBlockTimestampGasUsedAggregationsTestCase returns a test case for a block timestamp gas used aggregations query.
func getBlockTimestampGasUsedAggregationsTestCase(_ *testing.T) apiTestCase {
	agg := []types.HexUintTick{
		{Value: hexutil.Uint64(105_803_475), Time: 1_690_099_503},
		{Value: hexutil.Uint64(160_550_785), Time: 1_690_099_563},
		{Value: hexutil.Uint64(116_962_544), Time: 1_690_099_623},
		{Value: hexutil.Uint64(115_388_923), Time: 1_690_099_683},
		{Value: hexutil.Uint64(91_255_380), Time: 1_690_099_743},
	}
	endTime := uint64(1_690_100_448)
	return apiTestCase{
		testName:    "GetBlockTimestampAggregations",
		requestBody: `{"query": "query { blockTimestampAggregations(subject: GAS_USED, resolution: HOUR, ticks: 5, endTime: 1690100448) { timestamp, value }}"}`,
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetGasUsedAggByTimestamp(gomock.Eq(types.AggResolutionHour), gomock.Eq(uint(5)), gomock.Eq(&endTime)).Return(agg, nil)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into response
			var response struct {
				Aggregations []struct {
					Timestamp int32          `json:"timestamp"`
					Value     hexutil.Uint64 `json:"value"`
				} `json:"blockTimestampAggregations"`
			}
			if err := json.Unmarshal(apiRes.Data, &response); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate aggregations
			if len(response.Aggregations) != len(agg) {
				t.Errorf("expected aggregations count %v, got %v", len(agg), len(response.Aggregations))
			}
			for i, tick := range agg {
				if response.Aggregations[i].Timestamp != int32(tick.Time) {
					t.Errorf("expected timestamp %v, got %v", tick.Time, response.Aggregations[i].Timestamp)
				}
				if response.Aggregations[i].Value != tick.Value {
					t.Errorf("expected value %v, got %v", tick.Value, response.Aggregations[i].Value)
				}
			}
		},
	}
}

// getRecentBlocksTestCase returns a test case for a recent blocks query.
func getNumberOfAccountsTestCase(_ *testing.T) apiTestCase {
	var number uint64 = 4_250
	return apiTestCase{
		testName:    "GetNumberOfAccounts",
		requestBody: `{"query": "query { numberOfAccounts}"}`,
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetNumberOfAccounts().Return(number)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got: %s", apiRes.Errors[0].Message)
			}
			// decode raw data into response
			numberRes := struct {
				NumberOfAccounts int32 `json:"numberOfAccounts"`
			}{}
			if err := json.Unmarshal(apiRes.Data, &numberRes); err != nil {
				t.Errorf("failed to unmarshall data: %v", err)
			}
			// validate number of accounts
			if numberRes.NumberOfAccounts != int32(number) {
				t.Errorf("expected number of accounts %v, got %v", number, numberRes.NumberOfAccounts)
			}
		},
	}
}

// apiResponse represents a response from the API server.
type apiResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// decodeResponse decodes the response body into the given value.
// If the response contains errors, they will be returned as a slice of strings.
func decodeResponse(t *testing.T, resp *http.Response) apiResponse {
	t.Helper()
	var body apiResponse
	// decode response body
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Errorf("failed to unmarshall response: %v", err)
	}
	return body
}

// getTestTransaction returns a test transaction.
func getTestTransaction(t *testing.T) types.Transaction {
	t.Helper()
	blockHash := common.HexToHash("0x0003298e00000cfb8d70c2c4ff1861a5f12607ad642ed8529cce62d4afe9e1f7")
	blockNumber := hexutil.Uint64(100)
	to := common.HexToAddress("0x8f8ddaca443ceac1ee5676721d14cfc5c4548020")
	gasUsed := hexutil.Uint64(36_944)
	cumulativeGas := hexutil.Uint64(36_944)
	trxIndex := hexutil.Uint64(0)
	status := hexutil.Uint64(1)
	return types.Transaction{
		Hash:              common.HexToHash("0xcd38ab7c7c77700e3d225316f10b1e5cbdff4c294a7034142e5f4a0ef7eebd7f"),
		BlockHash:         &blockHash,
		BlockNumber:       &blockNumber,
		From:              common.HexToAddress("0x3ea8375b450e443a9bb5cac5f4be9f8f646d7318"),
		To:                &to,
		ContractAddress:   nil,
		Nonce:             hexutil.Uint64(813),
		Gas:               hexutil.Uint64(53_719),
		GasUsed:           &gasUsed,
		CumulativeGasUsed: &cumulativeGas,
		GasPrice:          hexutil.Big(*big.NewInt(221_759_311_181)),
		Value:             hexutil.Big(*big.NewInt(0)),
		Input:             common.Hex2Bytes("32dee40b0000000000000000000000002a2f9fdc05656a0bbbf0e373462d614ff1aeffbf"),
		TransactionIndex:  &trxIndex,
		Status:            &status,
	}
}

// getTestTransaction returns a test transaction.
func getTestBlock(t *testing.T) types.Block {
	t.Helper()
	return types.Block{
		Number:     hexutil.Uint64(65_797_494),
		Epoch:      hexutil.Uint64(223_648),
		Hash:       common.HexToHash("0x000369a000001e8b893fcc26dd34add080cc996746e468f728f6cff334722d65"),
		ParentHash: common.HexToHash("0x000369a000001e8062583ee1cbb8044713711793de93f5a9a1e0733ac43929c1"),
		Timestamp:  hexutil.Uint64(1_689_357_196),
		GasLimit:   hexutil.Uint64(281_474_976_710_655),
		GasUsed:    hexutil.Uint64(775_639),
		Transactions: []common.Hash{
			common.HexToHash("0x030b623eb60d312e5e866c401cb04959f450de6ce4184b76eacc2fe17bd29dee"),
			common.HexToHash("0x3ab52dc0c6a376933e78f95a69003130272b1b26b0374cf757de1a3c2f27e209"),
			common.HexToHash("0x3779f09001d187d076daf214c8b255acada58be1921d9bce648836f41753c607"),
			common.HexToHash("0x8cba9acf51367104ee232cf0180b7ddb95d6525783c88da4463cfd5a5d749834"),
			common.HexToHash("0x92f992ecbfeb174c0db9790ea844ec4c117f5fea1edb9a737524e7249b96869b"),
		},
	}
}

// validateBlock validates that the given block equals the expected block.
func validateBlock(t *testing.T, blk types.Block, blkRes types.Block) {
	t.Helper()
	// number
	if blk.Number != blkRes.Number {
		t.Errorf("expected block number %v, got %v", blk.Number, blkRes.Number)
	}
	// epoch
	if blk.Epoch != blkRes.Epoch {
		t.Errorf("expected block epoch %v, got %v", blk.Epoch, blkRes.Epoch)
	}
	// hash
	if blk.Hash != blkRes.Hash {
		t.Errorf("expected block hash %v, got %v", blk.Hash, blkRes.Hash)
	}
	// parent hash
	if blk.ParentHash != blkRes.ParentHash {
		t.Errorf("expected block parent hash %v, got %v", blk.ParentHash, blkRes.ParentHash)
	}
	// timestamp
	if blk.Timestamp != blkRes.Timestamp {
		t.Errorf("expected block timestamp %v, got %v", blk.Timestamp, blkRes.Timestamp)
	}
	// gas limit
	if blk.GasLimit != blkRes.GasLimit {
		t.Errorf("expected block gas limit %v, got %v", blk.GasLimit, blkRes.GasLimit)
	}
	// gas used
	if blk.GasUsed != blkRes.GasUsed {
		t.Errorf("expected block gas used %v, got %v", blk.GasUsed, blkRes.GasUsed)
	}
	// transactions
	if len(blk.Transactions) != len(blkRes.Transactions) {
		t.Errorf("expected block transactions length %v, got %v", len(blk.Transactions), len(blkRes.Transactions))
	} else {
		for i, hash := range blk.Transactions {
			if hash != blkRes.Transactions[i] {
				t.Errorf("expected block transaction hash %v, got %v", hash, blkRes.Transactions[i])
			}
		}
	}
}

// validateTransaction validates that the given transaction equals the expected transaction.
func validateTransaction(t *testing.T, trx types.Transaction, trxRes types.Transaction) {
	t.Helper()
	// hash
	if trx.Hash != trxRes.Hash {
		t.Errorf("expected transaction hash %v, got %v", trx.Hash, trxRes.Hash)
	}
	// block hash
	if trx.BlockHash == nil && trxRes.BlockHash != nil {
		t.Errorf("expected nil block hash, got %v", trxRes.BlockHash)
	}
	if trx.BlockHash != nil && trxRes.BlockHash == nil {
		t.Errorf("expected block hash %v, got nil", *trx.BlockHash)
	}
	if trx.BlockHash != nil && trxRes.BlockHash != nil && *trx.BlockHash != *trxRes.BlockHash {
		t.Errorf("expected block hash %v, got %v", *trx.BlockHash, *trxRes.BlockHash)
	}
	// block number
	if trx.BlockNumber == nil && trxRes.BlockNumber != nil {
		t.Errorf("expected nil block number, got %v", trxRes.BlockNumber)
	}
	if trx.BlockNumber != nil && trxRes.BlockNumber == nil {
		t.Errorf("expected block number %v, got nil", *trx.BlockNumber)
	}
	if trx.BlockNumber != nil && trxRes.BlockNumber != nil && *trx.BlockNumber != *trxRes.BlockNumber {
		t.Errorf("expected block number %v, got %v", *trx.BlockNumber, *trxRes.BlockNumber)
	}
	// from
	if trx.From != trxRes.From {
		t.Errorf("expected from %v, got %v", trx.From, trxRes.From)
	}
	// to
	if trx.To == nil && trxRes.To != nil {
		t.Errorf("expected nil to, got %v", trxRes.To)
	}
	if trx.To != nil && trxRes.To == nil {
		t.Errorf("expected to %v, got nil", *trx.To)
	}
	if trx.To != nil && trxRes.To != nil && *trx.To != *trxRes.To {
		t.Errorf("expected to %v, got %v", *trx.To, *trxRes.To)
	}
	// contract address
	if trx.ContractAddress == nil && trxRes.ContractAddress != nil {
		t.Errorf("expected nil contract address, got %v", trxRes.ContractAddress)
	}
	if trx.ContractAddress != nil && trxRes.ContractAddress == nil {
		t.Errorf("expected contract address %v, got nil", *trx.ContractAddress)
	}
	if trx.ContractAddress != nil && trxRes.ContractAddress != nil && *trx.ContractAddress != *trxRes.ContractAddress {
		t.Errorf("expected contract address %v, got %v", *trx.ContractAddress, *trxRes.ContractAddress)
	}
	// nonce
	if trx.Nonce != trxRes.Nonce {
		t.Errorf("expected nonce %v, got %v", trx.Nonce, trxRes.Nonce)
	}
	// gas
	if trx.Gas != trxRes.Gas {
		t.Errorf("expected gas %v, got %v", trx.Gas, trxRes.Gas)
	}
	// gas used
	if trx.GasUsed == nil && trxRes.GasUsed != nil {
		t.Errorf("expected nil gas used, got %v", trxRes.GasUsed)
	}
	if trx.GasUsed != nil && trxRes.GasUsed == nil {
		t.Errorf("expected gas used %v, got nil", *trx.GasUsed)
	}
	if trx.GasUsed != nil && trxRes.GasUsed != nil && *trx.GasUsed != *trxRes.GasUsed {
		t.Errorf("expected gas used %v, got %v", *trx.GasUsed, *trxRes.GasUsed)
	}
	// cumulative gas used
	if trx.CumulativeGasUsed == nil && trxRes.CumulativeGasUsed != nil {
		t.Errorf("expected nil cumulative gas used, got %v", trxRes.CumulativeGasUsed)
	}
	if trx.CumulativeGasUsed != nil && trxRes.CumulativeGasUsed == nil {
		t.Errorf("expected cumulative gas used %v, got nil", *trx.CumulativeGasUsed)
	}
	if trx.CumulativeGasUsed != nil && trxRes.CumulativeGasUsed != nil && *trx.CumulativeGasUsed != *trxRes.CumulativeGasUsed {
		t.Errorf("expected cumulative gas used %v, got %v", *trx.CumulativeGasUsed, *trxRes.CumulativeGasUsed)
	}
	// gas price
	if trx.GasPrice.ToInt().Cmp(trxRes.GasPrice.ToInt()) != 0 {
		t.Errorf("expected gas price %v, got %v", trx.GasPrice.String(), trxRes.GasPrice.String())
	}
	// value
	if trx.Value.ToInt().Cmp(trxRes.Value.ToInt()) != 0 {
		t.Errorf("expected value %v, got %v", trx.Value.String(), trxRes.Value.String())
	}
	// input data
	if !bytes.Equal(trx.Input, trxRes.Input) {
		t.Errorf("expected input data %v, got %v", trx.Input, trxRes.Input)
	}
	// transaction index
	if trx.TransactionIndex == nil && trxRes.TransactionIndex != nil {
		t.Errorf("expected nil transaction index, got %v", trxRes.TransactionIndex)
	}
	if trx.TransactionIndex != nil && trxRes.TransactionIndex == nil {
		t.Errorf("expected transaction index %v, got nil", *trx.TransactionIndex)
	}
	if trx.TransactionIndex != nil && trxRes.TransactionIndex != nil && *trx.TransactionIndex != *trxRes.TransactionIndex {
		t.Errorf("expected transaction index %v, got %v", *trx.TransactionIndex, *trxRes.TransactionIndex)
	}
	// status
	if trx.Status == nil && trxRes.Status != nil {
		t.Errorf("expected nil status, got %v", trxRes.Status)
	}
	if trx.Status != nil && trxRes.Status == nil {
		t.Errorf("expected status %v, got nil", *trx.Status)
	}
	if trx.Status != nil && trxRes.Status != nil && *trx.Status != *trxRes.Status {
		t.Errorf("expected status %v, got %v", *trx.Status, *trxRes.Status)
	}
}
