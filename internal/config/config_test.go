package config

import (
	"fmt"
	"os"
	"testing"
)

// Test that config is loaded from file.
func TestConfig_Load(t *testing.T) {
	cfgStr := `{
	  "explorer": {
		"blockBufferSize": 128964,
		"isPersisted": true,
		"maxTxsCount": 66999999
	  },
      "faucet": {
        "claimLimitSeconds": 1000,
        "claimTokensAmount": 0.5,
        "walletPrivateKey": "9s4d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58a6285",
        "claimsPerDay": 5,
	    "erc20MintAmountHex": "0x38d7ea4c68000",
        "erc20sPath": "%s"
      },
	  "maze": {
 		"visibilityRange": 3,
	    "configPaths": ["%s"]
	  },
	  "metaFetcher": {
		"numberOfAccountsUrl": "number-of-accounts-test-url",
		"diskSizePer100MTxsUrl": "disk-size-test-url",
		"diskSizePrunedPer100MTxsUrl": "disk-size-pruned-test-url",
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
	erc20CfgStr := `[
	  {
		"name":"Apatite",
		"address":"0x3bc666c4073853a59a7bfb0184298551d922f1df",
		"minter_key":"904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285"
	  },
	  {
		"name":"Epidote",
		"address":"0x1234567890123456789012345678901234567890",
		"minter_key":"904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285"
	  }
	]`
	mazeCfgStr := `{
	  "name": "Black Squirrel",
	  "address": "0x0000000000000000000000000000000000000000",
	  "width": 3,
	  "height": 3,
	  "entry": 0,
	  "exit": 6,
	  "tiles": [
		{
		  "id": 1941,
		  "position": {
			"x": 0,
			"y": 0
		  },
		  "paths": {
			"north": null,
			"south": 3,
			"east": null,
			"west": null
		  }
		},
		{
		  "id": 808,
		  "position": {
			"x": 1,
			"y": 0
		  },
		  "paths": {
			"north": null,
			"south": null,
			"east": 2,
			"west": null
		  }
		}
	  ]
	}`

	// store config into temporary file
	file, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	erc20File, err := os.CreateTemp("", "erc20s*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(erc20File.Name())

	mazeFile, err := os.CreateTemp("", "maze*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(mazeFile.Name())

	// Write the configuration string to the temporary file
	_, err = file.Write([]byte(fmt.Sprintf(cfgStr, erc20File.Name(), mazeFile.Name())))
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	// Write erc20 config to the temporary file
	_, err = erc20File.Write([]byte(erc20CfgStr))
	if err != nil {
		t.Fatal(err)
	}
	if err := erc20File.Close(); err != nil {
		t.Fatal(err)
	}

	// Write maze config to the temporary file
	_, err = mazeFile.Write([]byte(mazeCfgStr))
	if err != nil {
		t.Fatal(err)
	}
	if err := mazeFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the configuration from the temporary file
	cfg := Load(file.Name())

	// Check the configuration values
	if cfg.Explorer.BlockBufferSize != 128964 {
		t.Errorf("expected Explorer.BlockBufferSize to be 128964, got %d", cfg.Explorer.BlockBufferSize)
	}
	if !cfg.Explorer.IsPersisted {
		t.Errorf("expected Explorer.IsPersisted to be true, got %v", cfg.Explorer.IsPersisted)
	}
	if cfg.Explorer.MaxTxsCount != 66999999 {
		t.Errorf("expected Explorer.MaxTxsCount to be 66999999, got %d", cfg.Explorer.MaxTxsCount)
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
	if cfg.Faucet.ClaimsPerDay != 5 {
		t.Errorf("expected Faucet.ClaimsPerDay to be 5, got %d", cfg.Faucet.ClaimsPerDay)
	}
	if cfg.Faucet.Erc20sPath != erc20File.Name() {
		t.Errorf("expected Faucet.Erc20sPath to be %s, got %s", erc20File.Name(), cfg.Faucet.Erc20sPath)
	}
	if cfg.Faucet.Erc20MintAmountHex != "0x38d7ea4c68000" {
		t.Errorf("expected Faucet.Erc20MintAmountHex to be 0x38d7ea4c68000, got %s", cfg.Faucet.Erc20MintAmountHex)
	}
	// check erc20 config
	if len(cfg.Faucet.Erc20s) != 2 {
		t.Errorf("expected Faucet.Erc20s to have 2 elements, got %d", len(cfg.Faucet.Erc20s))
	}
	if cfg.Faucet.Erc20s[0].Address != "0x3bc666c4073853a59a7bfb0184298551d922f1df" {
		t.Errorf("expected Faucet.Erc20s[0].Address to be 0x3bc666c4073853a59a7bfb0184298551d922f1df, got %s", cfg.Faucet.Erc20s[0].Address)
	}
	if cfg.Faucet.Erc20s[0].MinterPk != "904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285" {
		t.Errorf("expected Faucet.Erc20s[0].MinterPk to be 904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285, got %s", cfg.Faucet.Erc20s[0].MinterPk)
	}
	if cfg.Faucet.Erc20s[1].Address != "0x1234567890123456789012345678901234567890" {
		t.Errorf("expected Faucet.Erc20s[1].Address to be 0x1234567890123456789012345678901234567890, got %s", cfg.Faucet.Erc20s[1].Address)
	}
	if cfg.Faucet.Erc20s[1].MinterPk != "904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285" {
		t.Errorf("expected Faucet.Erc20s[1].MinterPk to be 904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285, got %s", cfg.Faucet.Erc20s[1].MinterPk)
	}
	if cfg.MetaFetcher.NumberOfAccountsUrl != "number-of-accounts-test-url" {
		t.Errorf("expected MetaFetcher.NumberOfAccountsUrl to be number-of-accounts-test-url, got %s", cfg.MetaFetcher.NumberOfAccountsUrl)
	}
	if cfg.MetaFetcher.DiskSizePer100MTxsUrl != "disk-size-test-url" {
		t.Errorf("expected MetaFetcher.DiskSizePer100MTxsUrl to be disk-size-test-url, got %s", cfg.MetaFetcher.DiskSizePer100MTxsUrl)
	}
	if cfg.MetaFetcher.DiskSizePrunedPer100MTxsUrl != "disk-size-pruned-test-url" {
		t.Errorf("expected MetaFetcher.DiskSizePrunedPer100MTxsUrl to be disk-size-pruned-test-url, got %s", cfg.MetaFetcher.DiskSizePrunedPer100MTxsUrl)
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
	// check maze config
	if cfg.Maze == nil {
		t.Errorf("expected Maze to be not nil")
	}
	if cfg.Maze.VisibilityRange != 3 {
		t.Errorf("expected Maze.VisibilityRange to be 3, got %d", cfg.Maze.VisibilityRange)
	}
	if len(cfg.Maze.ConfigPaths) != 1 || cfg.Maze.ConfigPaths[0] != mazeFile.Name() {
		t.Errorf("expected Maze.ConfigPaths to be [%s], got %v", mazeFile.Name(), cfg.Maze.ConfigPaths)
	}
	if len(cfg.Maze.Configs) != 1 {
		t.Errorf("expected Maze.Configs to have 1 element, got %d", len(cfg.Maze.Configs))
	}
	if cfg.Maze.Configs[0].Name != "Black Squirrel" {
		t.Errorf("expected Maze.Configs[0].Name to be Black Squirrel, got %s", cfg.Maze.Configs[0].Name)
	}
	if cfg.Maze.Configs[0].Address != "0x0000000000000000000000000000000000000000" {
		t.Errorf("expected Maze.Configs[0].Address to be \"0x0000000000000000000000000000000000000000\", got %s", cfg.Maze.Configs[0].Address)
	}
	if cfg.Maze.Configs[0].Width != 3 {
		t.Errorf("expected Maze.Configs[0].Width to be 3, got %d", cfg.Maze.Configs[0].Width)
	}
	if cfg.Maze.Configs[0].Height != 3 {
		t.Errorf("expected Maze.Configs[0].Height to be 3, got %d", cfg.Maze.Configs[0].Height)
	}
	if cfg.Maze.Configs[0].Entry != 0 {
		t.Errorf("expected Maze.Configs[0].Entry to be 0, got %d", cfg.Maze.Configs[0].Entry)
	}
	if cfg.Maze.Configs[0].Exit != 6 {
		t.Errorf("expected Maze.Configs[0].Exit to be 6, got %d", cfg.Maze.Configs[0].Exit)
	}
	if len(cfg.Maze.Configs[0].Tiles) != 2 {
		t.Errorf("expected Maze.Configs[0].Tiles to have 2 elements, got %d", len(cfg.Maze.Configs[0].Tiles))
	}
	if cfg.Maze.Configs[0].Tiles[0].Id != 1941 {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Id to be 1941, got %d", cfg.Maze.Configs[0].Tiles[0].Id)
	}
	if cfg.Maze.Configs[0].Tiles[0].Position.X != 0 {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Position.X to be 0, got %d", cfg.Maze.Configs[0].Tiles[0].Position.X)
	}
	if cfg.Maze.Configs[0].Tiles[0].Position.Y != 0 {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Position.Y to be 0, got %d", cfg.Maze.Configs[0].Tiles[0].Position.Y)
	}
	if cfg.Maze.Configs[0].Tiles[0].Paths.North != nil {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Paths.North to be nil, got %d", *cfg.Maze.Configs[0].Tiles[0].Paths.North)
	}
	if cfg.Maze.Configs[0].Tiles[0].Paths.South == nil {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Paths.South to be not nil")
	}
	if *cfg.Maze.Configs[0].Tiles[0].Paths.South != 3 {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Paths.South to be 3, got %d", *cfg.Maze.Configs[0].Tiles[0].Paths.South)
	}
	if cfg.Maze.Configs[0].Tiles[0].Paths.East != nil {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Paths.East to be nil, got %d", *cfg.Maze.Configs[0].Tiles[0].Paths.East)
	}
	if cfg.Maze.Configs[0].Tiles[0].Paths.West != nil {
		t.Errorf("expected Maze.Configs[0].Tiles[0].Paths.West to be nil, got %d", *cfg.Maze.Configs[0].Tiles[0].Paths.West)
	}
	if cfg.Maze.Configs[0].Tiles[1].Id != 808 {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Id to be 808, got %d", cfg.Maze.Configs[0].Tiles[1].Id)
	}
	if cfg.Maze.Configs[0].Tiles[1].Position.X != 1 {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Position.X to be 1, got %d", cfg.Maze.Configs[0].Tiles[1].Position.X)
	}
	if cfg.Maze.Configs[0].Tiles[1].Position.Y != 0 {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Position.Y to be 0, got %d", cfg.Maze.Configs[0].Tiles[1].Position.Y)
	}
	if cfg.Maze.Configs[0].Tiles[1].Paths.North != nil {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Paths.North to be nil, got %d", *cfg.Maze.Configs[0].Tiles[1].Paths.North)
	}
	if cfg.Maze.Configs[0].Tiles[1].Paths.South != nil {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Paths.South to be nil, got %d", *cfg.Maze.Configs[0].Tiles[1].Paths.South)
	}
	if cfg.Maze.Configs[0].Tiles[1].Paths.East == nil {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Paths.East to be not nil")
	}
	if *cfg.Maze.Configs[0].Tiles[1].Paths.East != 2 {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Paths.East to be 2, got %d", *cfg.Maze.Configs[0].Tiles[1].Paths.East)
	}
	if cfg.Maze.Configs[0].Tiles[1].Paths.West != nil {
		t.Errorf("expected Maze.Configs[0].Tiles[1].Paths.West to be nil, got %d", *cfg.Maze.Configs[0].Tiles[1].Paths.West)
	}
}
