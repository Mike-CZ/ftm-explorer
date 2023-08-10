package main

import (
	"ftm-explorer/cmd/ftm-explorer-cli/ftm-explorer"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func initApp() *cli.App {
	return &cli.App{
		Name:     "Demonet Explorer",
		HelpName: "demonet-explorer",
		Usage:    "starts observing blocks and aggregating data, serving it via API.",
		Commands: []*cli.Command{
			&ftm_explorer.CmdRun,
			&ftm_explorer.CmdConfig,
			&ftm_explorer.CmdGenerateJwt,
		},
	}
}

func main() {
	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
