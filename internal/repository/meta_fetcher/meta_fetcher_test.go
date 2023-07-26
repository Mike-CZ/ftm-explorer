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

	mf := NewMetaFetcher(&config.MetaFetcher{Url: server.URL}, logger.NewMockLogger())

	num, err := mf.NumberOfAccounts()
	if err != nil {
		t.Errorf("failed to get number of accounts: %v", err)
	}
	if num != 68 {
		t.Errorf("expected number of accounts to be 68, got %d", num)
	}
}

// startServer starts a test server that will server metadata.
func startServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if _, err := w.Write([]byte("68")); err != nil {
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
