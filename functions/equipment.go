package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
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
	ID             uint                     `json:"id"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	PropertyNumber string                   `json:"property_number"`
	Location       LocationResponse         `json:"location"`
	UsageList      []GetEquipmentUsageResponse `json:"usage_list"`
}

type GetEquipmentUsageResponse struct {
	ID          uint   `json:"id"`
	ResearcherID uint   `json:"researcher_id"`
	ResearcherName string `json:"researcher_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type EditEquipmentRequest struct {
	ID             uint   `json:"id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	PropertyNumber string `json:"property_number"`
	LocationID     uint   `json:"location_id" validate:"required"`
}

func GetEquipments(w http.ResponseWriter, r *http.Request) {
	laboratoryID := r.URL.Query().Get("laboratoryId")

	if laboratoryID == "" {
		http.Error(w, "Laboratory ID is required", http.StatusBadRequest)
		return
	}

	labID, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid laboratory ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var equipment []models.Equipment

	result := db.Where("laboratory_id = ?", labID).Find(&equipment)

	if result.Error != nil {
		http.Error(w, "Failed to fetch equipment", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddEquipment(w http.ResponseWriter, r *http.Request) {
	var req AddEquipmentRequest
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

	equipment := models.Equipment{
		Name:           req.Name,
		Description:    req.Description,
		PropertyNumber: req.PropertyNumber,
		LocationID:     req.LocationID,
		LaboratoryID:   labID,
	}

	db := database.GetDB()
	result := db.Create(&equipment)

	if result.Error != nil {
		http.Error(w, "Failed to create equipment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s was created successfully.", req.Name)
}

func GetDetailedEquipment(w http.ResponseWriter, r *http.Request) {
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
	var equipment models.Equipment

	result := db.Where("id = ? AND laboratory_id = ?", eqID, labID).
		Preload("Location").
		Preload("EquipmentUsages").
		Preload("EquipmentUsages.Researcher").
		First(&equipment)

	if result.Error != nil {
		http.Error(w, "Equipment not found", http.StatusNotFound)
		return
	}

	var usageList []GetEquipmentUsageResponse
	for _, usage := range equipment.EquipmentUsages {
		usageList = append(usageList, GetEquipmentUsageResponse{
			ID:             usage.ID,
			ResearcherID:   usage.ResearcherID,
			ResearcherName: usage.Researcher.Name,
			StartDate:      usage.StartDate.Format("2006-01-02T15:04:05Z07:00"),
			EndDate:        usage.EndDate.Format("2006-01-02T15:04:05Z07:00"),
			Description:    usage.Description,
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func EditEquipment(w http.ResponseWriter, r *http.Request) {
	var req EditEquipmentRequest
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
	var equipment models.Equipment

	result := db.Where("id = ? AND laboratory_id = ?", req.ID, labID).First(&equipment)
	if result.Error != nil {
		http.Error(w, "Equipment not found", http.StatusNotFound)
		return
	}

	equipment.Name = req.Name
	equipment.Description = req.Description
	equipment.PropertyNumber = req.PropertyNumber
	equipment.LocationID = req.LocationID

	result = db.Save(&equipment)
	if result.Error != nil {
		http.Error(w, "Failed to update equipment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s was updated successfully.", req.Name)
}

func RegisterEquipmentRoutes(r chi.Router) {
	r.Route("/equipment", func(r chi.Router) {
		r.Get("/GetEquipments", GetEquipments)
		r.Post("/AddEquipment", AddEquipment)
		r.Get("/GetDetailedEquipment", GetDetailedEquipment)
		r.Patch("/EditEquipment", EditEquipment)
	})
}