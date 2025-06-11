package main

import (
	"context"
	"flag"

	zynra "github.com/Ow1Dev/Zynra/pkgs/sdk/Zynra"
	"github.com/rs/zerolog/log"
)

// Example on how to use a custom logger
type zeroLogger struct {}

func (l *zeroLogger) LogInfo(format string, v ...any) {
	log.Info().Msgf(format, v...)
}

func (l *zeroLogger) LogWarn(format string, v ...any) {
	log.Warn().Msgf(format, v...)
}

func (l *zeroLogger) LogError(format string, v ...any) {
	log.Error().Msgf(format, v...)
}


// Example on how to make a handler
func fooHandler(_ context.Context) (any, error) {
	type response struct {
    Status string `json:"status"`
    Data struct {
        Message string `json:"message"`
    } `json:"data"`
}

	var resp response
	resp.Status = "success"
	resp.Data.Message = "Hello, world!"

	return resp, nil
}

// The main function to show how to set it up
func main() {
	addr := flag.String("addr", "localhost:8081", "Address to connect to")
	port := flag.Uint("port", 1234, "Port to connect to")
	flag.Parse()

	// Setup context and service
	ctx := context.Background()
	service := zynra.NewService(*addr)
	service.SetLogger(&zeroLogger{})
	defer service.Stop()

	// Register actions
	service.AddAction("foo", fooHandler)

	log.Info().Msgf("Starting Zynra service at %s:%d...", *addr, *port)

	// Start listening (consider handling error if Listen can fail)
	if err := service.Listen(uint32(*port), ctx); err != nil {
		log.Fatal().Msgf("Failed to start service: %v", err)
	}
}
