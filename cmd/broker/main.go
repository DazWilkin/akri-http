package main

import (
	"flag"
	"log"
	"net"
	"os"

	pb "github.com/DazWilkin/akri-http/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	deviceEndpoint = "AKRI_HTTP_DEVICE_ENDPOINT"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", "", "The endpoint of this gRPC server.")
)

func main() {
	log.Println("[main] Starting gRPC server")

	flag.Parse()
	if *grpcEndpoint == "" {
		log.Fatal("[main] Unable to start server. Requires gRPC endpoint.")
	}

	deviceURL := os.Getenv(deviceEndpoint)
	if deviceURL == "" {
		log.Fatalf("Unable to determine Device URL using environment: %s", deviceEndpoint)
	}

	serverOpts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(serverOpts...)

	// Register this module's DeviceService
	pb.RegisterDeviceServiceServer(grpcServer, NewServer(deviceURL))

	// Register Golang Healthcheck service
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())

	listen, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[main] Starting gRPC Listener [%s]\n", *grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))
}
