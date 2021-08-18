package server

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
