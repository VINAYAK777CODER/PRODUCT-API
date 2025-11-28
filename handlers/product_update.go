package handlers

import (
	"WORKING-GO/data"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateProduct handles: PUT /product/:id
// ----------------------------------------------------------------------------
// Purpose:
//   Update an existing product identified by its ID.
//
// Steps:
//   1. Extract and convert :id route parameter.
//   2. Bind JSON body to Product struct.
//   3. Validate updated fields.
//   4. Call data.UpdateProduct() to modify stored product.
//   5. Return appropriate success or error messages.
//
func (p *Products) UpdateProduct(c *gin.Context) {
	p.l.Println("PUT /product/:id called")

	// Extract ID from URL (string → int)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Prepare empty struct for updated JSON data
	prod := &data.Product{}

	// Bind JSON
	if err := c.BindJSON(prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// Validate fields before updating
	if err := prod.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	// Try updating product
	err = data.UpdateProduct(id, prod)

	// Product not found
	if err == data.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Any unknown error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}
