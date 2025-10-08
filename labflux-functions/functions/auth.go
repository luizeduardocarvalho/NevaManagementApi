package functions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"labflux-functions/pkg/database"
	"labflux-functions/pkg/middleware"
	"labflux-functions/pkg/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

// Clerk webhook event types
type ClerkWebhookEvent struct {
	Type   string          `json:"type"`
	Object string          `json:"object"`
	Data   json.RawMessage `json:"data"`
}

type ClerkUser struct {
	ID             string              `json:"id"`
	EmailAddresses []ClerkEmailAddress `json:"email_addresses"`
	FirstName      *string             `json:"first_name"`
	LastName       *string             `json:"last_name"`
	CreatedAt      int64               `json:"created_at"`
	UpdatedAt      int64               `json:"updated_at"`
}

type ClerkEmailAddress struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
	Verification struct {
		Status string `json:"status"`
	} `json:"verification"`
}

type AcceptInvitationRequest struct {
	ClerkUserID string `json:"clerk_user_id" validate:"required"`
}

type UserResponse struct {
	ID           uint      `json:"id"`
	ClerkUserID  string    `json:"clerk_user_id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	LaboratoryID *uint     `json:"laboratory_id"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// ClerkWebhookHandler handles Clerk webhook events for user management
func ClerkWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Verify Clerk webhook signature for security
	// webhookSecret := os.Getenv("CLERK_WEBHOOK_SECRET")
	// if !verifyClerkSignature(r, webhookSecret) {
	//     middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Invalid webhook signature")
	//     return
	// }

	var event ClerkWebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	switch event.Type {
	case "user.created":
		handleUserCreated(w, r, event.Data)
	case "user.updated":
		handleUserUpdated(w, r, event.Data)
	case "user.deleted":
		handleUserDeleted(w, r, event.Data)
	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
		w.WriteHeader(http.StatusOK)
	}
}

func handleUserCreated(w http.ResponseWriter, r *http.Request, data json.RawMessage) {
	var clerkUser ClerkUser
	if err := json.Unmarshal(data, &clerkUser); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid user data")
		return
	}

	// Get primary email
	primaryEmail := ""
	for _, email := range clerkUser.EmailAddresses {
		if email.Verification.Status == "verified" {
			primaryEmail = email.EmailAddress
			break
		}
	}

	if primaryEmail == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "No verified email found")
		return
	}

	db := database.GetDB()

	// Create user in our database
	user := models.User{
		ClerkUserID: clerkUser.ID,
		Email:       primaryEmail,
		FirstName:   getStringValue(clerkUser.FirstName),
		LastName:    getStringValue(clerkUser.LastName),
		Role:        "user", // Default role
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error creating user")
		return
	}

	log.Printf("Created user %s (%s)", user.Email, user.ClerkUserID)
	w.WriteHeader(http.StatusOK)
}

func handleUserUpdated(w http.ResponseWriter, r *http.Request, data json.RawMessage) {
	var clerkUser ClerkUser
	if err := json.Unmarshal(data, &clerkUser); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid user data")
		return
	}

	db := database.GetDB()

	// Find and update user
	var user models.User
	if err := db.Where("clerk_user_id = ?", clerkUser.ID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User not found for update: %s", clerkUser.ID)
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Update user fields
	user.FirstName = getStringValue(clerkUser.FirstName)
	user.LastName = getStringValue(clerkUser.LastName)

	// Update primary email if changed
	for _, email := range clerkUser.EmailAddresses {
		if email.Verification.Status == "verified" {
			user.Email = email.EmailAddress
			break
		}
	}

	if err := db.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error updating user")
		return
	}

	log.Printf("Updated user %s (%s)", user.Email, user.ClerkUserID)
	w.WriteHeader(http.StatusOK)
}

func handleUserDeleted(w http.ResponseWriter, r *http.Request, data json.RawMessage) {
	var clerkUser ClerkUser
	if err := json.Unmarshal(data, &clerkUser); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid user data")
		return
	}

	db := database.GetDB()

	// Soft delete user
	if err := db.Where("clerk_user_id = ?", clerkUser.ID).Delete(&models.User{}).Error; err != nil {
		log.Printf("Error deleting user: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error deleting user")
		return
	}

	log.Printf("Deleted user %s", clerkUser.ID)
	w.WriteHeader(http.StatusOK)
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// AcceptInvitation handles accepting laboratory invitations (for Clerk users)
func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	invitationToken := chi.URLParam(r, "token")
	if invitationToken == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invitation token is required")
		return
	}

	var req AcceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.ClerkUserID == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Clerk user ID is required")
		return
	}

	db := database.GetDB()

	// Find invitation
	var invitation models.LaboratoryInvitation
	if err := db.Preload("Laboratory").Where("invitation_token = ? AND is_accepted = ? AND expires_at > ?",
		invitationToken, false, time.Now()).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "Invalid or expired invitation")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Find user by Clerk ID
	var user models.User
	if err := db.Where("clerk_user_id = ?", req.ClerkUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "User not found")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Verify email matches invitation
	if user.Email != invitation.Email {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Email mismatch")
		return
	}

	// Update user with laboratory assignment and accept invitation in transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Update user with laboratory assignment
		user.LaboratoryID = &invitation.LaboratoryID
		user.Role = invitation.Role
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Mark invitation as accepted
		now := time.Now()
		invitation.IsAccepted = true
		invitation.AcceptedAt = &now
		if err := tx.Save(&invitation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Transaction error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error accepting invitation")
		return
	}

	// Fetch updated user with laboratory info
	if err := db.Preload("Laboratory").First(&user, user.ID).Error; err != nil {
		log.Printf("Error fetching updated user: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error fetching user")
		return
	}

	userResponse := UserResponse{
		ID:           user.ID,
		ClerkUserID:  user.ClerkUserID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LaboratoryID: user.LaboratoryID,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, userResponse)
}

// GetCurrentUser returns current user information
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "No user claims found")
		return
	}

	db := database.GetDB()

	var user models.User
	if err := db.Preload("Laboratory").Where("clerk_user_id = ?", claims.ClerkUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "User not found")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	userResponse := UserResponse{
		ID:           user.ID,
		ClerkUserID:  user.ClerkUserID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LaboratoryID: user.LaboratoryID,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}

	render.JSON(w, r, userResponse)
}

