package api

import (
	"encoding/json"
	"fmt"
	"ftm-explorer/internal/api/graphql/resolvers"
	"ftm-explorer/internal/api/handlers"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

// apiTestCase represents a test case for the API server.
type apiTestCase struct {
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
	}
}

// getTransactionTestCase returns a test case for a transaction not found error.
func getTransactionTestCase(t *testing.T) apiTestCase {
	trx := getTestTransaction(t)
	return apiTestCase{
		requestBody: fmt.Sprintf(`{"query": "query { transaction(hash: \"%s\") { hash, blockHash, blockNumber, from, to, contractAddress, nonce, gas, gasUsed, cumulativeGasUsed, gasPrice, value, inputData, trxIndex, status }}"}`, trx.Hash.Hex()),
		buildStubs: func(mockRepository *repository.MockRepository) {
			mockRepository.EXPECT().GetTransactionByHash(gomock.Any()).Return(&trx, nil)
		},
		checkResponse: func(t *testing.T, resp *http.Response) {
			apiRes := decodeResponse(t, resp)
			if len(apiRes.Errors) != 0 {
				t.Errorf("expected no errors, got %v", len(apiRes.Errors))
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
	return types.Transaction{
		Hash:      common.HexToHash("0xcd38ab7c7c77700e3d225316f10b1e5cbdff4c294a7034142e5f4a0ef7eebd7f"),
		BlockHash: &common.Hash{common.FromHex("")},
	}
}

// validateTransaction validates that the given transaction equals the expected transaction.
func validateTransaction(t *testing.T, trx types.Transaction, trxRes types.Transaction) {
	t.Helper()
	if trx.Hash != trxRes.Hash {
		t.Errorf("expected transaction hash %v, got %v", trx.Hash, trxRes.Hash)
	}
}
