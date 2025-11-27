package handlers

import (
	"WORKING-GO/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Products holds a logger and groups all product-related handler methods.
type Products struct {
	l *log.Logger
}

// NewProducts returns a new Products handler with the provided logger.
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GET /product
// Returns all products
func (p *Products) GetProducts(c *gin.Context) {
	p.l.Println("GET /product called")

	// retrieve product list from data layer
	lp := data.GetProductList()

	// return list as JSON
	c.JSON(http.StatusOK, lp)
}

// POST /product
// Adds a new product
func (p *Products) AddProduct(c *gin.Context) {
	p.l.Println("POST /product called")

	// holds incoming JSON data
	prod := &data.Product{}

	// parse request body into struct
	if err := c.BindJSON(prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// 2️⃣ Validate product input  🔥 REQUIRED
	if err := prod.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	// save product into in-memory store
	data.AddProduct(prod)

	c.JSON(http.StatusCreated, gin.H{"message": "product added", "data": prod})
}

// PUT /product/:id
// Updates an existing product by ID
func (p *Products) UpdateProduct(c *gin.Context) {
	p.l.Println("PUT /product/:id called")

	// extract id from URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// holds updated product data
	prod := &data.Product{}

	// parse JSON body
	if err := c.BindJSON(prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// Validate updated product
	if err := prod.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	// try updating in data layer
	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}
