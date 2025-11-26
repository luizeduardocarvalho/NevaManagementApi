package functions

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

// Response Types

type ReplicaResponse struct {
	ID                 uint               `json:"id"`
	Name               string             `json:"name"`
	Sample             *ReplicaSampleInfo `json:"sample,omitempty"`
	SampleID           uint               `json:"sample_id"`
	Location           *ReplicaLocation   `json:"location,omitempty"`
	LocationID         *uint              `json:"location_id,omitempty"`
	Status             string             `json:"status"`
	LastSubcultureDate *string            `json:"lastSubcultureDate,omitempty"`
	NextSubcultureDate *string            `json:"nextSubcultureDate,omitempty"`
	LaboratoryID       uint               `json:"laboratory_id"`
}

type ReplicaSampleInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ReplicaLocation struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Handlers

// GET /api/replicas/sample/:sampleId
func GetReplicasBySampleID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	sampleIDStr := chi.URLParam(r, "sampleId")
	sampleID, err := strconv.ParseUint(sampleIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid sample ID")
		return
	}

	db := database.GetDB()

	// Verify sample exists and belongs to user's laboratory
	var sample models.Sample
	if err := db.Where("id = ? AND laboratory_id = ?", sampleID, claims.LaboratoryID).First(&sample).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "sample not found")
		return
	}

	// Get all replicas for the sample
	var replicas []models.Replica
	if err := db.Where("sample_id = ? AND laboratory_id = ?", sampleID, claims.LaboratoryID).
		Preload("Sample").
		Preload("Location").
		Order("created_at ASC").
		Find(&replicas).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch replicas")
		return
	}

	// Build response
	response := make([]ReplicaResponse, 0, len(replicas))
	for _, replica := range replicas {
		var sampleInfo *ReplicaSampleInfo
		if replica.Sample != nil {
			sampleInfo = &ReplicaSampleInfo{
				ID:   replica.Sample.ID,
				Name: replica.Sample.Name,
			}
		}

		var location *ReplicaLocation
		if replica.Location != nil {
			location = &ReplicaLocation{
				ID:          replica.Location.ID,
				Name:        replica.Location.Name,
				Description: replica.Location.Description,
			}
		}

		response = append(response, ReplicaResponse{
			ID:                 replica.ID,
			Name:               replica.Name,
			Sample:             sampleInfo,
			SampleID:           replica.SampleID,
			Location:           location,
			LocationID:         replica.LocationID,
			Status:             replica.Status,
			LastSubcultureDate: formatDateString(replica.LastSubcultureDate),
			NextSubcultureDate: formatDateString(replica.NextSubcultureDate),
			LaboratoryID:       replica.LaboratoryID,
		})
	}

	render.JSON(w, r, response)
}

// RegisterReplicaRoutes registers all replica routes
func RegisterReplicaRoutes(r chi.Router) {
	r.Route("/replicas", func(r chi.Router) {
		r.Get("/sample/{sampleId}", GetReplicasBySampleID)
	})
}
