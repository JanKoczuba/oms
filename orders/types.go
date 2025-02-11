package main

import (
	"context"
	pb "github.com/JanKoczuba/commons/api"
)

type OrdersService interface {
	CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(ctx context.Context) error
}
