package main

import (
	"context"
	pb "github.com/JanKoczuba/commons/api"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {

	return "", nil
}
