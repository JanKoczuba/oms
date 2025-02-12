package main

import (
	"context"
	"encoding/json"
	pb "github.com/JanKoczuba/commons/api"
	"github.com/JanKoczuba/commons/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, p)
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order recived! Order %v", p)
	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	items, err := h.service.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(ctx, p, items)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return o, nil
}
