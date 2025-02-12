package main

import (
	"context"
	common "github.com/JanKoczuba/commons"
	pb "github.com/JanKoczuba/commons/api"
	"log"
)

type service struct {
	store OrdersStore
}

func newService(store OrdersStore) *service {
	return &service{store}
}

func (s *service) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, o.ID, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	o, err := s.store.Get(ctx, p.OrderID, p.CustomerID)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {

	id, err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}

	//TODO remove hardcoded
	o := &pb.Order{
		ID:         id,
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
	}
	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Print(mergedItems)

	//TODO test only, remove after
	var itemWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemWithPrice = append(itemWithPrice, &pb.Item{
			PriceID:  "price_1QrFEnB4QS5L2w5cEBWbOYEw",
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	return itemWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
