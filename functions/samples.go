package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

// Request/Response Types

type CreateSampleRequest struct {
	Name                   string   `json:"name" validate:"required"`
	Description            string   `json:"description"`
	Origin                 string   `json:"origin"`
	IsolationDate          string   `json:"isolationDate"`
	Latitude               *float64 `json:"latitude"`
	Longitude              *float64 `json:"longitude"`
	SubcultureMedium       string   `json:"subcultureMedium"`
	SubcultureIntervalDays *int     `json:"subcultureIntervalDays"`
	ResearcherID           *uint    `json:"researcherId"`
	LocationID             *uint    `json:"locationId"`
	Tags                   []string `json:"tags"`
}

type UpdateSampleRequest struct {
	Name                   string   `json:"name" validate:"required"`
	Description            string   `json:"description"`
	Origin                 string   `json:"origin"`
	IsolationDate          string   `json:"isolationDate"`
	Latitude               *float64 `json:"latitude"`
	Longitude              *float64 `json:"longitude"`
	SubcultureMedium       string   `json:"subcultureMedium"`
	SubcultureIntervalDays *int     `json:"subcultureIntervalDays"`
	ResearcherID           *uint    `json:"researcherId"`
	LocationID             *uint    `json:"locationId"`
	Tags                   []string `json:"tags"`
}

type SampleListResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type SampleDetailResponse struct {
	ID                     uint      `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Origin                 string    `json:"origin"`
	IsolationDate          *string   `json:"isolationDate,omitempty"`
	Latitude               *float64  `json:"latitude,omitempty"`
	Longitude              *float64  `json:"longitude,omitempty"`
	SubcultureMedium       string    `json:"subcultureMedium"`
	SubcultureIntervalDays *int      `json:"subcultureIntervalDays,omitempty"`
	ResearcherID           *uint     `json:"researcherId,omitempty"`
	ResearcherName         string    `json:"researcherName"`
	LocationID             *uint     `json:"locationId,omitempty"`
	LocationName           string    `json:"locationName"`
	Tags                   []string  `json:"tags"`
	ReplicaCount           int64     `json:"replicaCount"`
	LaboratoryID           uint      `json:"laboratory_id"`
	CreatedAt              time.Time `json:"createdAt"`
}

// Helper Functions

func parseDateString(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// Try ISO 8601 format first
	t, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr)
	if err == nil {
		return &t, nil
	}

	// Try alternative ISO format
	t, err = time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
	if err == nil {
		return &t, nil
	}

	// Try date only format
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return &t, nil
	}

	return nil, fmt.Errorf("invalid date format")
}

func formatDateString(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	formatted := t.Format("2006-01-02T15:04:05.000Z")
	return &formatted
}

func tagsToJSON(tags []string) string {
	if tags == nil {
		tags = []string{}
	}
	b, _ := json.Marshal(tags)
	return string(b)
}

func tagsFromJSON(jsonStr string) []string {
	var tags []string
	if jsonStr == "" || jsonStr == "null" {
		return []string{}
	}
	json.Unmarshal([]byte(jsonStr), &tags)
	return tags
}

// Handlers

// GET /api/samples
func GetAllSamples(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	db := database.GetDB()
	query := db.Where("laboratory_id = ?", labID)

	// Apply search filter
	search := r.URL.Query().Get("search")
	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	// Apply tag filters (AND logic - must have ALL specified tags)
	tags := r.URL.Query()["tags"]
	if len(tags) > 0 {
		for _, tag := range tags {
			query = query.Where("tags::jsonb @> ?", fmt.Sprintf(`["%s"]`, tag))
		}
	}

	var samples []models.Sample
	if err := query.Order("created_at DESC").Find(&samples).Error; err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch samples"})
		return
	}

	// Return 404 if no samples found
	if len(samples) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "No samples found"})
		return
	}

	response := make([]SampleListResponse, 0, len(samples))
	for _, sample := range samples {
		response = append(response, SampleListResponse{
			ID:          sample.ID,
			Name:        sample.Name,
			Description: sample.Description,
			Tags:        tagsFromJSON(sample.Tags),
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// GET /api/samples/:id
func GetSampleByID(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	sampleIDStr := chi.URLParam(r, "id")
	sampleID, err := strconv.ParseUint(sampleIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid sample ID"})
		return
	}

	db := database.GetDB()
	var sample models.Sample

	if err := db.Where("id = ? AND laboratory_id = ?", sampleID, labID).
		Preload("Researcher").
		Preload("Location").
		First(&sample).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Sample not found"})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch sample"})
		return
	}

	// Count replicas
	var replicaCount int64
	db.Model(&models.Replica{}).Where("sample_id = ? AND laboratory_id = ?", sampleID, labID).Count(&replicaCount)

	// Build researcher name
	researcherName := ""
	if sample.Researcher != nil {
		researcherName = fmt.Sprintf("%s %s", sample.Researcher.FirstName, sample.Researcher.LastName)
	}

	// Build location name
	locationName := ""
	if sample.Location != nil {
		locationName = sample.Location.Name
	}

	response := SampleDetailResponse{
		ID:                     sample.ID,
		Name:                   sample.Name,
		Description:            sample.Description,
		Origin:                 sample.Origin,
		IsolationDate:          formatDateString(sample.IsolationDate),
		Latitude:               sample.Latitude,
		Longitude:              sample.Longitude,
		SubcultureMedium:       sample.SubcultureMedium,
		SubcultureIntervalDays: sample.SubcultureIntervalDays,
		ResearcherID:           sample.ResearcherID,
		ResearcherName:         researcherName,
		LocationID:             sample.LocationID,
		LocationName:           locationName,
		Tags:                   tagsFromJSON(sample.Tags),
		ReplicaCount:           replicaCount,
		LaboratoryID:           sample.LaboratoryID,
		CreatedAt:              sample.CreatedAt,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// POST /api/samples
func CreateSample(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	var req CreateSampleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	// Parse isolation date
	isolationDate, err := parseDateString(req.IsolationDate)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid isolation date format"})
		return
	}

	db := database.GetDB()

	// Use transaction to create sample and first replica atomically
	err = db.Transaction(func(tx *gorm.DB) error {
		sample := models.Sample{
			Name:                   req.Name,
			Description:            req.Description,
			Origin:                 req.Origin,
			IsolationDate:          isolationDate,
			Latitude:               req.Latitude,
			Longitude:              req.Longitude,
			SubcultureMedium:       req.SubcultureMedium,
			SubcultureIntervalDays: req.SubcultureIntervalDays,
			ResearcherID:           req.ResearcherID,
			LocationID:             req.LocationID,
			Tags:                   tagsToJSON(req.Tags),
			LaboratoryID:           labID,
		}

		if err := tx.Create(&sample).Error; err != nil {
			return err
		}

		// Auto-create first replica
		now := time.Now()
		var nextSubcultureDate *time.Time
		if req.SubcultureIntervalDays != nil && *req.SubcultureIntervalDays > 0 {
			next := now.AddDate(0, 0, *req.SubcultureIntervalDays)
			nextSubcultureDate = &next
		}

		replica := models.Replica{
			Name:               fmt.Sprintf("%s - Replica 1", req.Name),
			SampleID:           sample.ID,
			LocationID:         req.LocationID,
			Status:             "active",
			LastSubcultureDate: &now,
			NextSubcultureDate: nextSubcultureDate,
			LaboratoryID:       labID,
		}

		if err := tx.Create(&replica).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create sample"})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"message": "Sample created successfully with first replica"})
}

// PUT /api/samples/:id
func UpdateSample(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	sampleIDStr := chi.URLParam(r, "id")
	sampleID, err := strconv.ParseUint(sampleIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid sample ID"})
		return
	}

	var req UpdateSampleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	isolationDate, err := parseDateString(req.IsolationDate)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid isolation date format"})
		return
	}

	db := database.GetDB()
	var sample models.Sample

	if err := db.Where("id = ? AND laboratory_id = ?", sampleID, labID).First(&sample).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Sample not found"})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Database error"})
		return
	}

	// Update fields
	sample.Name = req.Name
	sample.Description = req.Description
	sample.Origin = req.Origin
	sample.IsolationDate = isolationDate
	sample.Latitude = req.Latitude
	sample.Longitude = req.Longitude
	sample.SubcultureMedium = req.SubcultureMedium
	sample.SubcultureIntervalDays = req.SubcultureIntervalDays
	sample.ResearcherID = req.ResearcherID
	sample.LocationID = req.LocationID
	sample.Tags = tagsToJSON(req.Tags)

	if err := db.Save(&sample).Error; err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update sample"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Sample updated successfully"})
}

// DELETE /api/samples/:id
func DeleteSample(w http.ResponseWriter, r *http.Request) {
	labID, err := getLaboratoryIDFromContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	sampleIDStr := chi.URLParam(r, "id")
	sampleID, err := strconv.ParseUint(sampleIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid sample ID"})
		return
	}

	db := database.GetDB()

	// Use transaction to delete sample and all its replicas
	err = db.Transaction(func(tx *gorm.DB) error {
		var sample models.Sample
		if err := tx.Where("id = ? AND laboratory_id = ?", sampleID, labID).First(&sample).Error; err != nil {
			return err
		}

		// Delete all replicas
		if err := tx.Where("sample_id = ? AND laboratory_id = ?", sampleID, labID).Delete(&models.Replica{}).Error; err != nil {
			return err
		}

		// Delete sample
		if err := tx.Delete(&sample).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Sample not found"})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete sample"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Sample deleted successfully"})
}

// RegisterSampleRoutes registers all sample routes
func RegisterSampleRoutes(r chi.Router) {
	r.Route("/samples", func(r chi.Router) {
		r.Get("/", GetAllSamples)
		r.Get("/{id}", GetSampleByID)
		r.Post("/", CreateSample)
		r.Put("/{id}", UpdateSample)
		r.Delete("/{id}", DeleteSample)
	})
}
