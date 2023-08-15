package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/v1/readiness", readiness)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Println("server listening on", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	StopGracefully(&server, serverErr)

}

func StopGracefully(server *http.Server, serverErr chan error) {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT)

	select {
	case sig := <-shutdownChannel:
		log.Println("signal:", sig)

		const timeout = 60 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
		}
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server: %v", err)
		}
	}
}

func readiness(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-REQUEST-ID")
	log.Println("start", requestID)
	defer log.Println("done", requestID)

	time.Sleep(5 * time.Second)

	response := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		panic(err)
	}
}
