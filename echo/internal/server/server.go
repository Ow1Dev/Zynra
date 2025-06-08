package internal

import (
	"context"

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
	message :="Hello, " + request.GetName() 
	return &pb.ExecuteResponse{
		Message: message,
	}, nil
}

func NewTunnelServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterGatewayServiceServer(s, &gatewayServiceServer{})
	return s
}
