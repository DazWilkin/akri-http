package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	host = "0.0.0.0"
)

var (
	portDisco = flag.Int("discovery_port", 9999, "A non-overlapping port on which the discovery service will run.")
	portStart = flag.Int("starting_port", 8000, "The lowest port on which a device should be created.")
	portCount = flag.Int("num_devices", 10, "The number of devices that should be created (on incrementing ports).")
)

var wg sync.WaitGroup

func main() {
	flag.Parse()
	// [TODO:dazwilkin] Add checks on portStart and portCount values

	// Create|Register Device(s)
	devices := make([]string, *portCount)
	for p := 0; p < *portCount; p++ {
		port := *portStart + p
		addr := fmt.Sprintf("%s:%d", host, port)

		// Create Device service
		log.Printf("[main] Creating Device: %s", addr)
		wg.Add(1)
		go createDeviceService(addr)

		// Register Device
		log.Printf("[main] Register Device")
		devices[p] = addr
	}

	// Create Discovery service
	addr := fmt.Sprintf("%s:%d", host, *portDisco)
	wg.Add(1)
	go createDiscoveryService(addr, devices)

	wg.Wait()
}
func createDeviceService(addr string) {
	defer wg.Done()

	seed := rand.NewSource(time.Now().UnixNano())
	entr := rand.New(seed)

	// Handlers: Devices, Healthz
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", healthz)
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[device:%s] Handler entered", addr)
		fmt.Fprint(w, entr.Float64())
	})

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[createDeviceService] Starting Device: %s", addr)
	log.Fatal(s.Serve(listen))
}
func createDiscoveryService(addr string, devices []string) {
	defer wg.Done()

	// Handlers: Devices, Healthz
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", healthz)
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[discovery] Handler entered")
		for _, device := range devices {
			fmt.Fprintf(w, "%s\n", html.EscapeString(device))
		}
	})

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[createDiscoveryService] Starting Discovery Service: %s", addr)
	log.Fatal(s.Serve(listen))
}
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
