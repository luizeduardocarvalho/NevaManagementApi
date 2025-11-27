package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type AddEquipmentRequest struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	PropertyNumber string `json:"property_number"`
	LocationID     uint   `json:"location_id" validate:"required"`
}

type GetEquipmentResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	PropertyNumber string `json:"property_number"`
	LocationID     uint   `json:"location_id"`
}

type GetDetailedEquipmentResponse struct {
	ID             uint                        `json:"id"`
	Name           string                      `json:"name"`
	Description    string                      `json:"description"`
	PropertyNumber string                      `json:"property_number"`
	Location       LocationResponse            `json:"location"`
	UsageList      []GetEquipmentUsageResponse `json:"usage_list"`
}

type GetEquipmentUsageResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	UserName    string `json:"user_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type EditEquipmentRequest struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	PropertyNumber string `json:"property_number"`
	LocationID     uint   `json:"location_id" validate:"required"`
}

func GetEquipments(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	// Pagination parameters
	page := 1
	pageSize := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	db := database.GetDB()

	// Get total count
	var total int64
	if err := db.Model(&models.Equipment{}).Where("laboratory_id = ?", claims.LaboratoryID).Count(&total).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to count equipment")
		return
	}

	// Get paginated equipment
	var equipment []models.Equipment
	offset := (page - 1) * pageSize
	if err := db.Where("laboratory_id = ?", claims.LaboratoryID).
		Preload("Location").
		Order("name").
		Offset(offset).
		Limit(pageSize).
		Find(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch equipment")
		return
	}

	var response []GetEquipmentResponse
	for _, eq := range equipment {
		response = append(response, GetEquipmentResponse{
			ID:             eq.ID,
			Name:           eq.Name,
			Description:    eq.Description,
			PropertyNumber: eq.PropertyNumber,
			LocationID:     eq.LocationID,
		})
	}

	render.JSON(w, r, map[string]interface{}{
		"equipment": response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AddEquipment(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	var req AddEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()

	// Validate location belongs to same laboratory
	var location models.Location
	if err := db.Where("id = ? AND laboratory_id = ?", req.LocationID, claims.LaboratoryID).First(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid location")
		return
	}

	equipment := models.Equipment{
		Name:           req.Name,
		Description:    req.Description,
		PropertyNumber: req.PropertyNumber,
		LocationID:     req.LocationID,
		LaboratoryID:   claims.LaboratoryID,
	}

	if err := db.Create(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to create equipment")
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"message":   fmt.Sprintf("%s was created successfully", req.Name),
		"equipment": equipment,
	})
}

func GetDetailedEquipment(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	equipmentID := chi.URLParam(r, "id")
	if equipmentID == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "equipment ID is required")
		return
	}

	eqID, err := strconv.ParseUint(equipmentID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid equipment ID")
		return
	}

	db := database.GetDB()
	var equipment models.Equipment

	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).
		Preload("Location").
		Preload("EquipmentUsages").
		Preload("EquipmentUsages.User").
		First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	var usageList []GetEquipmentUsageResponse
	for _, usage := range equipment.EquipmentUsages {
		userName := usage.User.FirstName + " " + usage.User.LastName
		usageList = append(usageList, GetEquipmentUsageResponse{
			ID:          usage.ID,
			UserID:      usage.UserID,
			UserName:    userName,
			StartDate:   usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:     usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
			Description: usage.Description,
		})
	}

	response := GetDetailedEquipmentResponse{
		ID:             equipment.ID,
		Name:           equipment.Name,
		Description:    equipment.Description,
		PropertyNumber: equipment.PropertyNumber,
		Location: LocationResponse{
			ID:          equipment.Location.ID,
			Name:        equipment.Location.Name,
			Description: equipment.Location.Description,
		},
		UsageList: usageList,
	}

	render.JSON(w, r, response)
}

func EditEquipment(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	equipmentID := chi.URLParam(r, "id")
	if equipmentID == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "equipment ID is required")
		return
	}

	eqID, err := strconv.ParseUint(equipmentID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid equipment ID")
		return
	}

	var req EditEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()

	// Validate location belongs to same laboratory
	var location models.Location
	if err := db.Where("id = ? AND laboratory_id = ?", req.LocationID, claims.LaboratoryID).First(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid location")
		return
	}

	var equipment models.Equipment
	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	equipment.Name = req.Name
	equipment.Description = req.Description
	equipment.PropertyNumber = req.PropertyNumber
	equipment.LocationID = req.LocationID

	if err := db.Save(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to update equipment")
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"message":   fmt.Sprintf("%s was updated successfully", req.Name),
		"equipment": equipment,
	})
}

func DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	equipmentID := chi.URLParam(r, "id")
	if equipmentID == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "equipment ID is required")
		return
	}

	eqID, err := strconv.ParseUint(equipmentID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid equipment ID")
		return
	}

	db := database.GetDB()
	var equipment models.Equipment

	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	// Soft delete
	if err := db.Delete(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to delete equipment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterEquipmentRoutes(r chi.Router) {
	r.Route("/equipment", func(r chi.Router) {
		r.Get("/", GetEquipments)
		r.Post("/", AddEquipment)
		r.Get("/{id}", GetDetailedEquipment)
		r.Put("/{id}", EditEquipment)
		r.Delete("/{id}", DeleteEquipment)

		// Equipment usage routes
		r.Post("/{id}/use", UseEquipment)
		r.Get("/{id}/usage-history", GetEquipmentUsageHistory)
		r.Get("/{id}/calendar", GetEquipmentUsageCalendar)
		r.Post("/{id}/check-overlap", CheckOverlap)
	})
}
