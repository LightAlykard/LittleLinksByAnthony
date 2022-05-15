package server

import (
	"context"
	"net/http"
	"time"

	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
)

type Server struct {
	srv http.Server
	lk  *link.Links
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(lk *link.Links) {
	s.lk = lk
	go s.srv.ListenAndServe()
}
