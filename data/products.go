package data

import (
	"encoding/json" // Used to convert Go data <-> JSON
	"fmt"
	"io"   // io.Reader + io.Writer (used by HTTP request & response)
	"time" // Used for timestamps (CreatedOn, UpdatedOn, etc.)
)

// ----------------------
// Product STRUCT
// ----------------------
//
// Represents ONE product.
// json:"fieldName" → how the field appears in JSON.
// json:"_"         → this field will NOT be sent in JSON output.
type Product struct {
	ID          int     `json:"id"`          // Included in JSON
	Name        string  `json:"name"`        // Included in JSON
	Description string  `json:"description"` // Included in JSON
	Price       float32 `json:"price"`       // Included in JSON

	// These fields are INTERNAL only → not shown in JSON response
	SKU       string `json:"_"`
	CreatedOn string `json:"_"`
	UpdatedOn string `json:"_"`
	DeletedOn string `json:"_"`
}

// ----------------------
// Products TYPE
// ----------------------
//
// Products is a NEW type made from []*Product (slice of product pointers).
//
// Why create this new type?
// 1. Cleaner code (we use Products instead of []*Product everywhere)
// 2. We can add methods on it (like ToJSON)
// 3. Makes the code easy to understand
type Products []*Product

// ----------------------
// FromJSON (Decode)
// ----------------------
//
// Converts JSON coming from REQUEST body → into a Product struct.
//
// Why receiver is (p *Product)?
// → Because JSON of ONE product should fill ONE Product.
// → "p" will store the decoded JSON data.
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ----------------------
// ToJSON (Encode)
// ----------------------
//
// Converts a Products list → JSON → writes directly to the output.
//
// Why receiver is (p *Products)?
// → Because we are encoding MANY products (slice)
// → And pointer avoids copying large data.
//
// Why use Encoder instead of Marshal?
// → Encoder writes directly to response (streaming)
// → Faster + less memory
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p) // writes JSON directly to 'w'
}

// ----------------------
// GetProductList
// ----------------------
//
// Returns our in-memory product slice (Products type).
func GetProductList() Products {
	return ProductList
}

// ----------------------
// AddProduct
// ----------------------
//
// Adds a new product to the ProductList.
// Auto-increases ID.
func AddProduct(p *Product) {
	p.ID = getNextId()
	ProductList = append(ProductList, p)
}

// ----------------------
// getNextId
// ----------------------
//
// Gets the last product ID and returns next ID.
func getNextId() int {
	lastProduct := ProductList[len(ProductList)-1]
	return lastProduct.ID + 1
}

// UpdateProduct updates the product that matches the given ID.
// This function belongs to "data" layer because updating records
// is a DATA responsibility (not handler responsibility).
func UpdateProduct(id int, p *Product) error {

	// STEP 1: Find the product and get its index
	// findProduct returns:
	//   existingProduct, indexInSlice, error
	_, pos, err := findProduct(id)
	if err != nil {
		// Means ID does not exist → return the same error
		return err
	}

	// STEP 2: Set the ID (to ensure user cannot change ID through JSON)
	p.ID = id

	// STEP 3: Replace old product with new updated product
	ProductList[pos] = p

	return nil
}

// Common error used when ID is not found
var ErrProductNotFound = fmt.Errorf("Product not found")

// findProduct searches for product with given ID.
// It loops through ProductList and returns:
//
//	product, index, nil → if found
//	nil, -1, ErrProductNotFound  → if not found
func findProduct(id int) (*Product, int, error) {
	for i, p := range ProductList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

// ----------------------
// ProductList
// ----------------------
//
// Our fake database stored in memory.
// In real apps, this comes from MySQL / MongoDB / PostgreSQL.
var ProductList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
