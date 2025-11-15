package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type ClerkClaims struct {
	Sub   string `json:"sub"`   // Clerk user ID
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	ClerkUserID  string `json:"clerk_user_id"`
	Email        string `json:"email"`
	LaboratoryID uint   `json:"laboratory_id"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, fmt.Sprintf("Invalid authorization header format. Got: %s", authHeader), http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		
		// Try to validate as test token first (for development)
		if claims, err := validateTestToken(tokenString); err == nil {
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Try to validate as Clerk token
		clerkClaims, err := validateClerkToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Create claims from Clerk token
		claims := &JWTClaims{
			ClerkUserID: clerkClaims.Sub,
			Email:       clerkClaims.Email,
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateTestToken validates test tokens for development
func validateTestToken(tokenString string) (*JWTClaims, error) {
	secret := "your-local-jwt-secret-key" // Same as in create_test_token.go
	
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid test token: %v", err)
	}

	// Convert to our internal claims format
	userClaims := &JWTClaims{
		ClerkUserID: claims["sub"].(string),
		Email:       claims["email"].(string),
	}

	return userClaims, nil
}

// validateClerkToken validates a Clerk JWT token
func validateClerkToken(tokenString string) (*ClerkClaims, error) {
	// Parse without verification first to get the header
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &ClerkClaims{})
	if err != nil {
		return nil, err
	}

	// Get the key ID from header
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("kid not found in token header")
	}

	// Get Clerk public key (in production, cache this)
	publicKey, err := getClerkPublicKey(kid)
	if err != nil {
		return nil, err
	}

	// Parse and verify with public key
	claims := &ClerkClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// getClerkPublicKey fetches the public key from Clerk's JWKS endpoint
func getClerkPublicKey(kid string) (*rsa.PublicKey, error) {
	// In production, implement caching and actual HTTP request to Clerk JWKS endpoint
	// For now, return a placeholder error to indicate this needs implementation
	return nil, fmt.Errorf("clerk public key fetching not implemented - set CLERK_PUBLISHABLE_KEY")
}

func GetUserClaims(r *http.Request) (*JWTClaims, bool) {
	claims, ok := r.Context().Value(UserClaimsKey).(*JWTClaims)
	return claims, ok
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, map[string]string{"error": message})
}