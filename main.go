package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/VINAYAK777CODER/PRODUCT-API/handlers"
	appmw "github.com/VINAYAK777CODER/PRODUCT-API/middleware"
	swaggerMw "github.com/go-openapi/runtime/middleware"
)

func main() {

	// 1️⃣ Create logger (dependency) for handlers
	// prefix = PRODUCT-API, add timestamp
	l := log.New(os.Stdout, "PRODUCT-API ", log.LstdFlags)

	// 2️⃣ Create handlers + inject logger
	// hh := handlers.NewHello(l)
	// gh := handlers.NewGoodBye(l)
	ph := handlers.NewProducts(l)

	// 3️⃣ Create router
	// sm := http.NewServeMux()

	// using gin creating router

	r := gin.Default()

	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", ph.GetProducts)
		productRoutes.POST("/", appmw.MethodCheck(), appmw.JSONCheck(), ph.AddProduct)
		productRoutes.PUT("/:id", appmw.MethodCheck(), appmw.JSONCheck(), ph.UpdateProduct)
		productRoutes.DELETE("/:id",ph.Delete)
	}

	// serve swagger file
	r.StaticFile("/swagger.yaml", "./swagger.yaml")

	// serve redoc UI
	opts := swaggerMw.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := swaggerMw.Redoc(opts, nil)
	r.GET("/docs", gin.WrapH(sh))

	// 4️⃣ Register routes
	// sm.Handle("/", hh)
	// sm.Handle("/goodbye", gh)
	// sm.Handle("/product",ph).

	// 5️⃣ Create a fully configured server
	s := &http.Server{
		Addr:         ":9090",           // server port
		Handler:      r,                 // router
		IdleTimeout:  120 * time.Second, // disconnect idle clients
		ReadTimeout:  1 * time.Second,   // read time limit
		WriteTimeout: 1 * time.Second,   // write time limit
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
