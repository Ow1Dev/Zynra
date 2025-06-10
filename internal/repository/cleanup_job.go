package repository

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Ow1Dev/Zynra/pkgs/api/gateway"
)

type CleanupJob struct {
	repo *ServiceRepository
}

func NewCleanupJob(repo *ServiceRepository) *CleanupJob {
	return &CleanupJob{
		repo: repo,
	}
}

func (j *CleanupJob) cleanupServiceWork() error {
	// Iterate over all services in the repository
	for name, entity := range j.repo.services {
		log.Debug().Msgf("Checking service %s at %s", name, entity.Address)

		// Check if the service is still reachable
		addr := entity.Address
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error().Err(err).Msgf("Failed to connect to service %s at %s", name, addr)
			log.Debug().Msgf("Removing service %s from repository", name)
			j.repo.RemoveService(name)
			continue
		}

		client := pb.NewGatewayServiceClient(conn)
		_, err = client.Ping(context.Background(), &pb.PingRequest{})

		// Close connection right after ping
		conn.Close()

		if err != nil {
			log.Error().Err(err).Msgf("Service %s at %s is not reachable", name, addr)
			log.Debug().Msgf("Removing service %s from repository", name)
			j.repo.RemoveService(name)
			continue
		}
	}

	return nil
}


func (j *CleanupJob) Start(ctx context.Context, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                log.Info().Msg("CleanupJob shutting down")
                return
            case <-ticker.C:
							  log.Debug().Msg("Running cleanup job for services")
                if err := j.cleanupServiceWork(); err != nil {
                    log.Error().Err(err).Msg("CleanupJob error")
                }
            }
        }
    }()
}
