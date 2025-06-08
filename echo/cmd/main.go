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

	server "github.com/Ow1Dev/Zynra/echo/internal/server"
	pb "github.com/Ow1Dev/Zynra/pkgs/api/managment"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	addr := flag.String("addr", "localhost:8081", "the address to connect to")
	flag.Parse()

	initLog(w, *debug)

	// connnect ot the management server
	log.Info().Msg("Connecting to management server...")
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to management server: %w", err)
	}
	defer conn.Close()
	c := pb.NewManagementServiceClient(conn)
	defer cancel()
	r, err := c.Connect(ctx, &pb.ConnectRequest{Name: "Echo Service"})
	if err != nil {
		return fmt.Errorf("could not greet: %w", err)
	}
	log.Printf("message: %s", r.GetMessage())
	log.Info().Msg("Connected to management server")

	TunnelServer := server.NewTunnelServer() 

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8082))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
		log.Info().Msgf("Management server listening on %s", lis.Addr().String())
		if err := TunnelServer.Serve(lis); err != nil {
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
