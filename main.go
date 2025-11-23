package main

import (
	"WORKING-GO/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 1️⃣ Create logger (dependency) for handlers
	// prefix = PRODUCT-API, add timestamp
	l := log.New(os.Stdout, "PRODUCT-API ", log.LstdFlags)

	// 2️⃣ Create handlers + inject logger
	// hh := handlers.NewHello(l)
	// gh := handlers.NewGoodBye(l)
	ph:=handlers.NewProducts(l)

	// 3️⃣ Create router
	sm := http.NewServeMux()

	// 4️⃣ Register routes
	// sm.Handle("/", hh)
	// sm.Handle("/goodbye", gh)
	sm.Handle("/",ph)


	// 5️⃣ Create a fully configured server
	s := &http.Server{
		Addr:         ":9090",               // server port
		Handler:      sm,                    // router
		IdleTimeout:  120 * time.Second,     // disconnect idle clients
		ReadTimeout:  1 * time.Second,       // read time limit
		WriteTimeout: 1 * time.Second,       // write time limit
	}

	// 6️⃣ Start server asynchronously (goroutine)
	// so main goroutine can run shutdown logic
	go func() {
		log.Println("Server starting on port 9090...")
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// 7️⃣ Create channel for receiving OS shutdown signals
	sigchan := make(chan os.Signal, 1)

	// 8️⃣ Notify for SIGINT (CTRL+C) and SIGTERM (K8s / Docker stop)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// 9️⃣ Block until signal is received
	<-sigchan
	log.Println("Shutdown signal received...")

	// 🔟 Graceful shutdown using timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1️⃣1️⃣ Attempt graceful shutdown
	err := s.Shutdown(ctx)
	if err != nil {
		log.Println("Server shutdown error:", err)
	} else {
		log.Println("Server is shutting down gracefully...")
	}
}
