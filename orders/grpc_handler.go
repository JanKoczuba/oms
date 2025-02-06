package main

import (
	"context"
	pb "github.com/JanKoczuba/commons/api"
	"google.golang.org/grpc"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Panicln("New order created")
	return nil, nil
}
