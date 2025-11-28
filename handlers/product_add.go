package handlers

import (
	"github.com/VINAYAK777CODER/PRODUCT-API/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// -----------------------------
// Swagger Response Definitions
// -----------------------------

// swagger:response productAddSuccess
type ProductAddSuccessResponse struct {
	// in:body
	Body struct {
		Message string       `json:"message"`
		Data    data.Product `json:"data"`
	}
}

// swagger:response productAddBadRequest
type ProductAddBadRequestResponse struct {
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// swagger:response productAddValidationError
type ProductAddValidationErrorResponse struct {
	// in:body
	Body struct {
		ValidationError string `json:"validation_error"`
	}
}

// ------------------------------------
// FIXED Swagger Route Documentation
// ------------------------------------

//
// swagger:route POST /products products addProduct
//
// Add a new product
//
// This endpoint creates a new product using JSON data sent in the request body.
//
// Responses:
//   201: productAddSuccess
//   400: productAddBadRequest
//   400: productAddValidationError
//

// AddProduct handles: POST /product
// ------------------------------------------------------------------------
// Purpose:
//   Add a new product using JSON data provided in the request.
//
// Steps:
//   1. Log the request.
//   2. Bind incoming JSON to a Product struct.
//   3. Validate product fields (name, price, SKU, etc).
//   4. Add product to the data layer.
//   5. Return success response.
//
func (p *Products) AddProduct(c *gin.Context) {
	p.l.Println("POST /product called")

	// Create empty product object for incoming JSON
	prod := &data.Product{}

	// Bind request body → struct
	if err := c.BindJSON(prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	// Validate product input using your custom Validate() method
	if err := prod.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	// Save product into memory storage
	data.AddProduct(prod)

	c.JSON(http.StatusCreated, gin.H{
		"message": "product added successfully",
		"data":    prod,
	})
}
