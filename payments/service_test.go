package main

import (
	"context"
	"github.com/JanKoczuba/oms-payments/gateway"
	"testing"

	"github.com/JanKoczuba/commons/api"
	inmemRegistry "github.com/JanKoczuba/commons/discovery/inmem"
	"github.com/JanKoczuba/oms-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()

	gateway := gateway.NewGateway(registry)
	svc := NewService(processor, gateway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
