package server

import (
	"context"
	"log"
	"net"

	"github.com/envoyproxy/go-control-plane/api"
	pb "github.com/taiki45/eds/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

// registration
func (s *server) RegisterHosts(ctx context.Context, in *pb.HostsRegistrationRequest) (*pb.HostsRegistrationResponse, error) {
	return &pb.HostsRegistrationResponse{}, nil
}

func (s *server) DeregisterHosts(ctx context.Context, in *pb.HostsDeregistrationRequest) (*pb.HostsDeregistrationResponse, error) {
	return &pb.HostsDeregistrationResponse{}, nil
}

// fetching
func (s *server) FetchEndpoints(ctx context.Context, in *api.DiscoveryRequest) (*api.DiscoveryResponse, error) {
	return &api.DiscoveryResponse{}, nil
}

func (s *server) StreamEndpoints(stream api.EndpointDiscoveryService_StreamEndpointsServer) error {
	return nil
}

func (s *server) StreamLoadStats(stream api.EndpointDiscoveryService_StreamLoadStatsServer) error {
	return nil
}

func Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRegistrationApiServer(s, &server{})
	api.RegisterEndpointDiscoveryServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
