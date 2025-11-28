package data

import (
	"encoding/json" // Used to convert Go data <-> JSON
	"fmt"
	"io" // io.Reader + io.Writer (used by HTTP request & response)
	"regexp"
	"time" // Used for timestamps (CreatedOn, UpdatedOn, etc.)

	"github.com/go-playground/validator/v10"
)


/*
------------------------------------------------------
GLOBAL VARIABLES (package-level)
------------------------------------------------------
We define these here because:

1. skuRegex is compiled ONLY ONCE → faster performance.
2. validate is created ONCE → avoids re-registering rules.
3. init() will prepare everything BEFORE main() runs.
*/
var (
	validate = validator.New() // Validator instance used by all products
	skuRegex = regexp.MustCompile(`^[a-z]+-[0-9]+-[a-z]+$`) // Precompiled regex
)

/*
------------------------------------------------------
init() FUNCTION (special function)
------------------------------------------------------
REASON: init() is automatically executed by Go *before* main().

Why use init() here?
✔ Register custom validators only once
✔ Avoid repeated setup in Product.Validate()
✔ Cleaner and faster

You NEVER call init() yourself.
Go runtime always calls it before main.
*/
func init() {

	// Register our custom validation rule named "sku"
	// Example: validate:"sku"
	err := validate.RegisterValidation("sku", validateSKU)

	// If registration fails (ex: duplicate tag), print an error
	if err != nil {
		fmt.Println("Error registering SKU validation:", err.Error())
	}
}

/*
------------------------------------------------------
Product Struct
------------------------------------------------------
The 'SKU' field uses custom validation tag: `sku`
So the validator will run validateSKU() on this field.
*/
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"` // "-" means DO NOT show in JSON
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

/*
------------------------------------------------------
Product.Validate()
------------------------------------------------------
This method calls validator.Struct(p)
which checks:

✔ Built-in rules (required, gt=0)
✔ Custom rule ("sku")

REASON:
We separated registration into init(),
so this method stays clean and fast.
*/
func (p *Product) Validate() error {
	return validate.Struct(p)
}

/*
------------------------------------------------------
validateSKU()
------------------------------------------------------
This is our custom validator function.

Purpose:
Validate that SKU matches format:
   letters - numbers - letters
   example: abc-123-xyz

Steps:
1. Extract the field value as string
2. Match it against precompiled regex
3. Return true (valid) or false (invalid)
*/
func validateSKU(fl validator.FieldLevel) bool {

	// Field value (e.g., "abc-123-xyz")
	sku := fl.Field().String()

	// Match with strict regex pattern
	return skuRegex.MatchString(sku)
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

func DeleteProduct(id int) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return ErrProductNotFound
	}

	// Remove product at index `pos`
	ProductList = append(ProductList[:pos], ProductList[pos+1:]...)

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
		SKU:         "abc-323-abc",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd-343-abd",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
