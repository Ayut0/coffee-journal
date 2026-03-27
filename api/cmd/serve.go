package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ayut0/coffee-journal/api/builder"
	"github.com/Ayut0/coffee-journal/api/config"
	server "github.com/Ayut0/coffee-journal/api/internal/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

// ServeCommand starts the HTTP server.
var ServeCommand = &cli.Command{
	Name:  "serve",
	Usage: "Start the HTTP API server",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		d := &builder.Dependency{Cfg: cfg}
		e := server.NewServer(d)

		// Signal context — tells us when to begin shutdown
		sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		// Start server in background goroutine
		go func() {
			if err := e.Start(":" + cfg.Port); !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("server error")
				stop() // trigger shutdown on unexpected error
			}
		}()

		// Block until signal received
		<-sigCtx.Done()
		log.Info().Msg("shutting down server...")

		// Shutdown context — gives in-flight requests 10s to finish
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(shutCtx); err != nil {
			log.Error().Err(err).Msg("shutdown error")
			return err
		}

		return nil
	},
}
