package functions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

// ResearcherSimple is a simplified researcher for dropdowns
type ResearcherSimple struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// ResearcherDetail is full researcher details
type ResearcherDetail struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	Status       string `json:"status"`
	LaboratoryID uint   `json:"laboratory_id"`
}

// CreateResearcherRequest for creating a new researcher (user in the lab)
type CreateResearcherRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=lab-coordinator technician student"`
}

// UpdateResearcherRequest for updating researcher details
type UpdateResearcherRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role" validate:"omitempty,oneof=lab-coordinator technician student"`
}

// GetAllResearchers returns researchers for the user's laboratory
// Supports ?simple=true for dropdown lists
func GetAllResearchers(w http.ResponseWriter, r *http.Request) {
	log.Println("[RESEARCHERS] GetAllResearchers endpoint called")

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		log.Println("[RESEARCHERS] No user claims found")
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	log.Printf("[RESEARCHERS] User: %s, Email: %s, LaboratoryID: %d, Role: %s",
		claims.ClerkUserID, claims.Email, claims.LaboratoryID, claims.Role)

	if claims.LaboratoryID == 0 {
		log.Println("[RESEARCHERS] User does not belong to a laboratory")
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	db := database.GetDB()

	// Check if simple format is requested
	simple := r.URL.Query().Get("simple") == "true"
	log.Printf("[RESEARCHERS] Fetching researchers for laboratory_id=%d, simple=%v", claims.LaboratoryID, simple)

	var users []models.User
	if err := db.Where("laboratory_id = ?", claims.LaboratoryID).
		Order("last_name, first_name").
		Find(&users).Error; err != nil {
		log.Printf("[RESEARCHERS] Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch researchers")
		return
	}

	log.Printf("[RESEARCHERS] Found %d researchers", len(users))

	if simple {
		// Return simplified format for dropdowns
		var researchers []ResearcherSimple
		for _, user := range users {
			researchers = append(researchers, ResearcherSimple{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Role:      user.Role,
			})
		}
		render.JSON(w, r, researchers)
	} else {
		// Return full details
		var researchers []ResearcherDetail
		for _, user := range users {
			labID := uint(0)
			if user.LaboratoryID != nil {
				labID = *user.LaboratoryID
			}
			researchers = append(researchers, ResearcherDetail{
				ID:           user.ID,
				Email:        user.Email,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				Role:         user.Role,
				Status:       user.Status,
				LaboratoryID: labID,
			})
		}
		render.JSON(w, r, researchers)
	}
}

// GetResearcherByID returns a specific researcher by ID
func GetResearcherByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid researcher ID")
		return
	}

	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ? AND laboratory_id = ?", id, claims.LaboratoryID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "researcher not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	labID := uint(0)
	if user.LaboratoryID != nil {
		labID = *user.LaboratoryID
	}

	researcher := ResearcherDetail{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		Status:       user.Status,
		LaboratoryID: labID,
	}

	render.JSON(w, r, researcher)
}

// CreateResearcher creates a new researcher (user) for the laboratory
// This sends an invitation to join the lab
func CreateResearcher(w http.ResponseWriter, r *http.Request) {
	var req CreateResearcherRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	// Only lab-coordinator can create researchers
	if claims.Role != "lab-coordinator" {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "lab coordinator access required")
		return
	}

	// Validate role
	validRoles := map[string]bool{"lab-coordinator": true, "technician": true, "student": true}
	if !validRoles[req.Role] {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid role. Must be one of: lab-coordinator, technician, student")
		return
	}

	db := database.GetDB()

	// Check if user with this email already exists in the lab
	var existingUser models.User
	err := db.Where("email = ? AND laboratory_id = ?", req.Email, claims.LaboratoryID).First(&existingUser).Error
	if err == nil {
		middleware.ErrorResponse(w, r, http.StatusConflict, "user with this email already exists in this laboratory")
		return
	}
	if err != gorm.ErrRecordNotFound {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	// Note: This endpoint creates a skeleton user record
	// In a real implementation, you'd want to send an invitation instead
	// For now, creating a user directly for simplicity
	user := models.User{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		ClerkUserID:  "", // Will be populated when they sign up
		LaboratoryID: &claims.LaboratoryID,
		Role:         req.Role,
		Status:       "pending",
	}

	if err := db.Create(&user).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to create researcher")
		return
	}

	labID := uint(0)
	if user.LaboratoryID != nil {
		labID = *user.LaboratoryID
	}

	researcher := ResearcherDetail{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		Status:       user.Status,
		LaboratoryID: labID,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, researcher)
}

// UpdateResearcher updates researcher details
func UpdateResearcher(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	// Only lab-coordinator can update researchers
	if claims.Role != "lab-coordinator" {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "lab coordinator access required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid researcher ID")
		return
	}

	var req UpdateResearcherRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate role if provided
	if req.Role != "" {
		validRoles := map[string]bool{"lab-coordinator": true, "technician": true, "student": true}
		if !validRoles[req.Role] {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid role. Must be one of: lab-coordinator, technician, student")
			return
		}
	}

	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ? AND laboratory_id = ?", id, claims.LaboratoryID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "researcher not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := db.Save(&user).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to update researcher")
		return
	}

	labID := uint(0)
	if user.LaboratoryID != nil {
		labID = *user.LaboratoryID
	}

	researcher := ResearcherDetail{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		Status:       user.Status,
		LaboratoryID: labID,
	}

	render.JSON(w, r, researcher)
}

// DeleteResearcher soft-deletes a researcher (user)
func DeleteResearcher(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	// Only lab-coordinator can delete researchers
	if claims.Role != "lab-coordinator" {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "lab coordinator access required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid researcher ID")
		return
	}

	// Prevent deleting yourself
	if uint(id) == claims.UserID {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "cannot delete yourself")
		return
	}

	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ? AND laboratory_id = ?", id, claims.LaboratoryID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "researcher not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	// Soft delete
	if err := db.Delete(&user).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to delete researcher")
		return
	}

	render.Status(r, http.StatusNoContent)
	w.Write(nil)
}

// RegisterResearcherRoutes registers researcher-related routes
func RegisterResearcherRoutes(r chi.Router) {
	log.Println("[RESEARCHERS] Registering researcher routes at /researchers")
	r.Route("/researchers", func(r chi.Router) {
		r.Get("/", GetAllResearchers)
		r.Post("/", CreateResearcher)
		r.Get("/{id}", GetResearcherByID)
		r.Put("/{id}", UpdateResearcher)
		r.Delete("/{id}", DeleteResearcher)
	})
	log.Println("[RESEARCHERS] Researcher routes registered successfully")
}
