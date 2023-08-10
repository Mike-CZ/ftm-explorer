package ftm_explorer

import (
	"fmt"
	"ftm-explorer/cmd/ftm-explorer-cli/flags"
	"ftm-explorer/internal/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urfave/cli/v2"
)

// CmdGenerateJwt represents the command used to generate a new JWT token.
var CmdGenerateJwt = cli.Command{
	Action:    generateJwt,
	Name:      "generate-jwt",
	Usage:     "Generate a JWT with the given expiration time in minutes",
	ArgsUsage: "[expiration]",
	Flags: []cli.Flag{
		&flags.Cfg,
	},
}

// run starts the fantom explorer.
func generateJwt(ctx *cli.Context) error {
	// parse expiration
	expiration := ctx.Args().Get(0)
	if expiration == "" {
		return fmt.Errorf("expiration is required")
	}
	expDuration, err := strconv.ParseInt(expiration, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid expiration value: %v", err)
	}

	// load config
	cfg := config.Load(ctx.String(flags.Cfg.Name))

	// create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Minute * time.Duration(expDuration)).Unix(), // Token expiration time
		"version": cfg.Api.Jwt.Version,
	})

	// sign and get the complete encoded token as a string using a secret
	tokenString, err := token.SignedString([]byte(cfg.Api.Jwt.Secret))
	if err != nil {
		return err
	}

	// print the token
	fmt.Println(tokenString)

	return nil
}
