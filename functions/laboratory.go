package functions

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/crypto"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/email"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

type CreateLaboratoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
	Address     string `json:"address,max=200"`
}

type CreateInvitationRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"role" validate:"required,oneof=coordinator technician student"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type InvitationResponse struct {
	ID              uint       `json:"id"`
	LaboratoryID    uint       `json:"laboratory_id"`
	LaboratoryName  string     `json:"laboratory_name,omitempty"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	FirstName       string     `json:"first_name,omitempty"`
	LastName        string     `json:"last_name,omitempty"`
	InvitedBy       uint       `json:"invited_by"`
	InvitedByName   string     `json:"invited_by_name,omitempty"`
	InvitationToken string     `json:"invitation_token"`
	Status          string     `json:"status"`
	ExpiresAt       time.Time  `json:"expires_at"`
	AcceptedAt      *time.Time `json:"accepted_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type InvitationPublicResponse struct {
	ID             uint   `json:"id"`
	LaboratoryID   uint   `json:"laboratory_id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	ExpiresAt      time.Time `json:"expires_at"`
	LaboratoryName string `json:"laboratory_name"`
	InvitedByName  string `json:"invited_by_name"`
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
// @Description Create an invitation for a user to join the laboratory (coordinator only)
// @Tags invitations
// @Accept json
// @Produce json
// @Param laboratoryId path int true "Laboratory ID"
// @Param invitation body CreateInvitationRequest true "Invitation data"
// @Success 201 {object} InvitationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
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

	// Validate email format
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Valid email is required")
		return
	}

	// Validate role
	validRoles := map[string]bool{"coordinator": true, "technician": true, "student": true}
	if !validRoles[req.Role] {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid role. Must be one of: coordinator, technician, student")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify user is coordinator of this laboratory
	if claims.LaboratoryID != uint(id) || (claims.Role != "coordinator" && claims.Role != "Admin") {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Coordinator access required")
		return
	}

	db := database.GetDB()

	// Normalize email to lowercase
	normalizedEmail := strings.ToLower(req.Email)

	// Check if user already exists with this email
	var existingUser models.User
	if err := db.Where("email = ?", normalizedEmail).First(&existingUser).Error; err == nil {
		middleware.ErrorResponse(w, r, http.StatusConflict, "User with this email already exists")
		return
	}

	// Check if pending invitation already exists for this email
	var existingInvitation models.LaboratoryInvitation
	if err := db.Where("email = ? AND laboratory_id = ? AND status = ?", normalizedEmail, id, "pending").First(&existingInvitation).Error; err == nil {
		middleware.ErrorResponse(w, r, http.StatusConflict, "Pending invitation already exists for this email")
		return
	}

	// Get current user info
	var currentUser models.User
	if err := db.First(&currentUser, claims.UserID).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to get user info")
		return
	}

	// Get laboratory info for email
	var laboratory models.Laboratory
	if err := db.First(&laboratory, id).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusNotFound, "Laboratory not found")
		return
	}

	// Generate secure token
	token, err := crypto.GenerateInvitationToken()
	if err != nil {
		log.Printf("Failed to generate invitation token: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to generate invitation token")
		return
	}

	// Create invitation
	invitation := models.LaboratoryInvitation{
		LaboratoryID:    uint(id),
		Email:           normalizedEmail,
		Role:            req.Role,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		InvitedBy:       claims.UserID,
		InvitationToken: token,
		Status:          "pending",
		ExpiresAt:       time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := db.Create(&invitation).Error; err != nil {
		log.Printf("Failed to create invitation: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to create invitation")
		return
	}

	// Send invitation email
	emailService := email.GetEmailService()
	invitedByName := currentUser.FirstName + " " + currentUser.LastName
	if invitedByName == " " {
		invitedByName = currentUser.Email
	}

	if err := emailService.SendInvitationEmail(
		normalizedEmail,
		req.FirstName,
		laboratory.Name,
		req.Role,
		invitedByName,
		token,
		invitation.ExpiresAt,
	); err != nil {
		log.Printf("Failed to send invitation email: %v", err)
		// Don't fail the request if email fails - invitation is still created
	}

	response := InvitationResponse{
		ID:              invitation.ID,
		LaboratoryID:    invitation.LaboratoryID,
		Email:           invitation.Email,
		Role:            invitation.Role,
		FirstName:       invitation.FirstName,
		LastName:        invitation.LastName,
		InvitedBy:       invitation.InvitedBy,
		InvitationToken: invitation.InvitationToken,
		Status:          invitation.Status,
		ExpiresAt:       invitation.ExpiresAt,
		CreatedAt:       invitation.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response)
}

// GetInvitations retrieves all invitations for a laboratory
// @Summary Get Laboratory Invitations
// @Description Get all invitations for the specified laboratory (coordinator only)
// @Tags invitations
// @Produce json
// @Param laboratoryId path int true "Laboratory ID"
// @Param status query string false "Filter by status (pending, accepted, expired, cancelled)"
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

	// Verify user is coordinator of this laboratory
	if claims.LaboratoryID != uint(id) || (claims.Role != "coordinator" && claims.Role != "Admin") {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Coordinator access required")
		return
	}

	db := database.GetDB()

	// Build query
	query := db.Where("laboratory_id = ?", id)

	// Filter by status if provided
	statusFilter := r.URL.Query().Get("status")
	if statusFilter != "" {
		validStatuses := map[string]bool{"pending": true, "accepted": true, "expired": true, "cancelled": true}
		if !validStatuses[statusFilter] {
			middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid status filter")
			return
		}
		query = query.Where("status = ?", statusFilter)
	}

	var invitations []models.LaboratoryInvitation
	if err := query.Preload("Laboratory").Preload("InvitedByUser").
		Order("created_at DESC").Find(&invitations).Error; err != nil {
		log.Printf("Failed to fetch invitations: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to fetch invitations")
		return
	}

	var responses []InvitationResponse
	for _, inv := range invitations {
		response := InvitationResponse{
			ID:              inv.ID,
			LaboratoryID:    inv.LaboratoryID,
			Email:           inv.Email,
			Role:            inv.Role,
			FirstName:       inv.FirstName,
			LastName:        inv.LastName,
			InvitedBy:       inv.InvitedBy,
			InvitationToken: inv.InvitationToken,
			Status:          inv.Status,
			ExpiresAt:       inv.ExpiresAt,
			AcceptedAt:      inv.AcceptedAt,
			CreatedAt:       inv.CreatedAt,
		}

		if inv.Laboratory.Name != "" {
			response.LaboratoryName = inv.Laboratory.Name
		}

		if inv.InvitedByUser != nil {
			invitedByName := inv.InvitedByUser.FirstName + " " + inv.InvitedByUser.LastName
			if invitedByName == " " {
				invitedByName = inv.InvitedByUser.Email
			}
			response.InvitedByName = invitedByName
		}

		responses = append(responses, response)
	}

	render.JSON(w, r, responses)
}

// GetInvitationByToken retrieves invitation details by token (public endpoint)
// @Summary Get Invitation by Token
// @Description Get invitation details using the invitation token (public, no auth required)
// @Tags invitations
// @Produce json
// @Param token path string true "Invitation Token"
// @Success 200 {object} InvitationPublicResponse
// @Failure 404 {object} map[string]string
// @Failure 410 {object} map[string]string
// @Router /invitations/token/{token} [get]
func GetInvitationByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Token is required")
		return
	}

	db := database.GetDB()

	var invitation models.LaboratoryInvitation
	if err := db.Preload("Laboratory").Preload("InvitedByUser").
		Where("invitation_token = ?", token).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "Invitation not found")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Check if invitation is expired or cancelled
	if invitation.Status == "expired" || invitation.Status == "cancelled" {
		middleware.ErrorResponse(w, r, http.StatusGone, "Invitation is no longer valid")
		return
	}

	// Check if invitation has expired
	if time.Now().After(invitation.ExpiresAt) && invitation.Status == "pending" {
		// Mark as expired
		invitation.Status = "expired"
		db.Save(&invitation)
		middleware.ErrorResponse(w, r, http.StatusGone, "Invitation has expired")
		return
	}

	invitedByName := ""
	if invitation.InvitedByUser != nil {
		invitedByName = invitation.InvitedByUser.FirstName + " " + invitation.InvitedByUser.LastName
		if invitedByName == " " {
			invitedByName = invitation.InvitedByUser.Email
		}
	}

	response := InvitationPublicResponse{
		ID:             invitation.ID,
		LaboratoryID:   invitation.LaboratoryID,
		Email:          invitation.Email,
		Role:           invitation.Role,
		Status:         invitation.Status,
		ExpiresAt:      invitation.ExpiresAt,
		LaboratoryName: invitation.Laboratory.Name,
		InvitedByName:  invitedByName,
	}

	render.JSON(w, r, response)
}

// ResendInvitation resends an invitation email
// @Summary Resend Invitation
// @Description Resend invitation email and extend expiration (coordinator only)
// @Tags invitations
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /invitations/{id}/resend [post]
func ResendInvitation(w http.ResponseWriter, r *http.Request) {
	invitationID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(invitationID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := database.GetDB()

	var invitation models.LaboratoryInvitation
	if err := db.Preload("Laboratory").Preload("InvitedByUser").
		First(&invitation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "Invitation not found")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Verify user is coordinator of the laboratory
	if claims.LaboratoryID != invitation.LaboratoryID || (claims.Role != "coordinator" && claims.Role != "Admin") {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Coordinator access required")
		return
	}

	// Check if invitation can be resent
	if invitation.Status != "pending" && invitation.Status != "expired" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Can only resend pending or expired invitations")
		return
	}

	// Generate new token and extend expiration
	newToken, err := crypto.GenerateInvitationToken()
	if err != nil {
		log.Printf("Failed to generate invitation token: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to generate invitation token")
		return
	}

	invitation.InvitationToken = newToken
	invitation.Status = "pending"
	invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	if err := db.Save(&invitation).Error; err != nil {
		log.Printf("Failed to update invitation: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to update invitation")
		return
	}

	// Send invitation email
	emailService := email.GetEmailService()
	invitedByName := ""
	if invitation.InvitedByUser != nil {
		invitedByName = invitation.InvitedByUser.FirstName + " " + invitation.InvitedByUser.LastName
		if invitedByName == " " {
			invitedByName = invitation.InvitedByUser.Email
		}
	}

	if err := emailService.SendInvitationEmail(
		invitation.Email,
		invitation.FirstName,
		invitation.Laboratory.Name,
		invitation.Role,
		invitedByName,
		newToken,
		invitation.ExpiresAt,
	); err != nil {
		log.Printf("Failed to send invitation email: %v", err)
		// Don't fail the request if email fails
	}

	render.JSON(w, r, map[string]interface{}{
		"message":    "Invitation resent successfully",
		"expires_at": invitation.ExpiresAt,
	})
}

// CancelInvitation cancels a pending invitation
// @Summary Cancel Invitation
// @Description Cancel a pending invitation (coordinator only)
// @Tags invitations
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /invitations/{id}/cancel [post]
func CancelInvitation(w http.ResponseWriter, r *http.Request) {
	invitationID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(invitationID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := database.GetDB()

	var invitation models.LaboratoryInvitation
	if err := db.First(&invitation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "Invitation not found")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Verify user is coordinator of the laboratory
	if claims.LaboratoryID != invitation.LaboratoryID || (claims.Role != "coordinator" && claims.Role != "Admin") {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Coordinator access required")
		return
	}

	// Check if invitation can be cancelled
	if invitation.Status == "accepted" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Cannot cancel an accepted invitation")
		return
	}

	if invitation.Status == "cancelled" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invitation is already cancelled")
		return
	}

	// Cancel invitation
	invitation.Status = "cancelled"
	if err := db.Save(&invitation).Error; err != nil {
		log.Printf("Failed to cancel invitation: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to cancel invitation")
		return
	}

	render.JSON(w, r, map[string]string{
		"message": "Invitation cancelled successfully",
	})
}

// DeleteInvitation permanently deletes an invitation
// @Summary Delete Invitation
// @Description Permanently delete an invitation (coordinator only)
// @Tags invitations
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /invitations/{id} [delete]
func DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	invitationID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(invitationID, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := database.GetDB()

	var invitation models.LaboratoryInvitation
	if err := db.First(&invitation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "Invitation not found")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Verify user is coordinator of the laboratory
	if claims.LaboratoryID != invitation.LaboratoryID || (claims.Role != "coordinator" && claims.Role != "Admin") {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "Coordinator access required")
		return
	}

	// Delete invitation (soft delete)
	if err := db.Delete(&invitation).Error; err != nil {
		log.Printf("Failed to delete invitation: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to delete invitation")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}