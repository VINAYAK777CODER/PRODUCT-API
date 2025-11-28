package handlers

import (
	"WORKING-GO/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
