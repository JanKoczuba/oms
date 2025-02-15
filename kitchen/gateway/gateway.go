package gateway

import (
	"context"

	pb "github.com/JanKoczuba/commons/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
