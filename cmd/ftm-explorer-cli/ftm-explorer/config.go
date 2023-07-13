package ftm_explorer

import (
	"encoding/json"
	"ftm-explorer/cmd/ftm-explorer-cli/flags"
	"ftm-explorer/internal/config"

	"github.com/urfave/cli/v2"
)

var CmdConfig = cli.Command{
	Name:  "config",
	Usage: "Prints default config",
	Action: func(ctx *cli.Context) error {
		cfg := config.Load(ctx.String(flags.Cfg.Name))
		enc := json.NewEncoder(ctx.App.Writer)
		enc.SetIndent("", "    ")
		if err := enc.Encode(cfg); err != nil {
			return err
		}
		return nil
	},
}
