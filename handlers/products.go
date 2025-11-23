package handlers

import (
	"WORKING-GO/data" // import our data package (contains Product list + JSON code)
	"log"             // used for printing logs
	"net/http"        // used for creating HTTP servers and handlers
	"regexp"
	"strconv"
)

// Products is a struct that holds a logger.
// We keep a struct so we can attach methods to it.
type Products struct {
	l *log.Logger // logger used to print messages in server
}

// NewProducts creates a new Products handler and returns it.
// We pass the logger so this handler can use it.
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP makes Products satisfy the http.Handler interface.
// This function runs automatically when a request comes to this handler.
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// ---------------------------
	// HANDLE GET REQUEST
	// If client sends GET /products → return all products
	// ---------------------------
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// ---------------------------
	// HANDLE POST REQUEST
	// If client sends POST /products → add a new product
	// ---------------------------
	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}

	// ---------------------------
	// HANDLE PUT REQUEST
	// PUT always requires an ID → PUT /products/3
	// ---------------------------
	if r.Method == http.MethodPut {

		// Create a regex to extract number (ID) from URL
		// Example URL: /products/10 → captures "10"
		reg := regexp.MustCompile(`/([0-9]+)`)

		// Find first match in the URL path
		// Example output: ["/10", "10"]
		g := reg.FindStringSubmatch(r.URL.Path)

		// If match is not exactly 2 items → URL is invalid
		// (means ID not found)
		if len(g) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		// Extract ID string from match
		idString := g[1]

		// Convert ID string → integer
		id, _ := strconv.Atoi(idString)

		// Log the ID (just for debugging)
		p.l.Println("Got ID:", id)

		// Call update handler to process the PUT request
		p.updateProducts(id, rw, r)
		return
	}

	// ---------------------------
	// HANDLE ALL OTHER METHODS
	// If user sends DELETE / PATCH / HEAD / OPTIONS etc.
	// → We do not support them, so return 405
	// ---------------------------
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts sends all product data in JSON format
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {

	// Get the list of products from data package
	lp := data.GetProductList()

	// Convert product list to JSON and write it to response
	err := lp.ToJSON(rw)

	// If something went wrong while converting to JSON
	if err != nil {
		http.Error(rw, "unable to convert data to json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Product POST request received")

	// create an empty Product struct
	prod := &data.Product{}

	// decode the JSON body into prod
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)

	// add the product to in-memory list
	data.AddProduct(prod)

}

// updateProducts handles PUT /products/{id}
// It receives ID extracted from URL, reads new JSON data,
// and asks the data layer to update the product.
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {

	// Just log that PUT request started
	p.l.Println("Handle PUT request")

	// Create an empty product struct where new data will be filled
	prod := &data.Product{}

	// STEP 1: Decode JSON body → fill "prod"
	// Example: {"name":"Tea", "price":50}
	err := prod.FromJSON(r.Body)
	if err != nil {
		// If client sends invalid JSON → return 400 Bad Request
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	// STEP 2: Ask the DATA LAYER to update this product
	// This function will:
	//   - Search product by ID
	//   - Replace old product with new one
	err = data.UpdateProduct(id, prod)

	// If ID does not exist → send 404 Not Found
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	// If any OTHER problem occurred → send 500
	if err != nil {
		http.Error(rw, "Unable to update product", http.StatusInternalServerError)
		return
	}

	// If everything is OK, no need to send anything (or you can send OK)
}
