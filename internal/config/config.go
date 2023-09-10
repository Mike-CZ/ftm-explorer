package config

import (
	"github.com/op/go-logging"
)

// Config is the base configuration structure.
type Config struct {
	Explorer    Explorer
	Faucet      Faucet
	MetaFetcher MetaFetcher
	Rpc         Rpc
	Api         ApiServer
	Logger      Logger
	MongoDb     MongoDb
}

// Explorer is the configuration structure for the explorer.
type Explorer struct {
	// BlockBufferSize is the size of the block buffer. The buffer is used to
	// store blocks in memory, so that they can be accessed quickly.
	BlockBufferSize uint
}

type Faucet struct {
	// ClaimLimitSeconds is the time limit between two claims.
	ClaimLimitSeconds uint
	// ClaimTokensAmount is the amount of tokens to be claimed.
	ClaimTokensAmount float32
	// WalletPrivateKey is the private key of the faucet wallet.
	WalletPrivateKey string
}

// MetaFetcher is the configuration structure for meta fetcher obtaining blockchain metadata.
type MetaFetcher struct {
	NumberOfAccountsUrl   string
	DiskSizePer100MTxsUrl string
	TimeToFinalityUrl     string
}

// Rpc is the configuration structure for RPC.
type Rpc struct {
	OperaRpcUrl string
	SfcAddress  string
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

// MongoDb is the configuration structure for MongoDB.
type MongoDb struct {
	Host     string
	Port     int
	Db       string
	User     *string
	Password *string
}
