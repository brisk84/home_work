package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	httpServer http.Server
	logg       Logger
	app        Application
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		ret := fmt.Sprintf("%s [%s] %s %s %s", ip, time.Now(), req.Method, req.URL.Path, req.UserAgent())
		fmt.Fprintf(w, "<p>%s</p>", ret)
	})

	server := &Server{
		httpServer: http.Server{
			Addr:    addr,
			Handler: loggingMiddleware(mux, logger),
		},
		logg: logger,
		app:  app,
	}
	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("Start server")
	err := s.httpServer.ListenAndServe()
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stop server")
	err := s.httpServer.Shutdown(ctx)
	return err
}
