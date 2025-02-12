package gateway

import (
	"context"
	pb "github.com/JanKoczuba/commons/api"
	"github.com/JanKoczuba/commons/discovery"
	"log"
)

type gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) UpdateOrderAfterPaymentLink(ctx context.Context, orderID, paymentLink string) error {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewOrderServiceClient(conn)

	_, err = ordersClient.UpdateOrder(ctx, &pb.Order{
		ID:          orderID,
		Status:      "waiting_payment",
		PaymentLink: paymentLink,
	})
	return err
}
