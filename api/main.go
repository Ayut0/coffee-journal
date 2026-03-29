package main

import (
	"context"
	"os"

	"github.com/Ayut0/coffee-journal/api/cmd"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "api",
		Usage: "Coffee Journal API server",
		Commands: []*cli.Command{
			cmd.ServeCommand,
			cmd.MigrateCommand,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Send()
	}
}
