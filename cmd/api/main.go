package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/vivek700/todo-server/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//listen for the interrupt signal
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	//the context is used to inform the server it has 5 seconds to finish
	//the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
	log.Println("Server exiting...")

	//Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	server := server.NewServer()

	//Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	//wait for the graceful shutdown to complete

	<-done
	log.Println("Graceful shutdown complete")

}