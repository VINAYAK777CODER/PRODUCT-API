package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Hello struct holds dependencies (like logger).
// So handler functions can use h.l to log things.
type Hello struct {
	l *log.Logger
}

// Constructor function - returns Hello with logger injected.
// This is Dependency Injection (DI).
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHTTP is AUTOMATICALLY called by Go's HTTP server.
// Any struct with ServeHTTP becomes an http.Handler.
//
// FLOW:
// - Client sends request: POST /
// - Go server receives it
// - Router (ServeMux) finds correct handler
// - And automatically calls:
//       h.ServeHTTP(responseWriter, request)
//
// You NEVER call ServeHTTP yourself.
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Log that the request came
	h.l.Println("helloworld!")

	// Read request body
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "unable to read request body", http.StatusBadRequest)
		return
	}

	// Send response back to client
	fmt.Fprintf(rw, "hello %s", d)
}
