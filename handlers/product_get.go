package handlers

import (
	"github.com/VINAYAK777CODER/PRODUCT-API/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListProductsResponse represents the response for GET /products
//
// swagger:response productsResponse
type ListProductsResponse struct {
	// in: body
	Body []data.Product
}

// swagger:route GET /products products listProducts
//
// Get all products
//
// This endpoint returns the complete list of products available in the system.
//
// Responses:
//   200: productsResponse
//



// GetProducts handles: GET /product
// ---------------------------------------------------------------
// Purpose:
//   Fetch all products from the data layer and return them as JSON.
//
// Steps:
//   1. Log the request.
//   2. Retrieve product list from the data package.
//   3. Send the product list in JSON format.
//
func (p *Products) GetProducts(c *gin.Context) {
	p.l.Println("GET /product called")

	// fetch product list (from in-memory or DB)
	list := data.GetProductList()

	// return JSON response
	c.JSON(http.StatusOK, list)
}
