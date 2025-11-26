package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	svix "github.com/svix/svix-webhooks/go"
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
	ID             uint      `json:"id"`
	ClerkUserID    string    `json:"clerk_user_id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	OrganizationID *uint     `json:"organization_id"`
	LaboratoryID   *uint     `json:"laboratory_id"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}

// Backend authentication request types
type SignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Clerk API response types
type ClerkSignUpResponse struct {
	ID                   string                 `json:"id"`
	Object               string                 `json:"object"`
	Status               string                 `json:"status"`
	EmailAddress         string                 `json:"email_address"`
	FirstName            string                 `json:"first_name"`
	LastName             string                 `json:"last_name"`
	CreatedSessionID     string                 `json:"created_session_id"`
	CreatedUserID        string                 `json:"created_user_id"`
	EmailAddressID       string                 `json:"email_address_id"`
	Verifications        map[string]interface{} `json:"verifications"`
}

type ClerkSessionResponse struct {
	ID               string    `json:"id"`
	ClientID         string    `json:"client_id"`
	UserID           string    `json:"user_id"`
	Status           string    `json:"status"`
	LastActiveAt     int64     `json:"last_active_at"`
	ExpireAt         int64     `json:"expire_at"`
	AbandonAt        int64     `json:"abandon_at"`
	CreatedAt        int64     `json:"created_at"`
	UpdatedAt        int64     `json:"updated_at"`
}

type ClerkSessionTokenResponse struct {
	JWT string `json:"jwt"`
}

// ClerkWebhookHandler handles Clerk webhook events for user management
func ClerkWebhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookSecret := os.Getenv("CLERK_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("CLERK_WEBHOOK_SECRET not set - webhook signature verification disabled")
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Server configuration error")
		return
	}

	// Read the request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Verify webhook signature
	headers := http.Header{}
	headers.Set("svix-id", r.Header.Get("svix-id"))
	headers.Set("svix-timestamp", r.Header.Get("svix-timestamp"))
	headers.Set("svix-signature", r.Header.Get("svix-signature"))

	wh, err := svix.NewWebhook(webhookSecret)
	if err != nil {
		log.Printf("Failed to create webhook verifier: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Webhook verification setup failed")
		return
	}

	var event ClerkWebhookEvent
	if err := wh.Verify(payload, headers); err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Invalid webhook signature")
		return
	}

	// Parse the verified payload
	if err := json.Unmarshal(payload, &event); err != nil {
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
	if err := db.Preload("Laboratory").Where("invitation_token = ? AND status = ? AND expires_at > ?",
		invitationToken, "pending", time.Now()).First(&invitation).Error; err != nil {
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
		invitation.Status = "accepted"
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
		ID:             user.ID,
		ClerkUserID:    user.ClerkUserID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		OrganizationID: user.OrganizationID,
		LaboratoryID:   user.LaboratoryID,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt,
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
		ID:             user.ID,
		ClerkUserID:    user.ClerkUserID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		OrganizationID: user.OrganizationID,
		LaboratoryID:   user.LaboratoryID,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt,
	}

	render.JSON(w, r, userResponse)
}

// RefreshToken returns a fresh JWT token for the authenticated user
// This allows the frontend to refresh tokens without re-authenticating
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "No user claims found")
		return
	}

	// Get fresh session token from Clerk
	token, err := createClerkSession(claims.ClerkUserID)
	if err != nil {
		log.Printf("Failed to refresh token for user %s: %v", claims.ClerkUserID, err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to refresh token")
		return
	}

	db := database.GetDB()
	var user models.User
	if err := db.Where("clerk_user_id = ?", claims.ClerkUserID).First(&user).Error; err != nil {
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	authResponse := AuthResponse{
		Token: token,
		User: UserResponse{
			ID:             user.ID,
			ClerkUserID:    user.ClerkUserID,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			OrganizationID: user.OrganizationID,
			LaboratoryID:   user.LaboratoryID,
			Role:           user.Role,
			CreatedAt:      user.CreatedAt,
		},
	}

	render.JSON(w, r, authResponse)
}

// SignUp creates a new user with Clerk Backend API
func SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Email and password are required")
		return
	}

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Println("CLERK_SECRET_KEY environment variable not set")
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Server configuration error")
		return
	}

	// Create user with Clerk Backend API
	clerkAPIURL := "https://api.clerk.com/v1/users"

	payload := map[string]interface{}{
		"email_address": []string{req.Email},
		"password":      req.Password,
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
	}

	jsonPayload, _ := json.Marshal(payload)

	clerkReq, err := http.NewRequest("POST", clerkAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error creating Clerk request: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error creating user")
		return
	}

	clerkReq.Header.Set("Authorization", "Bearer "+clerkSecretKey)
	clerkReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(clerkReq)
	if err != nil {
		log.Printf("Error calling Clerk API: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error creating user")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Clerk API error (status %d): %s", resp.StatusCode, string(body))

		// Parse Clerk error response
		var clerkError map[string]interface{}
		if err := json.Unmarshal(body, &clerkError); err == nil {
			if errors, ok := clerkError["errors"].([]interface{}); ok && len(errors) > 0 {
				if errorObj, ok := errors[0].(map[string]interface{}); ok {
					if message, ok := errorObj["message"].(string); ok {
						middleware.ErrorResponse(w, r, http.StatusBadRequest, message)
						return
					}
				}
			}
		}

		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Failed to create user")
		return
	}

	var clerkUser ClerkUser
	if err := json.Unmarshal(body, &clerkUser); err != nil {
		log.Printf("Error parsing Clerk response: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Error parsing response")
		return
	}

	// Create user in our database (webhook will also do this, but we do it here for immediate response)
	db := database.GetDB()

	user := models.User{
		ClerkUserID: clerkUser.ID,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        "user",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error creating user in database: %v (user created in Clerk: %s)", err, clerkUser.ID)
		// User is created in Clerk, webhook will sync it later
	}

	// Create a session for the user to get the JWT token
	token, err := createClerkSession(clerkUser.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		// Return user without token - they can login separately
		render.JSON(w, r, map[string]interface{}{
			"message": "User created successfully. Please login.",
			"user": UserResponse{
				ID:             user.ID,
				ClerkUserID:    user.ClerkUserID,
				Email:          user.Email,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				OrganizationID: user.OrganizationID,
				LaboratoryID:   user.LaboratoryID,
				Role:           user.Role,
				CreatedAt:      user.CreatedAt,
			},
		})
		return
	}

	authResponse := AuthResponse{
		Token: token,
		User: UserResponse{
			ID:             user.ID,
			ClerkUserID:    user.ClerkUserID,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			OrganizationID: user.OrganizationID,
			LaboratoryID:   user.LaboratoryID,
			Role:           user.Role,
			CreatedAt:      user.CreatedAt,
		},
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, authResponse)
}

// SignIn authenticates a user with Clerk Backend API
func SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.Email == "" || req.Password == "" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "Email and password are required")
		return
	}

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Println("CLERK_SECRET_KEY environment variable not set")
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Server configuration error")
		return
	}

	// Get user by email first to get Clerk user ID
	db := database.GetDB()
	var user models.User

	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		log.Printf("Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	// Verify password with Clerk by creating a session
	token, err := verifyPasswordAndCreateSession(user.ClerkUserID, req.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	authResponse := AuthResponse{
		Token: token,
		User: UserResponse{
			ID:           user.ID,
			ClerkUserID:  user.ClerkUserID,
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			LaboratoryID: user.LaboratoryID,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt,
		},
	}

	render.JSON(w, r, authResponse)
}

// Helper function to verify password and create session
func verifyPasswordAndCreateSession(clerkUserID, password string) (string, error) {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		return "", fmt.Errorf("CLERK_SECRET_KEY not set")
	}

	// Verify the user's password by making an Admin API call
	// Note: Clerk doesn't have a direct "verify password" endpoint
	// We create a session which implicitly verifies credentials
	verifyURL := fmt.Sprintf("https://api.clerk.com/v1/users/%s/verify_password", clerkUserID)

	payload := map[string]interface{}{
		"password": password,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", verifyURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+clerkSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("password verification failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Password verified, now create a session
	return createClerkSession(clerkUserID)
}

// Helper function to create a Clerk session and return JWT token
func createClerkSession(userID string) (string, error) {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		return "", fmt.Errorf("CLERK_SECRET_KEY not set")
	}

	// Create session using Clerk Backend API
	clerkAPIURL := "https://api.clerk.com/v1/sessions"

	payload := map[string]interface{}{
		"user_id": userID,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", clerkAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+clerkSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("clerk API error (status %d): %s", resp.StatusCode, string(body))
	}

	var sessionResp ClerkSessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return "", err
	}

	// Get the JWT token for this session
	return getSessionToken(sessionResp.ID)
}

// Helper function to get JWT token from session
func getSessionToken(sessionID string) (string, error) {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		return "", fmt.Errorf("CLERK_SECRET_KEY not set")
	}

	tokenURL := fmt.Sprintf("https://api.clerk.com/v1/sessions/%s/tokens", sessionID)

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+clerkSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("clerk token API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp ClerkSessionTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.JWT, nil
}
