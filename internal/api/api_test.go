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

func TestApiServer_Run(t *testing.T) {
	// initialize stubs
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockLogger := logger.NewMockLogger()

	// initialize test server
	handler := handlers.ApiHandler([]string{"*"}, resolvers.NewResolver(mockRepository, mockLogger), mockLogger)
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})

	// use table-driven testing to test multiple cases
	testCases := []apiTestCase{
		getTransactionTestCase(t),
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
			defer resp.Body.Close()
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
	if trx.From == trxRes.From {
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
