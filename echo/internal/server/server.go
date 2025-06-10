package internal

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
}

func doBar(ctx context.Context) (any, error) {
	_ = ctx // Use ctx if needed for context-aware operations

	return map[string]string{
		"message": "Hello from the bar action",
		"SomeData": "This is some data from the bar action",
	}, nil
}

func doFoo(ctx context.Context) (any, error) {
	_ = ctx // Use ctx if needed for context-aware operations

	return map[string]string{
		"message": "Hello from the foo action",
		"status":  "success",
	}, nil
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

	switch action {
	case "bar":
			result, err = doBar(ctx)
	case "foo":
			result, err = doFoo(ctx)
	default:
		 //TODO: should we return an error here?
			log.Error().Msgf("Unknown action: %s", action)
			return nil, fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
			log.Error().Err(err).Msgf("Failed to execute action: %s", action)
			return nil, fmt.Errorf("failed to execute action: %s, error: %w", action, err)
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

func NewTunnelServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterGatewayServiceServer(s, &gatewayServiceServer{})
	return s
}
