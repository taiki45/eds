package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	port = 4000
)

type nodeGroupHasher struct {
}

// Hash function that always returns same value.
func (h nodeGroupHasher) Hash(*api.Node) (cache.Key, error) {
	return cache.Key("TODO"), nil
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
	i := 0
	for {
		version := fmt.Sprintf("version%d", i)
		endpoint := makeEndpoint()

		glog.Infof("updating cache with %d-labelled responses", i)
		snapshot := cache.NewSnapshot(version,
			[]proto.Message{endpoint},
			[]proto.Message{},
			[]proto.Message{},
			[]proto.Message{})
		c.SetSnapshot(cache.Key("TODO"), snapshot)

		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			return
		}
		i++
	}

}

func makeEndpoint() *api.ClusterLoadAssignment {
	return &api.ClusterLoadAssignment{
		ClusterName: "test",
		Endpoints: []*api.LocalityLbEndpoints{{
			LbEndpoints: []*api.LbEndpoint{{
				Endpoint: &api.Endpoint{
					Address: &api.Address{
						Address: &api.Address_SocketAddress{
							SocketAddress: &api.SocketAddress{
								Protocol: api.SocketAddress_TCP,
								Address:  "127.0.0.1",
								PortSpecifier: &api.SocketAddress_PortValue{
									PortValue: 3000,
								},
							},
						},
					},
				},
			}},
		}},
	}
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
