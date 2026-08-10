package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kikudesuyo/room-finder/api/library"
	"github.com/kikudesuyo/room-finder/api/route"
)

func main() {
	portFlag := flag.String("port", "8081", "port to run HTTP server on")
	flag.Parse()

	port := *portFlag
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := library.InitDB(); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: route.NewRouter(),
	}

	log.Printf("API server is running on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("API server failed to start: %v", err)
	}
}
