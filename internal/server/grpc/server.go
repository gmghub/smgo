package internalgrpc

import (
	"errors"
	"log"
	"net"

	pb "github.com/gmghub/smgo/pkg/smgo/api"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server  *grpc.Server
	App     Application
	Service *Service
}

type Application interface {
	StatJSON(period int, collectors ...string) []byte
	MapCollectors() map[string]int
}

func NewServer(app Application) *GRPCServer {
	return &GRPCServer{App: app}
}

func (s *GRPCServer) Start(addr string) error {
	s.Server = grpc.NewServer(
		// grpc.KeepaliveParams(),
		grpc.ChainStreamInterceptor(
			loggingStreamInterceptor,
		),
	)
	s.Service = NewService(s.App)
	pb.RegisterSmgoServiceServer(s.Server, s.Service)
	// reflection.Register(s.Server)

	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer listener.Close()
		if err := s.Server.Serve(listener); err != nil {
			log.Fatal(err.Error())
		}
		log.Println("gRPC server stopped")
	}()

	return nil
}

func (s *GRPCServer) Stop() error {
	if s.Server == nil {
		return errors.New("gRPC server not started")
	}
	s.Service.Close()
	s.Server.GracefulStop()
	return nil
}
