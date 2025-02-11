package main

import (
	"context"
	common "github.com/JanKoczuba/commons"
	"github.com/JanKoczuba/commons/broker"
	"github.com/JanKoczuba/commons/discovery"
	"github.com/JanKoczuba/commons/discovery/consul"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	serviceName = "payment"
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	grpcAddr    = common.EnvString("GRPC_ADDRESS", "localhost:2001")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {

	// Register consul
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)

	ctx := context.Background()
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatalf("failed to health check %v", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Broker connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	svc := NewService()

	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("GRPC Server Started at ", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
