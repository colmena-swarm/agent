package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	colmena "colmena.bsc.es/role-selector/policy/grpc/colmena_consensus"
	"colmena.bsc.es/role-selector/types"
)

type Server struct {
	colmena.UnimplementedSchedulingServiceServer
	grpcServer *grpc.Server
	port       string
	trigger    func(types.Decision)
}

func NewServer(port string, trigger func(types.Decision)) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
		port:       port,
		trigger:    trigger,
	}
}

func (s *Server) TriggerRole(ctx context.Context, req *colmena.TriggerRoleRequest) (*emptypb.Empty, error) {
	fmt.Printf("TriggerRole called for role: %s, startOrStop: %v\n", req.GetRoleId(), req.GetStartOrStop())
	s.trigger(types.Decision{RoleId: req.RoleId, ServiceId: req.ServiceId, StartOrStop: req.StartOrStop})
	return &emptypb.Empty{}, nil
}

func (s *Server) Start(ctx context.Context) {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	colmena.RegisterSchedulingServiceServer(s.grpcServer, s)

	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()

	log.Printf("gRPC server listening on %s", s.port)
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}
