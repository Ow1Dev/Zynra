package internal

import (
	"context"
	"encoding/json"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type gatewayServiceServer struct {
	pb.UnimplementedGatewayServiceServer
}

// Execute implements gateway.GatewayServiceServer.
func (g *gatewayServiceServer) Execute(ctx context.Context, request *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Info().Msgf("Received request: %s", request.GetName())
	message := map[string]string{
		"message": "Hello from the echo service",
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message to JSON")
		return nil, err
	}

	response := &pb.ExecuteResponse{
		Message: string(messageJSON),
	}

	return response, nil
}

func NewTunnelServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterGatewayServiceServer(s, &gatewayServiceServer{})
	return s
}
