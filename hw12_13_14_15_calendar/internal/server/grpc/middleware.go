package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func (s *Server) loggingMiddleware(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (resp interface{}, err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return resp, errors.New("error peer from ctx")
	}
	t := time.Now()
	resp, err = next(ctx, req)
	latency := fmt.Sprintf("%dms", time.Since(t).Milliseconds())
	ret := fmt.Sprintf("middleware: %s %s %s\n", p.Addr, info.FullMethod, latency)
	s.logg.Info(ret)
	return resp, err
}
