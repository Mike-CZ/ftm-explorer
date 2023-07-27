package ftm_explorer

import (
	"ftm-explorer/cmd/ftm-explorer-cli/flags"
	"ftm-explorer/internal/api"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/repository/db"
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
		return err
	}

	// create services manager and run it
	mgr := svc.NewServiceManager(cfg, repo, log)
	mgr.Start()

	// create api server
	apiServer := api.NewApiServer(&cfg.Api, repo, log)

	// run api server
	apiServer.Start()

	return nil
}

// createRepository creates a new repository instance.
func createRepository(cfg *config.Config, log logger.ILogger) (*repository.Repository, error) {
	// create rpc connection
	operaRpc, err := rpc.NewOperaRpc(&cfg.Rpc)
	if err != nil {
		return nil, err
	}

	// create db connection
	database, err := db.NewMongoDb(&cfg.MongoDb, log)
	if err != nil {
		return nil, err
	}

	// create repository
	return repository.NewRepository(cfg.Explorer.BlockBufferSize, operaRpc, database), nil
}
