package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	port = 4000
)

type nodeGroupHasher struct {
}

// Hash function that always returns same value.
func (h nodeGroupHasher) Hash(node *api.Node) (cache.Key, error) {
	return cache.Key(node.Id), nil
}

// Run server
func Run() {
	ctx := context.Background()
	c := cache.NewSimpleCache(nodeGroupHasher{}, nil)
	go runGrpcServer(ctx, c, port)
	go runResouceUpdator(ctx, c)

	<-ctx.Done()
}

func runResouceUpdator(ctx context.Context, c cache.Cache) {
}

// RunGrpcServer starts an xDS server at the given port.
func runGrpcServer(ctx context.Context, c cache.Cache, port uint) {
	server := xds.NewServer(c)
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	api.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	log.Printf("Server listening on %d", port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			glog.Error(err)
		}
	}()
	<-ctx.Done()
	grpcServer.GracefulStop()
}
