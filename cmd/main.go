package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Ow1Dev/Zynra/internal/config"
	"github.com/Ow1Dev/Zynra/internal/repository"
	"github.com/Ow1Dev/Zynra/internal/server"
	"github.com/Ow1Dev/Zynra/pkgs/httpsutils"
)

// TODO: Make this to a package
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

	config := config.Config{
		Host: "0.0.0.0",
	}

	repo := repository.NewServiceRepository()

	job := repository.NewCleanupJob(repo)
	job.Start(ctx, 5*time.Second)

	httpRouterServer := httpsutils.NewHTTPServer(
		server.NewRouterServer(repo), 
		"8080",
		config)
	httpMngServer := server.NewManagementServer(repo)

	httpRouterServer.ListenAndServe()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8081))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
		log.Info().Msgf("Management server listening on %s", lis.Addr().String())
		if err := httpMngServer.Serve(lis); err != nil {
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
		if err := httpRouterServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		httpMngServer.Stop()
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
