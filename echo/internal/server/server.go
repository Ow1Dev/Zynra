package internal

import (
	"context"
	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
	"google.golang.org/grpc"
)

type gatewayServiceServer struct {
	pb.UnimplementedGatewayServiceServer
}

// Execute implements gateway.GatewayServiceServer.
func (g *gatewayServiceServer) Execute(ctx context.Context, request *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
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
