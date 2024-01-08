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

// Test meta fetcher disk size per 100M transactions
func TestMetaFetcher_DiskSizePer100MTxs(t *testing.T) {
	server := startServer(t)

	mf := NewMetaFetcher(
		&config.MetaFetcher{DiskSizePer100MTxsUrl: server.URL + "/disk_size_per_100m_txs"}, logger.NewMockLogger())

	num, err := mf.DiskSizePer100MTxs()
	if err != nil {
		t.Errorf("failed to get disk size per 100M transactions: %v", err)
	}
	if num != 72799695667 {
		t.Errorf("expected disk size per 100M transactions to be 72799695667, got %d", num)
	}
}

// Test meta fetcher disk size pruned per 100M transactions
func TestMetaFetcher_DiskSizePrunedPer100MTxs(t *testing.T) {
	server := startServer(t)

	mf := NewMetaFetcher(
		&config.MetaFetcher{DiskSizePer100MTxsUrl: server.URL + "/disk_size_pruned_per_100m_txs"}, logger.NewMockLogger())

	num, err := mf.DiskSizePer100MTxs()
	if err != nil {
		t.Errorf("failed to get disk size pruned per 100M transactions: %v", err)
	}
	if num != 62799695667 {
		t.Errorf("expected disk size pruned per 100M transactions to be 62799695667, got %d", num)
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

// Test meta fetcher idle status
func TestMetaFetcher_IsIdleStatus(t *testing.T) {
	server := startServer(t)

	mf := NewMetaFetcher(&config.MetaFetcher{IsIdleStatusUrl: server.URL + "/is_idle_status"}, logger.NewMockLogger())

	idle, err := mf.IsIdleStatus()
	if err != nil {
		t.Errorf("failed to get time to finality: %v", err)
	}
	if !idle {
		t.Errorf("expected is idle status to be true, got %t", idle)
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
		if r.URL.Path == "/disk_size_per_100m_txs" {
			if _, err := w.Write([]byte("72799695667")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
			return
		}
		if r.URL.Path == "/disk_size_pruned_per_100m_txs" {
			if _, err := w.Write([]byte("62799695667")); err != nil {
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
		if r.URL.Path == "/is_idle_status" {
			if _, err := w.Write([]byte("1")); err != nil {
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
