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
      "faucet": {
        "claimLimitSeconds": 1000,
        "claimTokensAmount": 0.5,
        "walletPrivateKey": "9s4d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58a6285"
      },
	  "metaFetcher": {
		"numberOfAccountsUrl": "number-of-accounts-test-url",
		"diskSizePer100MTxsUrl": "disk-size-test-url",
		"timeToFinalityUrl": "time-to-finality-test-url"
	  },
	  "rpc": {
		"operaRpcUrl": "opera-rpc",
		"sfcAddress": "0x1234567890123456789012345678901234567890"
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
		"db": "mongodb",
		"user": "testUser",
		"password": "testPassword"
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
	if cfg.Faucet.ClaimLimitSeconds != 1000 {
		t.Errorf("expected Faucet.ClaimLimitSeconds to be 1000, got %d", cfg.Faucet.ClaimLimitSeconds)
	}
	if cfg.Faucet.ClaimTokensAmount != 0.5 {
		t.Errorf("expected Faucet.ClaimTokensAmount to be 0.5, got %f", cfg.Faucet.ClaimTokensAmount)
	}
	if cfg.Faucet.WalletPrivateKey != "9s4d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58a6285" {
		t.Errorf("expected Faucet.WalletPrivateKey to be 9s4d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58a6285, got %s", cfg.Faucet.WalletPrivateKey)
	}
	if cfg.MetaFetcher.NumberOfAccountsUrl != "number-of-accounts-test-url" {
		t.Errorf("expected MetaFetcher.NumberOfAccountsUrl to be number-of-accounts-test-url, got %s", cfg.MetaFetcher.NumberOfAccountsUrl)
	}
	if cfg.MetaFetcher.DiskSizePer100MTxsUrl != "disk-size-test-url" {
		t.Errorf("expected MetaFetcher.DiskSizePer100MTxsUrl to be disk-size-test-url, got %s", cfg.MetaFetcher.DiskSizePer100MTxsUrl)
	}
	if cfg.MetaFetcher.TimeToFinalityUrl != "time-to-finality-test-url" {
		t.Errorf("expected MetaFetcher.TimeToFinalityUrl to be time-to-finality-test-url, got %s", cfg.MetaFetcher.TimeToFinalityUrl)
	}
	if cfg.Rpc.OperaRpcUrl != "opera-rpc" {
		t.Errorf("expected RPC.OperaRpcUrl to be opera-rpc, got %s", cfg.Rpc.OperaRpcUrl)
	}
	if cfg.Rpc.SfcAddress != "0x1234567890123456789012345678901234567890" {
		t.Errorf("expected RPC.SfcAddress to be 0x1234567890123456789012345678901234567890, got %s", cfg.Rpc.SfcAddress)
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
	if cfg.MongoDb.User == nil {
		t.Errorf("expected Mongodb.User to be not nil")
	}
	if *cfg.MongoDb.User != "testUser" {
		t.Errorf("expected Mongodb.User to be testUser, got %s", *cfg.MongoDb.User)
	}
	if cfg.MongoDb.Password == nil {
		t.Errorf("expected Mongodb.Password to be not nil")
	}
	if *cfg.MongoDb.Password != "testPassword" {
		t.Errorf("expected Mongodb.Password to be testPassword, got %s", *cfg.MongoDb.Password)
	}
}
