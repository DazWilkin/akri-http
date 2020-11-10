package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/DazWilkin/akri-http/protos"

	"google.golang.org/grpc"
)

var _ pb.DeviceServiceServer = (*Server)(nil)

const (
	deviceEndpoint = "AKRI_HTTP_DEVICE_ENDPOINT"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", "", "The endpoint of this gRPC server.")
)

// Server is a type that implements pb.DeviceServiceServer
type Server struct {
	DeviceURL string
}

// NewServer is a function that returns a new Server
func NewServer(deviceURL string) *Server {
	return &Server{
		DeviceURL: deviceURL,
	}
}

// ReadSensor is a method that implements the pb.HTTPServer interface
func (s *Server) ReadSensor(ctx context.Context, rqst *pb.ReadSensorRequest) (*pb.ReadSensorResponse, error) {
	log.Println("[read_sensor] Entered")
	resp, err := http.Get(s.DeviceURL)
	if err != nil {
		return &pb.ReadSensorResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Printf("[read_sensor] Response status: %d", resp.StatusCode)
		return &pb.ReadSensorResponse{}, fmt.Errorf("response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &pb.ReadSensorResponse{}, err
	}

	log.Printf("[read_sensor] Response body: %s", body)
	return &pb.ReadSensorResponse{
		Value: string(body),
	}, nil
}

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
	pb.RegisterDeviceServiceServer(grpcServer, NewServer(fmt.Sprintf("http://%s", deviceURL)))

	listen, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[main] Starting gRPC Listener [%s]\n", *grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))
}
