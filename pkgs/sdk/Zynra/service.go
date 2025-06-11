package zynra

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/managment"
)

type ActionHandler func(ctx context.Context) (any, error)

type ZynraService struct {
	mdrAddr string
	grpcServer *grpc.Server
	actions map[string]ActionHandler
	logger Logger
}

func NewService(mngAddr string) *ZynraService {
	return &ZynraService{
		mdrAddr: mngAddr,
		actions: make(map[string]ActionHandler),
		logger: &stdLogger{},
	}
}

func (s *ZynraService) SetLogger(l Logger) {
	if l != nil {
		s.logger = l
	}
}

func (s *ZynraService) AddAction(action string, handler ActionHandler) {
	s.actions[action] = handler
}

func (s *ZynraService) Listen(port uint32, ctx context.Context) error {
	err := s.connectToManagementServer(ctx, port, &s.mdrAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to management server: %w", err)
	}

	s.grpcServer = newTunnelServer(s.actions) 

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	s.logger.LogInfo("Gateway server listening on %s", lis.Addr().String())
	if err := s.grpcServer.Serve(lis); err != nil {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}

	// Implementation for listening to an action
	return nil
}

func (s *ZynraService) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		s.logger.LogInfo("Zynra service stopped gracefully")
	} else {
		s.logger.LogWarn("Zynra service was not running")
	}
}

func (s *ZynraService) connectToManagementServer(ctx context.Context, port uint32, addr *string) error {
	s.logger.LogInfo("Connecting to management server at %s", *addr)
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to management server: %w", err)
	}
	defer conn.Close()
	c := pb.NewManagementServiceClient(conn)
	r, err := c.Connect(ctx, &pb.ConnectRequest{Name: "test", Port: port})
	if err != nil {
		return fmt.Errorf("could not greet: %w", err)
	}
	s.logger.LogInfo("Connected to management server: %s", r.GetMessage())
	conn.Close()
	s.logger.LogInfo("Connected to management server at %s", *addr)

	return nil
}
