package main

import (
	"context"
	common "github.com/JanKoczuba/commons"
	"github.com/JanKoczuba/commons/discovery"
	"github.com/JanKoczuba/commons/discovery/consul"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:3000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Printf("health check failed: %v", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to connect to grpc server: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := newService(store)
	NewGrpcHandler(grpcServer, svc)

	svc.CreateOrder(context.Background())

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf(err.Error())
	}
}
