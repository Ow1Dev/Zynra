package server

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/Ow1Dev/Zynra/internal/repository"
	"github.com/Ow1Dev/Zynra/pkgs/httpsutils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
)

func sendAction(addr string, ctx context.Context) (*string, error) {
	start := time.Now()
	log.Info().Msgf("Connecting to service at %s", addr)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}
	defer conn.Close()
	c := pb.NewGatewayServiceClient(conn)

	r, err := c.Execute(ctx, &pb.ExecuteRequest{Name: "Action"})
	if err != nil {
		return nil, fmt.Errorf("could not greet: %w", err)
	}
	log.Info().Msgf("Action sent to service at %s, duration: %s", addr, time.Since(start))

	msg := r.GetMessage()
	return &msg, nil
}

func handleRunner(repo *repository.ServiceRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		cleanPath = strings.Trim(cleanPath, "/")

		segments := strings.Split(cleanPath, "/")

		log.Debug().Any("segments", segments).Msg("Request path segments")

		if len(segments) != 1 || segments[0] == "" {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		action := segments[0]
		service, exists := repo.GetService(action)
		if !exists {
			http.Error(w, fmt.Sprintf("Service '%s' not found", action), http.StatusNotFound)
			return
		}

		msg, err := sendAction(service.Address, r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to send action")
			http.Error(w, "Failed to send action", http.StatusInternalServerError)
			return
		}

		httpsutils.Encode(w, http.StatusOK, msg)
	})
}

func NewRouterServer(repo *repository.ServiceRepository) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleRunner(repo))
	var handler http.Handler = mux
	return handler
}
