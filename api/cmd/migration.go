package cmd

import (
	"context"
	"os"

	"github.com/Ayut0/coffee-journal/api/config"
	"github.com/Ayut0/coffee-journal/api/migration"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var MigrateCommand = &cli.Command{
	Name: "migrate",
	Usage: "Run database migrations",
	Commands: []*cli.Command{
		{
			Name: "up",
			Usage: "Run all pending migrations",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

				cfg, err := config.Load()
				if err != nil {
					return err
				}

				log.Info().Msg("running migrations up...")
				if err := migration.Up(cfg.DatabaseURL); err !=nil {
					return err
				}
				log.Info().Msg("migrations up completed successfully")
				return nil
			},
		},
		{
			Name: "down",
			Usage: "Roll back the last migration",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name: "steps",
					Usage: "Number of migrations to roll back",
					Value: 1,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				log.Logger= log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
				cfg, err := config.Load()
				if err != nil {
					return err
				}

				steps := int(cmd.Int("steps"))
				log.Info().Int("steps", steps).Msg("rolling back migrations...")
				if err := migration.Down(cfg.DatabaseURL, steps); err != nil {
					return err
				}
				log.Info().Int("steps", steps).Msg("migrations down completed successfully")
				return nil
			},
		},
	},
}