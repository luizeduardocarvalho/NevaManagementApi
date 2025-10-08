package functions

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"labflux-functions/pkg/database"
	"labflux-functions/pkg/middleware"
	"labflux-functions/pkg/models"
)

type CreateLaboratoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
	Address     string `json:"address,max=200"`
}

type CreateInvitationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=Admin Member"`
}

type InvitationResponse struct {
	ID              uint      `json:"id"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	InvitationToken string    `json:"invitation_token"`
	ExpiresAt       time.Time `json:"expires_at"`
	IsAccepted      bool      `json:"is_accepted"`
}

// CreateLaboratory creates a new laboratory for the current user
// @Summary Create Laboratory
// @Description Creates a new laboratory and assigns the current user as admin
// @Tags laboratories
// @Accept json
// @Produce json
// @Param laboratory body CreateLaboratoryRequest true "Laboratory data"
// @Success 201 {object} models.Laboratory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /laboratories [post]
func CreateLaboratory(w http.ResponseWriter, r *http.Request) {
	var req CreateLaboratoryRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := database.GetDB()

	// Check if user already has a laboratory
	var user models.User
	if err := db.First(&user, claims.UserID).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "User not found")
		return
	}

	if user.LaboratoryID != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "User already belongs to a laboratory")
		return
	}

	// Create laboratory
	laboratory := models.Laboratory{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
	}

	if err := db.Create(&laboratory).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to create laboratory")
		return
	}

	// Assign user to laboratory as admin
	user.LaboratoryID = &laboratory.ID
	user.Role = "Admin"
	if err := db.Save(&user).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to assign user to laboratory")
		return
	}

	render.JSON(w, r, laboratory)
}

// GetLaboratory retrieves laboratory information
// @Summary Get Laboratory
// @Description Get laboratory details by ID (user must belong to this laboratory)
// @Tags laboratories
// @Produce json
// @Param laboratoryId path int true "Laboratory ID"
// @Success 200 {object} models.Laboratory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /laboratories/{laboratoryId} [get]
func GetLaboratory(w http.ResponseWriter, r *http.Request) {
	laboratoryID := chi.URLParam(r, "laboratoryId")
	id, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid laboratory ID")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify user belongs to this laboratory
	if claims.LaboratoryID != uint(id) {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Access denied")
		return
	}

	db := database.GetDB()
	var laboratory models.Laboratory
	if err := db.First(&laboratory, id).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "Laboratory not found")
		return
	}

	render.JSON(w, r, laboratory)
}

// CreateInvitation creates a new invitation to join a laboratory
// @Summary Create Laboratory Invitation
// @Description Create an invitation for a user to join the laboratory (admin only)
// @Tags invitations
// @Accept json
// @Produce json
// @Param laboratoryId path int true "Laboratory ID"
// @Param invitation body CreateInvitationRequest true "Invitation data"
// @Success 201 {object} InvitationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /laboratories/{laboratoryId}/invitations [post]
func CreateInvitation(w http.ResponseWriter, r *http.Request) {
	laboratoryID := chi.URLParam(r, "laboratoryId")
	id, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid laboratory ID")
		return
	}

	var req CreateInvitationRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify user is admin of this laboratory
	if claims.LaboratoryID != uint(id) || claims.Role != "Admin" {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Admin access required")
		return
	}

	db := database.GetDB()

	// Check if invitation already exists for this email
	var existingInvitation models.LaboratoryInvitation
	if err := db.Where("email = ? AND laboratory_id = ? AND is_accepted = false", req.Email, id).First(&existingInvitation).Error; err == nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invitation already exists for this email")
		return
	}

	// Create invitation
	invitation := models.LaboratoryInvitation{
		LaboratoryID:    uint(id),
		Email:           req.Email,
		Role:            req.Role,
		InvitationToken: uuid.New().String(),
		ExpiresAt:       time.Now().Add(7 * 24 * time.Hour), // 7 days
		IsAccepted:      false,
	}

	if err := db.Create(&invitation).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to create invitation")
		return
	}

	response := InvitationResponse{
		ID:              invitation.ID,
		Email:           invitation.Email,
		Role:            invitation.Role,
		InvitationToken: invitation.InvitationToken,
		ExpiresAt:       invitation.ExpiresAt,
		IsAccepted:      invitation.IsAccepted,
	}

	render.JSON(w, r, response)
}

// GetInvitations retrieves all invitations for a laboratory
// @Summary Get Laboratory Invitations
// @Description Get all invitations for the specified laboratory
// @Tags invitations
// @Produce json
// @Param laboratoryId path int true "Laboratory ID"
// @Success 200 {array} InvitationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /laboratories/{laboratoryId}/invitations [get]
func GetInvitations(w http.ResponseWriter, r *http.Request) {
	laboratoryID := chi.URLParam(r, "laboratoryId")
	id, err := strconv.ParseUint(laboratoryID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid laboratory ID")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify user belongs to this laboratory
	if claims.LaboratoryID != uint(id) {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Access denied")
		return
	}

	db := database.GetDB()
	var invitations []models.LaboratoryInvitation
	if err := db.Where("laboratory_id = ?", id).Find(&invitations).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to fetch invitations")
		return
	}

	var responses []InvitationResponse
	for _, inv := range invitations {
		responses = append(responses, InvitationResponse{
			ID:              inv.ID,
			Email:           inv.Email,
			Role:            inv.Role,
			InvitationToken: inv.InvitationToken,
			ExpiresAt:       inv.ExpiresAt,
			IsAccepted:      inv.IsAccepted,
		})
	}

	render.JSON(w, r, responses)
}