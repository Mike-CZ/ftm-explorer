package config

import (
	"log"

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

	return cfg, nil
}
