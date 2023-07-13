package config

import "github.com/op/go-logging"

// Config is the base configuration structure.
type Config struct {
	Rpc    Rpc
	Api    ApiServer
	Logger Logger
}

// Rpc is the configuration structure for RPC.
type Rpc struct {
	OperaRpcUrl string
}

// ApiServer is the configuration structure for API server.
type ApiServer struct {
	BindAddress     string
	DomainAddress   string
	ReadTimeout     int
	WriteTimeout    int
	IdleTimeout     int
	HeaderTimeout   int
	ResolverTimeout int
	CorsOrigin      []string
}

// Logger is the configuration structure for logging.
type Logger struct {
	LoggingLevel logging.Level
	LogFormat    string
}
