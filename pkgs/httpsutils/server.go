package httpsutils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/Ow1Dev/Zynra/internal/config"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(handler http.Handler, port string, config config.Config) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    net.JoinHostPort(config.Host, port),
			Handler: handler,
		},
	}
}

func (s *HTTPServer) ListenAndServe() {
	go func() {
		log.Info().Msgf("listening on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
}

func (s *HTTPServer) Shutdown(shutdownCtx context.Context) error {
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}
