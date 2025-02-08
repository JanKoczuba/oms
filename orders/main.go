package main

import (
	"context"
	common "github.com/JanKoczuba/commons"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:3000")
)

func main() {

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
