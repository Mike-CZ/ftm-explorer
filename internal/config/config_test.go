package config

import (
	"os"
	"testing"
)

// Test that config is loaded from file.
func TestConfig_Load(t *testing.T) {
	cfgStr := `{
	  "explorer": {
		"blockBufferSize": 128964
	  },
	  "metaFetcher": {
		"url": "metafetcher-test-url"
	  },
	  "rpc": {
		"operaRpcUrl": "opera-rpc"
	  },
	  "api": {
		"readTimeout": 200,
		"writeTimeout": 100,
		"idleTimeout": 58,
		"headerTimeout": 46,
		"resolverTimeout": 12,
		"bindAddress": "bindAddress",
		"domainAddress": "domainAddress",
		"corsOrigin": ["cors1", "cors2"]
	  },
	  "logger": {
		"loggingLevel": 1,
		"logFormat": "some-log-format"
	  },
	  "mongodb": {
		"host": "mongohost",
		"port": 1111,
		"db": "mongodb"
	  }
	}`
	// store config into temporary file
	file, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write the configuration string to the temporary file
	_, err = file.Write([]byte(cfgStr))
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the configuration from the temporary file
	cfg := Load(file.Name())

	// Check the configuration values
	if cfg.Explorer.BlockBufferSize != 128964 {
		t.Errorf("expected Explorer.BlockBufferSize to be 128964, got %d", cfg.Explorer.BlockBufferSize)
	}
	if cfg.MetaFetcher.Url != "metafetcher-test-url" {
		t.Errorf("expected MetaFetcher.Url to be metafetcher-test-url, got %s", cfg.MetaFetcher.Url)
	}
	if cfg.Rpc.OperaRpcUrl != "opera-rpc" {
		t.Errorf("expected RPC.OperaRpcUrl to be opera-rpc, got %s", cfg.Rpc.OperaRpcUrl)
	}
	if cfg.Api.ReadTimeout != 200 {
		t.Errorf("expected Api.ReadTimeout to be 200, got %d", cfg.Api.ReadTimeout)
	}
	if cfg.Api.WriteTimeout != 100 {
		t.Errorf("expected Api.WriteTimeout to be 100, got %d", cfg.Api.WriteTimeout)
	}
	if cfg.Api.IdleTimeout != 58 {
		t.Errorf("expected Api.IdleTimeout to be 58, got %d", cfg.Api.IdleTimeout)
	}
	if cfg.Api.HeaderTimeout != 46 {
		t.Errorf("expected Api.HeaderTimeout to be 46, got %d", cfg.Api.HeaderTimeout)
	}
	if cfg.Api.ResolverTimeout != 12 {
		t.Errorf("expected Api.ResolverTimeout to be 12, got %d", cfg.Api.ResolverTimeout)
	}
	if cfg.Api.BindAddress != "bindAddress" {
		t.Errorf("expected Api.BindAddress to be bindAddress, got %s", cfg.Api.BindAddress)
	}
	if cfg.Api.DomainAddress != "domainAddress" {
		t.Errorf("expected Api.DomainAddress to be domainAddress, got %s", cfg.Api.DomainAddress)
	}
	if len(cfg.Api.CorsOrigin) != 2 || cfg.Api.CorsOrigin[0] != "cors1" || cfg.Api.CorsOrigin[1] != "cors2" {
		t.Errorf("expected Api.CorsOrigin to be [cors1, cors2], got %v", cfg.Api.CorsOrigin)
	}
	if cfg.Logger.LoggingLevel != 1 {
		t.Errorf("expected Logger.LoggingLevel to be 1, got %d", cfg.Logger.LoggingLevel)
	}
	if cfg.Logger.LogFormat != "some-log-format" {
		t.Errorf("expected Logger.LogFormat to be some-log-format, got %s", cfg.Logger.LogFormat)
	}
	if cfg.MongoDb.Host != "mongohost" {
		t.Errorf("expected Mongodb.Host to be mongohost, got %s", cfg.MongoDb.Host)
	}
	if cfg.MongoDb.Port != 1111 {
		t.Errorf("expected Mongodb.Port to be 1111, got %d", cfg.MongoDb.Port)
	}
	if cfg.MongoDb.Db != "mongodb" {
		t.Errorf("expected Mongodb.Db to be mongodb, got %s", cfg.MongoDb.Db)
	}
}
