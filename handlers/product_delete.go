package handlers

import (
	"net/http"
	"strconv"

	"github.com/VINAYAK777CODER/PRODUCT-API/data"
	"github.com/gin-gonic/gin"
)

// ----------------------------
// Swagger Models for DELETE
// ----------------------------

// swagger:parameters deleteProduct
type ProductIDParam struct {
	// Product ID to delete
	// in: path
	// required: true
	// minimum: 1
	// example: 1
	ID int `json:"id"`
}

// Successful delete response
//
// swagger:response deleteProductResponse
type DeleteProductResponse struct {
	// in: body
	Body struct {
		Message string `json:"message"`
		ID      int    `json:"id"`
	}
}

// Error response (not found)
//
// swagger:response deleteProductNotFound
type DeleteProductNotFound struct {
	// in: body
	Body struct {
		Error string `json:"error"`
	}
}

// Error response (invalid id / internal)
//
// swagger:response deleteProductError
type DeleteProductError struct {
	// in: body
	Body struct {
		Error string `json:"error"`
	}
}

// swagger:route DELETE /products/{id} products deleteProduct
//
// # Delete a product by ID
//
// This endpoint removes the product with the provided ID from the system.
//
// Responses:
//
//	200: deleteProductResponse
//	400: deleteProductError
//	404: deleteProductNotFound
//	500: deleteProductError
func (p *Products) Delete(c *gin.Context) {
	p.l.Println("DELETE /products/:id called")

	// Extract and convert ID from URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	p.l.Println("Deleting product with ID:", id)

	// Try to delete the product
	err = data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if err != nil {
		// Unknown internal error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}

	// Successfully deleted
	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
		"id":      id,
	})
}
