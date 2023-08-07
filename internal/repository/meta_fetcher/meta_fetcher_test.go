package meta_fetcher

import (
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test meta fetcher number of accounts
func TestMetaFetcher_NumberOfAccounts(t *testing.T) {
	server := startServer(t)

	mf := NewMetaFetcher(
		&config.MetaFetcher{NumberOfAccountsUrl: server.URL + "/number_of_accounts"}, logger.NewMockLogger())

	num, err := mf.NumberOfAccounts()
	if err != nil {
		t.Errorf("failed to get number of accounts: %v", err)
	}
	if num != 68 {
		t.Errorf("expected number of accounts to be 68, got %d", num)
	}
}

// Test meta fetcher time to finality
func TestMetaFetcher_TimeToFinality(t *testing.T) {
	server := startServer(t)

	mf := NewMetaFetcher(&config.MetaFetcher{TimeToFinalityUrl: server.URL + "/ttf"}, logger.NewMockLogger())

	num, err := mf.TimeToFinality()
	if err != nil {
		t.Errorf("failed to get time to finality: %v", err)
	}
	if num != 0.9009191708000002 {
		t.Errorf("expected time to finality to be 0.9009191708000002, got %f", num)
	}
}

// startServer starts a test server that will server metadata.
func startServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/number_of_accounts" {
			if _, err := w.Write([]byte("68")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
			return
		}
		if r.URL.Path == "/ttf" {
			if _, err := w.Write([]byte(`{"ttf":0.9009191708000002}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(func() {
		srv.Close()
	})
	return srv
}
