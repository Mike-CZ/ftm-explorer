package config

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Load loads the configuration from the config file and command line flags.
func Load(path string) *Config {
	var config Config

	cfg, err := readConfigFile(path)
	if err != nil {
		log.Fatalf("can not read configuration file. Err: %v", err)
	}

	if err = cfg.Unmarshal(&config); err != nil {
		log.Fatalf("can not extract configuration. Err: %v", err)
	}

	return &config
}

// readConfigFile reads the configuration file from the given path.
func readConfigFile(path string) (*viper.Viper, error) {
	cfg := viper.New()
	cfg.SetConfigName("ftm_explorer")
	cfg.AddConfigPath("/etc")
	cfg.AddConfigPath("/etc/ftm_explorer")

	if path != "" {
		log.Println("loading config: ", path)
		cfg.SetConfigFile(path)
	}
	applyDefaults(cfg)

	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("configuration not found at %s", cfg.ConfigFileUsed())
			return nil, err
		}

		// config file not found; ignore the error, we may not need the config file
		log.Println("configuration file not found, using default values")
	}

	// load erc20s
	erc20sPath := cfg.GetString("faucet.erc20sPath")
	if erc20sPath == "" {
		log.Println("erc20sPath is empty")
	} else {
		log.Println("loading erc20s from: ", erc20sPath)
		erc20s, err := readErc20sFile(erc20sPath)
		if err != nil {
			log.Printf("can not read erc20s file. Err: %v", err)
			return nil, err
		}
		cfg.Set("faucet.erc20s", erc20s)
	}

	return cfg, nil
}

// readErc20sFile reads the erc20s file from the given path.
func readErc20sFile(path string) ([]FaucetErc20, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Printf("can not open erc20s file. Err: %v", err)
		return nil, err
	}
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Printf("can not read erc20s file. Err: %v", err)
		return nil, err
	}

	var erc20s []FaucetErc20
	if err := json.Unmarshal(byteValue, &erc20s); err != nil {
		log.Printf("can not unmarshal erc20s file. Err: %v", err)
		return nil, err
	}

	return erc20s, nil
}
