package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Ow1Dev/Zynra/internal/config"
	"github.com/Ow1Dev/Zynra/internal/server"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	_ = args // Unused args, can be used for command line arguments

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

  zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	log.Logger = zerolog.New(w).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	config := config.Config{
		Host: "0.0.0.0",
		Port: "8080",
	}

	srv := server.NewServer()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	go func() {
		log.Info().Msgf("listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
