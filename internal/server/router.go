package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Ow1Dev/Zynra/internal/repository"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
	"github.com/Ow1Dev/Zynra/pkgs/httpsutils"
)

func sendAction(addr string, action string, ctx context.Context) (*string, error) {
	start := time.Now()
	log.Info().Msgf("Connecting to service at %s", addr)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}
	defer conn.Close()
	c := pb.NewGatewayServiceClient(conn)

	r, err := c.Execute(ctx, &pb.ExecuteRequest{Name: action})
	if err != nil {
		return nil, fmt.Errorf("could not greet: %w", err)
	}
	log.Info().Msgf("Action sent to service at %s, duration: %s", addr, time.Since(start))

	msg := r.GetMessage()
	return &msg, nil
}

func handleRunner(repo *repository.ServiceRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

    // Check if we have at least two parts: service and action
    if len(parts) < 2 {
        http.Error(w, "URL must be in the format /service/action", http.StatusBadRequest)
        return
    }

    servicePath := strings.TrimSpace(parts[0])
    if servicePath == "" {
        http.Error(w, "Service parameter is required", http.StatusBadRequest)
        return
    }

    actionPath := strings.TrimSpace(parts[1])
    if actionPath == "" {
        http.Error(w, "Action parameter is required", http.StatusBadRequest)
        return
    }

		service, exists := repo.GetService(servicePath)
		if !exists {
			http.Error(w, fmt.Sprintf("Service '%s' not found", actionPath), http.StatusNotFound)
			return
		}

		msg, err := sendAction(service.Address, actionPath, r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to send action")
			http.Error(w, "Failed to send action", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(*msg))
	})
}

func HealthCheckHandler(repo *repository.ServiceRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpsutils.Encode(w, http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})
}

func NewRouterServer(repo *repository.ServiceRepository) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", HealthCheckHandler(repo))
	mux.HandleFunc("GET /{service}/{action}", handleRunner(repo))
	var handler http.Handler = mux
	return handler
}
