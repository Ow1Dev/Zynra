package server

import (
	"context"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/managment"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type managementServiceServer struct {
	 pb.UnimplementedManagementServiceServer
}

// Connect implements managment.ManagementServiceServer.
func (m *managementServiceServer) Connect(context.Context, *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	log.Info().Msg("Client connected")
	return &pb.ConnectResponse{
		Message: "Connected successfully",
	}, nil
}

func NewManagementServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterManagementServiceServer(s, &managementServiceServer{})
	return s
}
