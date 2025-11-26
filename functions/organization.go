package functions

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type CreateOrganizationRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

type OrganizationResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LabCount    int       `json:"lab_count"`
	UserCount   int       `json:"user_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrganizationDetailResponse struct {
	ID           uint                `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Laboratories []LaboratorySummary `json:"laboratories"`
	CreatedAt    time.Time           `json:"created_at"`
}

type LaboratorySummary struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserCount   int    `json:"user_count"`
}

// CreateOrganization creates a new organization and assigns current user as org-coordinator
func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req CreateOrganizationRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Organization name is required")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := database.GetDB()

	// Check if user already belongs to an organization or laboratory
	var user models.User
	if err := db.First(&user, claims.UserID).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "User not found")
		return
	}

	if user.OrganizationID != nil || user.LaboratoryID != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "User already belongs to an organization or laboratory")
		return
	}

	// Create organization in transaction
	var createdOrg models.Organization
	err := db.Transaction(func(tx *gorm.DB) error {
		org := models.Organization{
			Name:        req.Name,
			Description: req.Description,
		}

		if err := tx.Create(&org).Error; err != nil {
			return err
		}

		createdOrg = org

		// Assign user as org-coordinator
		user.OrganizationID = &org.ID
		user.Role = "org-coordinator"
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	response := OrganizationResponse{
		ID:          createdOrg.ID,
		Name:        createdOrg.Name,
		Description: createdOrg.Description,
		LabCount:    0,
		UserCount:   1,
		CreatedAt:   createdOrg.CreatedAt,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// GetMyOrganization returns the organization of the current user
func GetMyOrganization(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !middleware.IsOrganizationCoordinator(claims) {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Organization coordinator access required")
		return
	}

	db := database.GetDB()

	var org models.Organization
	if err := db.Preload("Laboratories").First(&org, claims.OrganizationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "Organization not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Get lab summaries with user counts
	var labs []LaboratorySummary
	for _, lab := range org.Laboratories {
		var userCount int64
		db.Model(&models.User{}).Where("laboratory_id = ?", lab.ID).Count(&userCount)
		labs = append(labs, LaboratorySummary{
			ID:          lab.ID,
			Name:        lab.Name,
			Description: lab.Description,
			UserCount:   int(userCount),
		})
	}

	response := OrganizationDetailResponse{
		ID:           org.ID,
		Name:         org.Name,
		Description:  org.Description,
		Laboratories: labs,
		CreatedAt:    org.CreatedAt,
	}

	render.JSON(w, r, response)
}

// RegisterOrganizationRoutes registers organization-related routes
func RegisterOrganizationRoutes(r chi.Router) {
	r.Route("/organizations", func(r chi.Router) {
		r.Post("/", CreateOrganization)
		r.Get("/me", GetMyOrganization)
	})
}
