package config

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// applyDefaults applies default values to the configuration.
func applyDefaults(cfg *viper.Viper) {
	// explorer
	cfg.SetDefault("explorer.blockBufferSize", 10_000)

	// rpc
	cfg.SetDefault("rpc.operaRpcUrl", "https://rpcapi.fantom.network")

	// apiserver server
	cfg.SetDefault("api.readTimeout", 2)
	cfg.SetDefault("api.writeTimeout", 15)
	cfg.SetDefault("api.idleTimeout", 1)
	cfg.SetDefault("api.headerTimeout", 1)
	cfg.SetDefault("api.resolverTimeout", 30)
	cfg.SetDefault("api.bindAddress", "localhost:16761")
	cfg.SetDefault("api.domainAddress", "localhost:16761")
	cfg.SetDefault("api.corsOrigin", []string{"*"})

	// logger
	cfg.SetDefault("logger.loggingLevel", logging.INFO)
	cfg.SetDefault("logger.logFormat", "%{color}%{level:-8s} %{shortpkg}/%{shortfunc}%{color:reset}: %{message}")

	// mongodb
	cfg.SetDefault("mongodb.host", "localhost")
	cfg.SetDefault("mongodb.port", 27017)
	cfg.SetDefault("mongodb.database", "ftm-explorer")
}
