// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package handlers

import "log"

// Products groups all product-related handler methods.
// It stores a logger that is shared among all handler functions.
type Products struct {
	l *log.Logger
}

// NewProducts returns a new Products handler.
// This function injects a logger dependency into Products.
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}
