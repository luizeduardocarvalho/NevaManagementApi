package functions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type CreateLocationRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	SubLocationID *uint  `json:"sub_location_id"`
}

type UpdateLocationRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	SubLocationID *uint  `json:"sub_location_id"`
}

// GET /api/locations
func GetAllLocations(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	db := database.GetDB()
	var locations []models.Location

	if err := db.Where("laboratory_id = ?", claims.LaboratoryID).
		Preload("SubLocation").
		Order("name").
		Find(&locations).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch locations")
		return
	}

	render.JSON(w, r, locations)
}

// GET /api/locations/{id}
func GetLocationByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	locationID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid location ID")
		return
	}

	db := database.GetDB()
	var location models.Location

	if err := db.Where("id = ? AND laboratory_id = ?", id, claims.LaboratoryID).
		Preload("SubLocation").
		First(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "location not found")
		return
	}

	render.JSON(w, r, location)
}

// POST /api/locations
func CreateLocation(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()

	// Validate sub-location belongs to same laboratory if provided
	if req.SubLocationID != nil {
		var subLocation models.Location
		if err := db.Where("id = ? AND laboratory_id = ?", *req.SubLocationID, claims.LaboratoryID).
			First(&subLocation).Error; err != nil {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid sub-location")
			return
		}
	}
	location := models.Location{
		Name:          req.Name,
		Description:   req.Description,
		SubLocationID: req.SubLocationID,
		LaboratoryID:  claims.LaboratoryID,
	}

	if err := db.Create(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to create location")
		return
	}

	// Reload with relationships
	db.Preload("SubLocation").First(&location, location.ID)

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, location)
}

// PUT /api/locations/{id}
func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	locationID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid location ID")
		return
	}

	var req UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()
	var location models.Location

	if err := db.Where("id = ? AND laboratory_id = ?", id, claims.LaboratoryID).
		First(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "location not found")
		return
	}

	// Validate sub-location belongs to same laboratory if provided
	if req.SubLocationID != nil {
		var subLocation models.Location
		if err := db.Where("id = ? AND laboratory_id = ?", *req.SubLocationID, claims.LaboratoryID).
			First(&subLocation).Error; err != nil {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid sub-location")
			return
		}
	}

	// Update fields
	location.Name = req.Name
	location.Description = req.Description
	location.SubLocationID = req.SubLocationID

	if err := db.Save(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to update location")
		return
	}

	// Reload with relationships
	db.Preload("SubLocation").First(&location, location.ID)

	render.JSON(w, r, location)
}

// DELETE /api/locations/{id}
func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	locationID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid location ID")
		return
	}

	db := database.GetDB()
	var location models.Location

	if err := db.Where("id = ? AND laboratory_id = ?", id, claims.LaboratoryID).
		First(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "location not found")
		return
	}

	// Soft delete
	if err := db.Delete(&location).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to delete location")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterLocationRoutes registers all location routes
func RegisterLocationRoutes(r chi.Router) {
	r.Route("/locations", func(r chi.Router) {
		r.Get("/", GetAllLocations)
		r.Get("/{id}", GetLocationByID)
		r.Post("/", CreateLocation)
		r.Put("/{id}", UpdateLocation)
		r.Delete("/{id}", DeleteLocation)
	})
}
