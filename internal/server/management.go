package server

import (
	"context"
	"fmt"
	"net"

	"github.com/Ow1Dev/Zynra/internal/repository"
	pb "github.com/Ow1Dev/Zynra/pkgs/api/managment"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type managementServiceServer struct {
	 pb.UnimplementedManagementServiceServer
	 repo *repository.ServiceRepository
}

// Connect implements managment.ManagementServiceServer.
func (m *managementServiceServer) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	var clientIP string
	if p, ok := peer.FromContext(ctx); ok {
		var ipStr string
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			ip := addr.IP
			if ip.To4() == nil {
				// It's IPv6
				ipStr = "[" + ip.String() + "]"
			} else {
				// It's IPv4
				ipStr = ip.String()
			}
		} else {
			// Fallback: use Addr.String() but detect IPv6 format manually
			addrStr := p.Addr.String()
			host, _, err := net.SplitHostPort(addrStr)
			if err == nil && net.ParseIP(host) != nil {
				if net.ParseIP(host).To4() == nil {
					ipStr = "[" + host + "]"
				} else {
					ipStr = host
				}
			} else {
				ipStr = addrStr // totally fallback
			}
		}

		log.Info().Msgf("Received connection request from %s for service %s on port %d", ipStr, request.GetName(), request.GetPort())
	} else {
		log.Error().Msg("Failed to get peer information from context")
		return nil, status.Errorf(codes.Internal, "failed to get peer information from context")
	}
	url := fmt.Sprintf("%s:%d", clientIP, request.Port)
	m.repo.AddService(request.GetName(), url)

	return &pb.ConnectResponse{
		Message: "Connected successfully",
	}, nil
}

func NewManagementServer(repo *repository.ServiceRepository) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterManagementServiceServer(s, &managementServiceServer{
		repo: repo,
	})
	return s
}
