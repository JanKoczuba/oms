package main

import (
	"context"
	common "github.com/JanKoczuba/commons"
	"github.com/JanKoczuba/commons/discovery"
	"github.com/JanKoczuba/commons/discovery/consul"
	"github.com/JanKoczuba/oms-gateway/gateway"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"time"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Printf("health check failed: %v", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server")
	}
}
