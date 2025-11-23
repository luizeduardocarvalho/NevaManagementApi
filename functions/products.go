package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
	"gorm.io/gorm"
)

// Request/Response Types

type CreateProductRequest struct {
	Name           string  `json:"name" validate:"required"`
	Description    string  `json:"description"`
	LocationID     uint    `json:"location_id" validate:"required"`
	Quantity       float64 `json:"quantity"`
	Formula        string  `json:"formula"`
	Unit           string  `json:"unit" validate:"required"`
	ExpirationDate string  `json:"expiration_date"`
	// LaboratoryID is not required in request - it's taken from authenticated user
}

type UpdateProductRequest struct {
	Name           string  `json:"name" validate:"required"`
	Description    string  `json:"description"`
	LocationID     uint    `json:"location_id" validate:"required"`
	Quantity       float64 `json:"quantity"`
	Formula        string  `json:"formula"`
	Unit           string  `json:"unit" validate:"required"`
	ExpirationDate string  `json:"expiration_date"`
	// LaboratoryID is not required in request - it's taken from authenticated user
}

type AddQuantityRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
}

type UseProductRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
	Unit     string  `json:"unit" validate:"required"`
	Notes    string  `json:"notes"`
}

type LocationResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductResponse struct {
	ID                               uint             `json:"id"`
	Name                             string           `json:"name"`
	Description                      string           `json:"description"`
	Formula                          string           `json:"formula"`
	Quantity                         float64          `json:"quantity"`
	Unit                             string           `json:"unit"`
	Location                         LocationResponse `json:"location"`
	LocationID                       uint             `json:"location_id,omitempty"`
	ExpirationDate                   string           `json:"expiration_date"`
	QuantityUsedInTheLastThreeMonths float64          `json:"quantity_used_in_the_last_three_months"`
	LaboratoryID                     uint             `json:"laboratory_id"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	NextPage   *int              `json:"nextPage"`
	TotalCount int64             `json:"totalCount"`
}

// Helper Functions

func getLaboratoryIDFromContext(r *http.Request) (uint, error) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		return 0, fmt.Errorf("authentication required")
	}

	if claims.LaboratoryID == 0 {
		return 0, fmt.Errorf("user does not belong to a laboratory")
	}

	return claims.LaboratoryID, nil
}

func parseExpirationDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	// Try ISO 8601 format first
	t, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr)
	if err == nil {
		return t, nil
	}

	// Try alternative ISO format
	t, err = time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
	if err == nil {
		return t, nil
	}

	// Try date only format
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}

func formatExpirationDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02T15:04:05.000Z")
}

func buildProductResponse(product models.Product) ProductResponse {
	// Calculate 3-month usage
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	var threeMonthUsage struct {
		TotalUsed float64
	}

	db := database.GetDB()
	db.Model(&models.ProductUsage{}).
		Select("COALESCE(SUM(quantity_used), 0) as total_used").
		Where("product_id = ? AND used_at >= ?", product.ID, threeMonthsAgo).
		Scan(&threeMonthUsage)

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Formula:     product.Formula,
		Quantity:    product.Quantity,
		Unit:        product.Unit,
		Location: LocationResponse{
			ID:          product.Location.ID,
			Name:        product.Location.Name,
			Description: product.Location.Description,
		},
		ExpirationDate:                   formatExpirationDate(product.ExpirationDate),
		QuantityUsedInTheLastThreeMonths: threeMonthUsage.TotalUsed,
		LaboratoryID:                     product.LaboratoryID,
	}
}

func buildDetailedProductResponse(product models.Product) ProductResponse {
	resp := buildProductResponse(product)
	resp.LocationID = product.LocationID
	return resp
}

// Handlers

// GET /api/products
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 9
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	db := database.GetDB()

	// Get total count
	var totalCount int64
	db.Model(&models.Product{}).Where("laboratory_id = ?", labID).Count(&totalCount)

	// Get products with pagination
	var products []models.Product
	offset := (page - 1) * pageSize
	result := db.Where("laboratory_id = ?", labID).
		Preload("Location").
		Offset(offset).
		Limit(pageSize).
		Find(&products)

	if result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch products"})
		return
	}

	// Build response
	productResponses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, buildProductResponse(product))
	}

	// Calculate next page
	var nextPage *int
	if int64(page*pageSize) < totalCount {
		np := page + 1
		nextPage = &np
	}

	response := ProductListResponse{
		Products:   productResponses,
		NextPage:   nextPage,
		TotalCount: totalCount,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// GET /api/products/low-stock
func GetLowStockProducts(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var products []models.Product

	// TODO: Implement proper low stock logic based on 3-month average usage
	// For now, consider products with quantity <= 10 as low stock
	result := db.Where("laboratory_id = ? AND quantity <= ?", labID, 10).
		Preload("Location").
		Find(&products)

	if result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch low stock products"})
		return
	}

	productResponses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, buildProductResponse(product))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, productResponses)
}

// GET /api/products/expired
func GetExpiredProducts(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var products []models.Product

	result := db.Where("laboratory_id = ? AND expiration_date < ? AND expiration_date IS NOT NULL", labID, time.Now()).
		Preload("Location").
		Find(&products)

	if result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch expired products"})
		return
	}

	productResponses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, buildProductResponse(product))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, productResponses)
}

// GET /api/products/:id
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Get product ID from URL path
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", productID, labID).
		Preload("Location").
		First(&product)

	if result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	response := buildDetailedProductResponse(product)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// POST /api/products
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	// Parse expiration date
	expirationDate, err := parseExpirationDate(req.ExpirationDate)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid expiration date format"})
		return
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
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create product"})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"message": "Product created successfully"})
}

// PUT /api/products/:id
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Get product ID from URL path
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	expirationDate, err := parseExpirationDate(req.ExpirationDate)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid expiration date format"})
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", productID, labID).First(&product)
	if result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.LocationID = req.LocationID
	product.Quantity = req.Quantity
	product.Formula = req.Formula
	product.Unit = req.Unit
	product.ExpirationDate = expirationDate

	result = db.Save(&product)
	if result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update product"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Product updated successfully"})
}

// POST /api/products/:id/add-quantity
func AddQuantityToProduct(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Get product ID from URL path
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	var req AddQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Quantity <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Quantity must be positive"})
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", productID, labID).First(&product)
	if result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	product.Quantity += req.Quantity
	result = db.Save(&product)

	if result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update product quantity"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Quantity added successfully"})
}

// POST /api/products/:id/use
func UseProduct(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "authentication required"})
		return
	}

	labID := claims.LaboratoryID
	if labID == 0 {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, map[string]string{"error": "user does not belong to a laboratory"})
		return
	}

	// Get product ID from URL path
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	var req UseProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Quantity <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Quantity must be positive"})
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", productID, labID).First(&product)
	if result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	if product.Quantity < req.Quantity {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Insufficient quantity available"})
		return
	}

	// Use transaction to ensure atomicity
	err = db.Transaction(func(tx *gorm.DB) error {
		// Deduct quantity from product
		product.Quantity -= req.Quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Track usage history
		usage := models.ProductUsage{
			ProductID:    uint(productID),
			UserID:       claims.UserID,
			QuantityUsed: req.Quantity,
			Unit:         req.Unit,
			UsedAt:       time.Now(),
			Notes:        req.Notes,
			LaboratoryID: labID,
		}

		if err := tx.Create(&usage).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to record product usage"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Product usage recorded successfully"})
}

// DELETE /api/products/:id
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Get product ID from URL path
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	db := database.GetDB()
	var product models.Product

	result := db.Where("id = ? AND laboratory_id = ?", productID, labID).First(&product)
	if result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	// Soft delete
	result = db.Delete(&product)
	if result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete product"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Product deleted successfully"})
}

// GET /api/products/:id/usage-history
func GetProductUsageHistory(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	db := database.GetDB()

	// Verify product exists and belongs to user's laboratory
	var product models.Product
	if err := db.Where("id = ? AND laboratory_id = ?", productID, labID).First(&product).Error; err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	// Get usage history
	var usages []models.ProductUsage
	if err := db.Where("product_id = ? AND laboratory_id = ?", productID, labID).
		Preload("User").
		Order("used_at DESC").
		Find(&usages).Error; err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch usage history"})
		return
	}

	render.JSON(w, r, usages)
}

// GET /api/products/:id/usage-stats
func GetProductUsageStats(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid product ID"})
		return
	}

	db := database.GetDB()

	// Verify product exists and belongs to user's laboratory
	var product models.Product
	if err := db.Where("id = ? AND laboratory_id = ?", productID, labID).First(&product).Error; err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Product not found"})
		return
	}

	// Calculate 3-month usage
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	var threeMonthUsage struct {
		TotalUsed float64
	}

	db.Model(&models.ProductUsage{}).
		Select("COALESCE(SUM(quantity_used), 0) as total_used").
		Where("product_id = ? AND laboratory_id = ? AND used_at >= ?", productID, labID, threeMonthsAgo).
		Scan(&threeMonthUsage)

	// Calculate 30-day usage
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var thirtyDayUsage struct {
		TotalUsed float64
	}

	db.Model(&models.ProductUsage{}).
		Select("COALESCE(SUM(quantity_used), 0) as total_used").
		Where("product_id = ? AND laboratory_id = ? AND used_at >= ?", productID, labID, thirtyDaysAgo).
		Scan(&thirtyDayUsage)

	// Calculate total usage ever
	var totalUsage struct {
		TotalUsed float64
	}

	db.Model(&models.ProductUsage{}).
		Select("COALESCE(SUM(quantity_used), 0) as total_used").
		Where("product_id = ? AND laboratory_id = ?", productID, labID).
		Scan(&totalUsage)

	stats := map[string]interface{}{
		"product_id":                productID,
		"product_name":              product.Name,
		"current_quantity":          product.Quantity,
		"unit":                      product.Unit,
		"usage_last_30_days":        thirtyDayUsage.TotalUsed,
		"usage_last_3_months":       threeMonthUsage.TotalUsed,
		"total_usage_all_time":      totalUsage.TotalUsed,
		"average_monthly_usage":     threeMonthUsage.TotalUsed / 3,
		"estimated_days_remaining":  calculateEstimatedDaysRemaining(product.Quantity, threeMonthUsage.TotalUsed),
	}

	render.JSON(w, r, stats)
}

func calculateEstimatedDaysRemaining(currentQty float64, threeMonthUsage float64) int {
	if threeMonthUsage <= 0 {
		return -1 // Unknown (no usage data)
	}

	dailyAverage := threeMonthUsage / 90 // 3 months ≈ 90 days
	if dailyAverage <= 0 {
		return -1
	}

	daysRemaining := currentQty / dailyAverage
	return int(daysRemaining)
}

// RegisterProductRoutes registers all product routes
func RegisterProductRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", GetAllProducts)
		r.Get("/low-stock", GetLowStockProducts)
		r.Get("/expired", GetExpiredProducts)
		r.Get("/{id}", GetProductByID)
		r.Get("/{id}/usage-history", GetProductUsageHistory)
		r.Get("/{id}/usage-stats", GetProductUsageStats)
		r.Post("/", CreateProduct)
		r.Put("/{id}", UpdateProduct)
		r.Post("/{id}/add-quantity", AddQuantityToProduct)
		r.Post("/{id}/use", UseProduct)
		r.Delete("/{id}", DeleteProduct)
	})
}
