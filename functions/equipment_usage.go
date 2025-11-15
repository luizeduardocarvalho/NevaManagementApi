package functions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type UseEquipmentRequest struct {
	EquipmentID  uint   `json:"equipment_id" validate:"required"`
	ResearcherID uint   `json:"researcher_id" validate:"required"`
	Description  string `json:"description"`
	StartDate    string `json:"start_date" validate:"required"`
	EndDate      string `json:"end_date" validate:"required"`
}

type EquipmentUsageCalendarResponse struct {
	ID             uint   `json:"id"`
	EquipmentID    uint   `json:"equipment_id"`
	ResearcherID   uint   `json:"researcher_id"`
	ResearcherName string `json:"researcher_name"`
	Description    string `json:"description"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}

func GetEquipmentUsageCalendar(w http.ResponseWriter, r *http.Request) {
	equipmentID := r.URL.Query().Get("id")
	laboratoryID := r.URL.Query().Get("laboratoryId")

	if equipmentID == "" || laboratoryID == "" {
		http.Error(w, "Equipment ID and Laboratory ID are required", http.StatusBadRequest)
		return
	}

	eqID, err := strconv.ParseUint(equipmentID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		return
	}

	labID, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid laboratory ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var usages []models.EquipmentUsage

	result := db.Where("equipment_id = ? AND laboratory_id = ?", eqID, labID).
		Preload("Researcher").
		Find(&usages)

	if result.Error != nil {
		http.Error(w, "Failed to fetch equipment usage calendar", http.StatusInternalServerError)
		return
	}

	var response []EquipmentUsageCalendarResponse
	for _, usage := range usages {
		response = append(response, EquipmentUsageCalendarResponse{
			ID:             usage.ID,
			EquipmentID:    usage.EquipmentID,
			ResearcherID:   usage.ResearcherID,
			ResearcherName: usage.Researcher.Name,
			Description:    usage.Description,
			StartDate:      usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:        usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetEquipmentUsageHistory(w http.ResponseWriter, r *http.Request) {
	equipmentID := r.URL.Query().Get("id")
	laboratoryID := r.URL.Query().Get("laboratoryId")

	if equipmentID == "" || laboratoryID == "" {
		http.Error(w, "Equipment ID and Laboratory ID are required", http.StatusBadRequest)
		return
	}

	eqID, err := strconv.ParseUint(equipmentID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		return
	}

	labID, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid laboratory ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var usages []models.EquipmentUsage

	// Get past usage (end_date < now), ordered by end_date desc
	result := db.Where("equipment_id = ? AND laboratory_id = ? AND end_date < ?", eqID, labID, time.Now()).
		Preload("Researcher").
		Order("end_date desc").
		Find(&usages)

	if result.Error != nil {
		http.Error(w, "Failed to fetch equipment usage history", http.StatusInternalServerError)
		return
	}

	var response []EquipmentUsageCalendarResponse
	for _, usage := range usages {
		response = append(response, EquipmentUsageCalendarResponse{
			ID:             usage.ID,
			EquipmentID:    usage.EquipmentID,
			ResearcherID:   usage.ResearcherID,
			ResearcherName: usage.Researcher.Name,
			Description:    usage.Description,
			StartDate:      usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:        usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UseEquipment(w http.ResponseWriter, r *http.Request) {
	var req UseEquipmentRequest
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

	// Parse dates
	startDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.StartDate)
	if err != nil {
		startDate, err = time.Parse("2006-01-02T15:04:05Z", req.StartDate)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	}

	endDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EndDate)
	if err != nil {
		endDate, err = time.Parse("2006-01-02T15:04:05Z", req.EndDate)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	}

	// Validate dates
	if endDate.Before(startDate) {
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	// Check if equipment exists and belongs to the laboratory
	var equipment models.Equipment
	result := db.Where("id = ? AND laboratory_id = ?", req.EquipmentID, labID).First(&equipment)
	if result.Error != nil {
		http.Error(w, "Equipment not found", http.StatusNotFound)
		return
	}

	// Check if researcher exists and belongs to the laboratory
	var researcher models.Researcher
	result = db.Where("id = ? AND laboratory_id = ?", req.ResearcherID, labID).First(&researcher)
	if result.Error != nil {
		http.Error(w, "Researcher not found", http.StatusNotFound)
		return
	}

	// Check for conflicts (overlapping usage periods)
	var conflictingUsage models.EquipmentUsage
	result = db.Where(`equipment_id = ? AND laboratory_id = ? AND 
		((start_date <= ? AND end_date >= ?) OR 
		 (start_date <= ? AND end_date >= ?) OR 
		 (start_date >= ? AND end_date <= ?))`,
		req.EquipmentID, labID,
		startDate, startDate,
		endDate, endDate,
		startDate, endDate).
		First(&conflictingUsage)

	if result.RowsAffected > 0 {
		http.Error(w, "Equipment is already in use during the specified time period", http.StatusConflict)
		return
	}

	// Create equipment usage record
	usage := models.EquipmentUsage{
		EquipmentID:  req.EquipmentID,
		ResearcherID: req.ResearcherID,
		Description:  req.Description,
		StartDate:    startDate,
		EndDate:      endDate,
		LaboratoryID: labID,
	}

	result = db.Create(&usage)
	if result.Error != nil {
		http.Error(w, "Failed to record equipment usage", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Equipment was used successfully."))
}

func GetEquipmentUsage(w http.ResponseWriter, r *http.Request) {
	equipmentID := r.URL.Query().Get("equipmentId")
	laboratoryID := r.URL.Query().Get("laboratoryId")
	page := r.URL.Query().Get("page")

	if equipmentID == "" || laboratoryID == "" {
		http.Error(w, "Equipment ID and Laboratory ID are required", http.StatusBadRequest)
		return
	}

	eqID, err := strconv.ParseUint(equipmentID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
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
	var usages []models.EquipmentUsage

	offset := (pageNum - 1) * 10
	result := db.Where("equipment_id = ? AND laboratory_id = ?", eqID, labID).
		Preload("Researcher").
		Order("start_date desc").
		Offset(offset).
		Limit(10).
		Find(&usages)

	if result.Error != nil {
		http.Error(w, "Failed to fetch equipment usage", http.StatusInternalServerError)
		return
	}

	var response []EquipmentUsageCalendarResponse
	for _, usage := range usages {
		response = append(response, EquipmentUsageCalendarResponse{
			ID:             usage.ID,
			EquipmentID:    usage.EquipmentID,
			ResearcherID:   usage.ResearcherID,
			ResearcherName: usage.Researcher.Name,
			Description:    usage.Description,
			StartDate:      usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:        usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterEquipmentUsageRoutes(r chi.Router) {
	r.Route("/equipmentusage", func(r chi.Router) {
		r.Get("/GetEquipmentUsageCalendar", GetEquipmentUsageCalendar)
		r.Get("/GetEquipmentUsageHistory", GetEquipmentUsageHistory)
		r.Post("/UseEquipment", UseEquipment)
		r.Get("/GetEquipmentUsage", GetEquipmentUsage)
	})
}