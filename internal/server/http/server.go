package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/IKolyas/image-previewer/internal/logger"
	"github.com/IKolyas/image-previewer/internal/storage/source"
)

type Server struct {
	server      *http.Server
	storage     source.Storage
	middlewares []func(next http.Handler) http.Handler
	logger      *logger.Logger
}

type Option func(*Server)

func WithMaxBodySize(size int64) Option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, limitBodySize(size))
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, timeoutMiddleware(timeout))
	}
}

func NewServer(addr string, storage source.Storage, logger *logger.Logger, opts ...Option) (*Server, error) {
	srv := &Server{
		storage: storage,
		logger:  logger,
	}

	for _, opt := range opts {
		opt(srv)
	}

	handler := srv.setupRoutes()
	srv.server = &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return srv, nil
}

func (s *Server) setupRoutes() http.Handler {
	router := http.NewServeMux()

	h := &PreviewerHandler{
		server: *s,
	}

	router.HandleFunc("/fill/", h.fill)

	var handler http.Handler = router
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}

	return handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.server.Addr, err)
	}

	s.logger.Info(fmt.Sprintf("Starting HTTP server on %s", s.server.Addr))

	if serveErr := s.server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", serveErr)
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
