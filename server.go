package main

import (
	"context"
	"net/http"
)

type server struct {
	httpServer *http.Server
}

func newServer(addr string) *server {
	return &server{
		httpServer: &http.Server{
			Addr: addr,
		},
	}
}

func (s *server) start() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
