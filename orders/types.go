package main

import (
	"context"
	pb "github.com/JanKoczuba/commons/api"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(ctx context.Context, o Order) (bson.ObjectID, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(ctx context.Context, id string, o *pb.Order) error
}

type Order struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	CustomerID  string        `bson:"customerID,omitempty"`
	Status      string        `bson:"status,omitempty"`
	PaymentLink string        `bson:"paymentLink,omitempty"`
	Items       []*pb.Item    `bson:"items,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		ID:          o.ID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
