package main

import (
	"context"
	"errors"
	pb "github.com/JanKoczuba/commons/api"
)

var orders = make([]*pb.Order, 0)

type store struct {
	// add mongoDB instance
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	orders = append(orders, &pb.Order{
		ID:          "22",
		CustomerID:  p.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})
	return "22", nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {
	for _, o := range orders {
		if o.ID == id && o.CustomerID == customerID {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, newOrder *pb.Order) error {
	for i, o := range orders {
		if o.ID == id {
			orders[i].Status = o.Status
			orders[i].PaymentLink = o.PaymentLink
			return nil
		}
	}

	return nil
}
