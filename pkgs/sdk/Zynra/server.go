package zynra

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type gatewayServiceServer struct {
	pb.UnimplementedGatewayServiceServer
	actions map[string]ActionHandler
}

// Ping implements gateway.GatewayServiceServer.
func (g *gatewayServiceServer) Ping(context.Context, *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{ Message: "pong" }, nil
}

// Execute implements gateway.GatewayServiceServer.
func (g *gatewayServiceServer) Execute(ctx context.Context, request *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Info().Msgf("Received request: %s", request.GetName())
	action := request.GetName()

	var (
			result any
			err    error
	)

	if handler, exists := g.actions[action]; exists {
		result, err = handler(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to execute action: %s", action)
			return nil, fmt.Errorf("failed to execute action: %s, error: %w", action, err)
		}
	} else {
		log.Error().Msgf("Unknown action: %s", action)
		return nil, fmt.Errorf("unknown action: %s", action)
	}

	messageJSON, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message to JSON")
		return nil, err
	}

	response := &pb.ExecuteResponse{
		Message: string(messageJSON),
	}

	return response, nil
}

func newTunnelServer(actions map[string]ActionHandler) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterGatewayServiceServer(s, &gatewayServiceServer{
		actions: actions,
	})
	return s
}
