package server

import (
	"context"
	"mailganer/internal/handler"
	"mailganer/pkg/logger"
	"net/http"
	"time"
)

type Server struct {
	Logger     logger.Logger
	httpServer *http.Server
}

func NewServer(logger logger.Logger, handler handler.Handler, host, port string) *Server {
	srv := &http.Server{
		Addr:           host + ":" + port,
		Handler:        handler.InitRoutes(),
		MaxHeaderBytes: 1 << 20, //1 Mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return &Server{
		Logger:     logger,
		httpServer: srv,
	}
}

func (s *Server) Start() error {

	s.Logger.Infof("Server starts at %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
