package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	addr = ":8080"
)

// Paths is an alias to use repeated flags with flag
type Paths []string

// String is a method required by flag.Value interface
func (e *Paths) String() string {
	result := strings.Join(*e, "\n")
	return result
}

// Set is a method required by flag.Value interface
func (e *Paths) Set(value string) error {
	log.Printf("[endpoints:Set] %s", value)
	*e = append(*e, value)
	return nil
}

var _ flag.Value = (*Paths)(nil)
var paths Paths

func main() {
	flag.Var(&paths, "path", "Repeat this flag to add paths for the device")
	flag.Parse()

	// At a minimum, respond on `/`
	if len(paths) == 0 {
		paths = []string{"/"}
	}
	log.Printf("[main] Paths: %d", len(paths))

	seed := rand.NewSource(time.Now().UnixNano())
	entr := rand.New(seed)

	// Handlers: Devices, Healthz
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[main:healthz] Handler entered")
		fmt.Fprint(w, "ok")
	})
	// Create handler for each endpoint
	for _, path := range paths {
		log.Printf("[main] Creating handler: %s", path)
		handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[main:handler] Handler entered: %s", path)
			fmt.Fprint(w, entr.Float64())
		})
	}

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[main] Starting Device: [%s]", addr)
	log.Fatal(s.Serve(listen))
}
