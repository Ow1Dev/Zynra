package server

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/Ow1Dev/Zynra/pkgs/httpsutils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
)

func sendAction(ctx context.Context) error {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to management server: %w", err)
	}
	defer conn.Close()
	c := pb.NewGatewayServiceClient(conn)

	r, err := c.Execute(ctx, &pb.ExecuteRequest{Name: "Echo Service"})
	if err != nil {
		return fmt.Errorf("could not greet: %w", err)
	}
	log.Printf("message: %s", r.GetMessage())

	return nil
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	cleanPath = strings.Trim(cleanPath, "/")

	segments := strings.Split(cleanPath, "/")

	log.Debug().Any("segments", segments).Msg("Request path segments")

	if len(segments) != 1 || segments[0] == "" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	action := segments[0]
	if(action != "echo") {
		http.Error(w, "Only 'echo' is allowed", http.StatusBadRequest)
		return
	}

	err := sendAction(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send action")
		http.Error(w, "Failed to send action", http.StatusInternalServerError)
		return
	}

	httpsutils.Encode(w, http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}

func NewRouterServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getRoot)
	var handler http.Handler = mux
	return handler
}
