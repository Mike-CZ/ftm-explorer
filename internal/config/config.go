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
	Maze        *Maze
}

// Explorer is the configuration structure for the explorer.
type Explorer struct {
	// BlockBufferSize is the size of the block buffer. The buffer is used to
	// store blocks in memory, so that they can be accessed quickly.
	BlockBufferSize uint
	// IsPersisted is the flag indicating whether the explorer data is persisted.
	IsPersisted bool
	// MaxTxsCount is the maximum number of transactions to be stored in the database.
	// If the number of transactions in the database exceeds this value, the oldest
	// transactions are removed.
	MaxTxsCount uint
}

type Faucet struct {
	// ClaimLimitSeconds is the time limit between two claims.
	ClaimLimitSeconds uint
	// ClaimTokensAmount is the amount of tokens to be claimed.
	ClaimTokensAmount float32
	// WalletPrivateKey is the private key of the faucet wallet.
	WalletPrivateKey string
	// ClaimsPerDay is the number of claims per day allowed from the same ip address.
	ClaimsPerDay uint
	// Erc20sPath is the path to the erc20 tokens configuration file.
	Erc20sPath string
	// Erc20MintAmountHex is the amount of erc20 tokens to be minted.
	Erc20MintAmountHex string
	// Erc20s is the list of erc20 tokens to be claimed.
	Erc20s []FaucetErc20
}

// FaucetErc20 is the configuration structure for the faucet erc20 token.
type FaucetErc20 struct {
	Address  string `json:"address"`
	MinterPk string `json:"minter_key"`
}

// MetaFetcher is the configuration structure for meta fetcher obtaining blockchain metadata.
type MetaFetcher struct {
	NumberOfAccountsUrl         string
	DiskSizePer100MTxsUrl       string
	DiskSizePrunedPer100MTxsUrl string
	TimeToFinalityUrl           string
	IsIdleStatusUrl             string
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

// Maze is the configuration structure for the maze.
type Maze struct {
	VisibilityRange uint
	ConfigPaths     []string
	Configs         []MazeConfig
}

// MazeConfig is the configuration structure for the maze.
type MazeConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Width   int32  `json:"width"`
	Height  int32  `json:"height"`
	Entry   int32  `json:"entry"`
	Exit    int32  `json:"exit"`
	Tiles   []struct {
		Id       int32 `json:"id"`
		Position struct {
			X int32 `json:"x"`
			Y int32 `json:"y"`
		} `json:"position"`
		Paths struct {
			North *int32 `json:"north"`
			East  *int32 `json:"east"`
			South *int32 `json:"south"`
			West  *int32 `json:"west"`
		} `json:"paths"`
	} `json:"tiles"`
}
