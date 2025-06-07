package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


func initLog(w io.Writer, debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(w).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	_ = args // Unused args, can be used for command line arguments

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	initLog(w, *debug)

	go func() {
		log.Info().Msg("Starting a new echo service")
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
		defer cancel()
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
