package server

import (
	"context"
	"net"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/managment"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type managementServiceServer struct {
	 pb.UnimplementedManagementServiceServer
}

// Connect implements managment.ManagementServiceServer.
func (m *managementServiceServer) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	var clientIP string
	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP = addr.IP.String()
		} else {
			clientIP = p.Addr.String() // fallback
		}
		log.Info().Msgf("Received connection request from %s:%d", clientIP, request.GetPort())
	} else {
		log.Error().Msg("Failed to get peer information from context")
		return nil, status.Errorf(codes.Internal, "failed to get peer information from context")
	}

	return &pb.ConnectResponse{
		Message: "Connected successfully",
	}, nil
}

func NewManagementServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterManagementServiceServer(s, &managementServiceServer{})
	return s
}
