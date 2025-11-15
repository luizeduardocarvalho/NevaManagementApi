package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type CreateProductRequest struct {
	Name           string  `json:"name" validate:"required"`
	Description    string  `json:"description"`
	LocationID     uint    `json:"location_id" validate:"required"`
	Quantity       float64 `json:"quantity"`
	Formula        string  `json:"formula"`
	Unit           string  `json:"unit" validate:"required"`
	ExpirationDate string  `json:"expiration_date"`
}

type GetDetailedProductResponse struct {
	ID                                  uint    `json:"id"`
	Name                               string  `json:"name"`
	Description                        string  `json:"description"`
	Quantity                           float64 `json:"quantity"`
	QuantityUsedInTheLastThreeMonths   float64 `json:"quantity_used_in_the_last_three_months"`
	Formula                            string  `json:"formula"`
	Unit                               string  `json:"unit"`
	ExpirationDate                     string  `json:"expiration_date"`
	Location                           LocationResponse `json:"location"`
}

type GetProductResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Quantity       float64 `json:"quantity"`
	Formula        string  `json:"formula"`
	Unit           string  `json:"unit"`
	ExpirationDate string  `json:"expiration_date"`
	LocationID     uint    `json:"location_id"`
}

type LocationResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	SubLocationID *uint  `json:"sub_location_id"`
}

type AddQuantityToProductRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required"`
}

type UseProductRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required"`
	Unit      string  `json:"unit" validate:"required"`
}

type EditProductRequest struct {
	ID             uint    `json:"id" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	Description    string  `json:"description"`
	LocationID     uint    `json:"location_id" validate:"required"`
	Quantity       float64 `json:"quantity"`
	Formula        string  `json:"formula"`
	Unit           string  `json:"unit" validate:"required"`
	ExpirationDate string  `json:"expiration_date"`
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	laboratoryID := r.URL.Query().Get("laboratoryId")
	page := r.URL.Query().Get("page")

	if laboratoryID == "" {
		http.Error(w, "Laboratory ID is required", http.StatusBadRequest)
		return
	}

	labID, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid laboratory ID", http.StatusBadRequest)
		return
	}

	pageNum := 1
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	db := database.GetDB()
	var products []models.Product

	offset := (pageNum - 1) * 10
	result := db.Where("laboratory_id = ?", labID).
		Preload("Location").
		Offset(offset).
		Limit(10).
		Find(&products)

	if result.Error != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	var response []GetProductResponse
	for _, product := range products {
		response = append(response, GetProductResponse{
			ID:             product.ID,
			Name:           product.Name,
			Description:    product.Description,
			Quantity:       product.Quantity,
			Formula:        product.Formula,
			Unit:           product.Unit,
			ExpirationDate: product.ExpirationDate.Format("2006-01-02T15:04:05Z07:00"),
			LocationID:     product.LocationID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailedProductByID(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	laboratoryID := r.URL.Query().Get("laboratoryId")

	if productID == "" || laboratoryID == "" {
		http.Error(w, "Product ID and Laboratory ID are required", http.StatusBadRequest)
		return
	}

	prodID, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	labID, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid laboratory ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", prodID, labID).
		Preload("Location").
		First(&product)

	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := GetDetailedProductResponse{
		ID:                                product.ID,
		Name:                             product.Name,
		Description:                      product.Description,
		Quantity:                         product.Quantity,
		QuantityUsedInTheLastThreeMonths: 0, // TODO: Calculate from usage data
		Formula:                          product.Formula,
		Unit:                             product.Unit,
		ExpirationDate:                   product.ExpirationDate.Format("2006-01-02T15:04:05Z07:00"),
		Location: LocationResponse{
			ID:            product.Location.ID,
			Name:          product.Location.Name,
			Description:   product.Location.Description,
			SubLocationID: product.Location.SubLocationID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	laboratoryID := r.URL.Query().Get("laboratoryId")

	if productID == "" || laboratoryID == "" {
		http.Error(w, "Product ID and Laboratory ID are required", http.StatusBadRequest)
		return
	}

	prodID, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	labID, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid laboratory ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", prodID, labID).First(&product)

	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := GetProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Description:    product.Description,
		Quantity:       product.Quantity,
		Formula:        product.Formula,
		Unit:           product.Unit,
		ExpirationDate: product.ExpirationDate.Format("2006-01-02T15:04:05Z07:00"),
		LocationID:     product.LocationID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get laboratory ID from JWT claims (middleware should have set this)
	laboratoryID := r.Context().Value("laboratory_id")
	if laboratoryID == nil {
		http.Error(w, "Laboratory ID not found in token", http.StatusUnauthorized)
		return
	}

	labID, ok := laboratoryID.(uint)
	if !ok {
		http.Error(w, "Invalid laboratory ID in token", http.StatusUnauthorized)
		return
	}

	var expirationDate time.Time
	if req.ExpirationDate != "" {
		var err error
		expirationDate, err = time.Parse("2006-01-02T15:04:05Z07:00", req.ExpirationDate)
		if err != nil {
			expirationDate, err = time.Parse("2006-01-02", req.ExpirationDate)
			if err != nil {
				http.Error(w, "Invalid expiration date format", http.StatusBadRequest)
				return
			}
		}
	}

	product := models.Product{
		Name:           req.Name,
		Description:    req.Description,
		LocationID:     req.LocationID,
		Quantity:       req.Quantity,
		Formula:        req.Formula,
		Unit:           req.Unit,
		ExpirationDate: expirationDate,
		LaboratoryID:   labID,
	}

	db := database.GetDB()
	result := db.Create(&product)

	if result.Error != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Successfully created %s.", req.Name)
}

func AddQuantityToProduct(w http.ResponseWriter, r *http.Request) {
	var req AddQuantityToProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get laboratory ID from JWT claims
	laboratoryID := r.Context().Value("laboratory_id")
	if laboratoryID == nil {
		http.Error(w, "Laboratory ID not found in token", http.StatusUnauthorized)
		return
	}

	labID, ok := laboratoryID.(uint)
	if !ok {
		http.Error(w, "Invalid laboratory ID in token", http.StatusUnauthorized)
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", req.ProductID, labID).First(&product)
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product.Quantity += req.Quantity
	result = db.Save(&product)

	if result.Error != nil {
		http.Error(w, "Failed to update product quantity", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully added %.2f to product.", req.Quantity)
}

func UseProduct(w http.ResponseWriter, r *http.Request) {
	var req UseProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get laboratory ID from JWT claims
	laboratoryID := r.Context().Value("laboratory_id")
	if laboratoryID == nil {
		http.Error(w, "Laboratory ID not found in token", http.StatusUnauthorized)
		return
	}

	labID, ok := laboratoryID.(uint)
	if !ok {
		http.Error(w, "Invalid laboratory ID in token", http.StatusUnauthorized)
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", req.ProductID, labID).First(&product)
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if product.Quantity < req.Quantity {
		http.Error(w, "Insufficient quantity available", http.StatusBadRequest)
		return
	}

	product.Quantity -= req.Quantity
	result = db.Save(&product)

	if result.Error != nil {
		http.Error(w, "Failed to update product quantity", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully used %.2f%s.", req.Quantity, req.Unit)
}

func EditProduct(w http.ResponseWriter, r *http.Request) {
	var req EditProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get laboratory ID from JWT claims
	laboratoryID := r.Context().Value("laboratory_id")
	if laboratoryID == nil {
		http.Error(w, "Laboratory ID not found in token", http.StatusUnauthorized)
		return
	}

	labID, ok := laboratoryID.(uint)
	if !ok {
		http.Error(w, "Invalid laboratory ID in token", http.StatusUnauthorized)
		return
	}

	var expirationDate time.Time
	if req.ExpirationDate != "" {
		var err error
		expirationDate, err = time.Parse("2006-01-02T15:04:05Z07:00", req.ExpirationDate)
		if err != nil {
			expirationDate, err = time.Parse("2006-01-02", req.ExpirationDate)
			if err != nil {
				http.Error(w, "Invalid expiration date format", http.StatusBadRequest)
				return
			}
		}
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", req.ID, labID).First(&product)
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.LocationID = req.LocationID
	product.Quantity = req.Quantity
	product.Formula = req.Formula
	product.Unit = req.Unit
	product.ExpirationDate = expirationDate

	result = db.Save(&product)
	if result.Error != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully edited %s.", req.Name)
}

func GetLowInStockProducts(w http.ResponseWriter, r *http.Request) {
	// Get laboratory ID from JWT claims
	laboratoryID := r.Context().Value("laboratory_id")
	if laboratoryID == nil {
		http.Error(w, "Laboratory ID not found in token", http.StatusUnauthorized)
		return
	}

	labID, ok := laboratoryID.(uint)
	if !ok {
		http.Error(w, "Invalid laboratory ID in token", http.StatusUnauthorized)
		return
	}

	db := database.GetDB()
	var products []models.Product

	// Consider products with quantity <= 10 as low in stock
	result := db.Where("laboratory_id = ? AND quantity <= ?", labID, 10).
		Preload("Location").
		Find(&products)

	if result.Error != nil {
		http.Error(w, "Failed to fetch low stock products", http.StatusInternalServerError)
		return
	}

	var response []GetProductResponse
	for _, product := range products {
		response = append(response, GetProductResponse{
			ID:             product.ID,
			Name:           product.Name,
			Description:    product.Description,
			Quantity:       product.Quantity,
			Formula:        product.Formula,
			Unit:           product.Unit,
			ExpirationDate: product.ExpirationDate.Format("2006-01-02T15:04:05Z07:00"),
			LocationID:     product.LocationID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterProductRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/GetAll", GetAllProducts)
		r.Get("/GetDetailedProductById", GetDetailedProductByID)
		r.Get("/GetProductById", GetProductByID)
		r.Post("/Create", CreateProduct)
		r.Patch("/AddQuantity", AddQuantityToProduct)
		r.Patch("/UseProduct", UseProduct)
		r.Patch("/EditProduct", EditProduct)
		r.Get("/GetLowInStockProducts", GetLowInStockProducts)
	})
}