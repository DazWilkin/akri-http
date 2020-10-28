package main

import (
	"context"
	"flag"
	"log"
	"net"

	pb "github.com/DazWilkin/akri-http/protos"

	"google.golang.org/grpc"
)

var _ pb.HTTPServer = (*Server)(nil)

var (
	grpcEndpoint = flag.String("grpc_endpoint", "", "The endpoint of this gRPC server.")
)

// Server is a type that implements pb.HTTPServer
type Server struct{}

// NewServer is a function that returns a new Server
func NewServer() *Server {
	return &Server{}
}

// Service is a method that implements the pb.HTTPServer interface
func (s *Server) Service(ctx context.Context, rqst *pb.ServiceRequest) (*pb.ServiceResponse, error) {
	return &pb.ServiceResponse{}, nil
}

func main() {
	log.Println("[main] Starting gRPC server")

	flag.Parse()
	if *grpcEndpoint == "" {
		log.Fatal("[main] Unable to start server. Requires gRPC endpoint.")
	}

	serverOpts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterHTTPServer(grpcServer, NewServer())

	listen, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[main] Starting gRPC Listener [%s]\n", *grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))
}
