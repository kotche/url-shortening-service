package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/kotche/url-shortening-service/internal/app/config"
	grpcHandler "github.com/kotche/url-shortening-service/internal/app/transport/grpc"
	"github.com/kotche/url-shortening-service/internal/app/transport/grpc/interceptors"
	pb "github.com/kotche/url-shortening-service/internal/app/transport/grpc/proto"
	"google.golang.org/grpc"
)

type Server struct {
	cfg        *config.Config
	handler    *grpcHandler.Handler
	grpcServer *grpc.Server
}

func NewServer(cfg *config.Config, handler *grpcHandler.Handler) *Server {
	authInterceptor := grpc.UnaryInterceptor(interceptors.UnaryCookieInterceptor)

	return &Server{
		cfg:        cfg,
		handler:    handler,
		grpcServer: grpc.NewServer(authInterceptor),
	}
}

func (s *Server) Run() error {
	pb.RegisterShortenerServer(s.grpcServer, s.handler)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.GRPCPort))
	if err != nil {
		log.Fatal(err)
	}
	return s.grpcServer.Serve(listen)
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}
