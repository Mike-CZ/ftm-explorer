package ftm_explorer

import (
	"fmt"
	"ftm-explorer/cmd/ftm-explorer-cli/flags"
	"ftm-explorer/internal/api"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/faucet"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/maze"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/repository/db"
	"ftm-explorer/internal/repository/meta_fetcher"
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/svc"

	"github.com/urfave/cli/v2"
)

// CmdRun defines a CLI command for running the gas ftm explorer.
var CmdRun = cli.Command{
	Action: run,
	Name:   "run",
	Usage:  `Runs the demonet explorer.`,
	Flags: []cli.Flag{
		&flags.Cfg,
	},
}

// run starts the fantom explorer.
func run(ctx *cli.Context) error {
	// load config
	cfg := config.Load(ctx.String(flags.Cfg.Name))

	// create logger
	log := logger.New(ctx.App.Writer, &cfg.Logger)

	// create repository
	repo, err := createRepository(cfg, log)
	if err != nil {
		return fmt.Errorf("can not create repository: %v", err)
	}

	// create services manager and run it
	mgr := svc.NewServiceManager(cfg, repo, log)
	mgr.Start()

	// create faucet
	fct, err := createFaucet(cfg, repo, log)
	if err != nil {
		return fmt.Errorf("can not create faucet: %v", err)
	}

	// create maze if enabled
	var m maze.IMaze
	if cfg.Maze != nil {
		m = maze.NewMaze(cfg.Maze)
	}

	// create api server
	apiServer := api.NewApiServer(cfg, repo, fct, m, log)

	// run api server
	apiServer.Start()

	return nil
}

// createRepository creates a new repository instance.
func createRepository(cfg *config.Config, log logger.ILogger) (*repository.Repository, error) {
	// create rpc connection
	operaRpc, err := rpc.NewOperaRpc(&cfg.Rpc)
	if err != nil {
		return nil, fmt.Errorf("can not create rpc connection: %v", err)
	}

	// create db connection
	database, err := db.NewMongoDb(&cfg.MongoDb, log)
	if err != nil {
		return nil, fmt.Errorf("can not create database connection: %v", err)
	}

	// create meta fetcher
	metaFetcher := meta_fetcher.NewMetaFetcher(&cfg.MetaFetcher, log)

	// create repository
	return repository.NewRepository(cfg.Explorer.BlockBufferSize, operaRpc, database, metaFetcher), nil
}

// createFaucet creates a new faucet instance.
func createFaucet(cfg *config.Config, repo *repository.Repository, log logger.ILogger) (*faucet.Faucet, error) {
	wallet, err := faucet.NewWallet(repo, log, cfg.Faucet.WalletPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("can not create faucet wallet: %v", err)
	}
	erc20s, err := faucet.NewFaucetErc20s(&cfg.Faucet, repo, log)
	if err != nil {
		return nil, fmt.Errorf("can not create faucet erc20s: %v", err)
	}
	f, err := faucet.NewFaucet(&cfg.Faucet, faucet.NewPhraseGenerator(), wallet, erc20s, repo)
	if err != nil {
		return nil, fmt.Errorf("can not create faucet: %v", err)
	}
	return f, nil
}
