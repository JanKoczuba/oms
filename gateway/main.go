package main

import (
	common "github.com/JanKoczuba/commons"
	pb "github.com/JanKoczuba/commons/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	httpAddr         = common.EnvString("HTTP_ADDR", ":8080")
	orderServiceAddr = "localhost:3000"
)

func main() {
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to make new client: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server")
	}
}
