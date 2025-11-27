package functions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type UseEquipmentRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	Description string `json:"description"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
}

type CheckOverlapRequest struct {
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type EquipmentUsageCalendarResponse struct {
	ID          uint   `json:"id"`
	EquipmentID uint   `json:"equipment_id"`
	UserID      uint   `json:"user_id"`
	UserName    string `json:"user_name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

func GetEquipmentUsageCalendar(w http.ResponseWriter, r *http.Request) {
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

	// Verify equipment belongs to user's laboratory
	var equipment models.Equipment
	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	var usages []models.EquipmentUsage
	if err := db.Where("equipment_id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).
		Preload("User").
		Order("start_date").
		Find(&usages).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch equipment usage calendar")
		return
	}

	var response []EquipmentUsageCalendarResponse
	for _, usage := range usages {
		userName := usage.User.FirstName + " " + usage.User.LastName
		response = append(response, EquipmentUsageCalendarResponse{
			ID:          usage.ID,
			EquipmentID: usage.EquipmentID,
			UserID:      usage.UserID,
			UserName:    userName,
			Description: usage.Description,
			StartDate:   usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:     usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	render.JSON(w, r, response)
}

func GetEquipmentUsageHistory(w http.ResponseWriter, r *http.Request) {
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

	// Verify equipment belongs to user's laboratory
	var equipment models.Equipment
	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	var usages []models.EquipmentUsage

	// Get past usage (end_date < now), ordered by end_date desc
	if err := db.Where("equipment_id = ? AND laboratory_id = ? AND end_date < ?", eqID, claims.LaboratoryID, time.Now()).
		Preload("User").
		Order("end_date desc").
		Find(&usages).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch equipment usage history")
		return
	}

	var response []EquipmentUsageCalendarResponse
	for _, usage := range usages {
		userName := usage.User.FirstName + " " + usage.User.LastName
		response = append(response, EquipmentUsageCalendarResponse{
			ID:          usage.ID,
			EquipmentID: usage.EquipmentID,
			UserID:      usage.UserID,
			UserName:    userName,
			Description: usage.Description,
			StartDate:   usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:     usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	render.JSON(w, r, response)
}

func UseEquipment(w http.ResponseWriter, r *http.Request) {
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

	var req UseEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.StartDate)
	if err != nil {
		startDate, err = time.Parse("2006-01-02T15:04:05Z", req.StartDate)
		if err != nil {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid start date format")
			return
		}
	}

	endDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EndDate)
	if err != nil {
		endDate, err = time.Parse("2006-01-02T15:04:05Z", req.EndDate)
		if err != nil {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid end date format")
			return
		}
	}

	// Validate dates
	if endDate.Before(startDate) {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "end date must be after start date")
		return
	}

	db := database.GetDB()

	// Check if equipment exists and belongs to the laboratory
	var equipment models.Equipment
	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	// Check if user exists and belongs to the laboratory
	var user models.User
	if err := db.Where("id = ? AND laboratory_id = ?", req.UserID, claims.LaboratoryID).First(&user).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "user not found")
		return
	}

	// Check for conflicts (overlapping usage periods)
	var conflictingUsage models.EquipmentUsage
	result := db.Where(`equipment_id = ? AND laboratory_id = ? AND
		((start_date <= ? AND end_date >= ?) OR
		 (start_date <= ? AND end_date >= ?) OR
		 (start_date >= ? AND end_date <= ?))`,
		eqID, claims.LaboratoryID,
		startDate, startDate,
		endDate, endDate,
		startDate, endDate).
		First(&conflictingUsage)

	if result.RowsAffected > 0 {
		middleware.ErrorResponse(w, r, http.StatusConflict, "equipment is already in use during the specified time period")
		return
	}

	// Create equipment usage record
	usage := models.EquipmentUsage{
		EquipmentID:  uint(eqID),
		UserID:       req.UserID,
		Description:  req.Description,
		StartDate:    startDate,
		EndDate:      endDate,
		LaboratoryID: claims.LaboratoryID,
	}

	if err := db.Create(&usage).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to record equipment usage")
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"message": "equipment reservation created successfully",
		"usage":   usage,
	})
}

func CheckOverlap(w http.ResponseWriter, r *http.Request) {
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

	var req CheckOverlapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.StartDate)
	if err != nil {
		startDate, err = time.Parse("2006-01-02T15:04:05Z", req.StartDate)
		if err != nil {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid start date format")
			return
		}
	}

	endDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EndDate)
	if err != nil {
		endDate, err = time.Parse("2006-01-02T15:04:05Z", req.EndDate)
		if err != nil {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid end date format")
			return
		}
	}

	// Validate dates
	if endDate.Before(startDate) {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "end date must be after start date")
		return
	}

	db := database.GetDB()

	// Verify equipment belongs to user's laboratory
	var equipment models.Equipment
	if err := db.Where("id = ? AND laboratory_id = ?", eqID, claims.LaboratoryID).First(&equipment).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "equipment not found")
		return
	}

	// Check for conflicts (overlapping usage periods)
	var conflictingUsages []models.EquipmentUsage
	if err := db.Where(`equipment_id = ? AND laboratory_id = ? AND
		((start_date <= ? AND end_date >= ?) OR
		 (start_date <= ? AND end_date >= ?) OR
		 (start_date >= ? AND end_date <= ?))`,
		eqID, claims.LaboratoryID,
		startDate, startDate,
		endDate, endDate,
		startDate, endDate).
		Preload("User").
		Find(&conflictingUsages).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to check overlap")
		return
	}

	if len(conflictingUsages) > 0 {
		var conflicts []map[string]interface{}
		for _, usage := range conflictingUsages {
			userName := usage.User.FirstName + " " + usage.User.LastName
			conflicts = append(conflicts, map[string]interface{}{
				"id":          usage.ID,
				"user_id":     usage.UserID,
				"user_name":   userName,
				"start_date":  usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
				"end_date":    usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
				"description": usage.Description,
			})
		}

		render.JSON(w, r, map[string]interface{}{
			"has_conflict": true,
			"conflicts":    conflicts,
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"has_conflict": false,
		"message":      "no conflicts found",
	})
}
