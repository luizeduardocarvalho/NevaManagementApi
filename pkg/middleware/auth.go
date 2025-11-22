package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
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

// JWKS cache
var (
	jwksCache      *JWKS
	jwksCacheMutex sync.RWMutex
	jwksCacheTime  time.Time
	jwksCacheTTL   = 1 * time.Hour
)

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

		// Look up user in database to get full claims
		db := database.GetDB()
		var user models.User
		if err := db.Where("clerk_user_id = ?", clerkClaims.Sub).First(&user).Error; err != nil {
			http.Error(w, "User not found in database", http.StatusUnauthorized)
			return
		}

		// Create claims from database user
		claims := &JWTClaims{
			UserID:       user.ID,
			ClerkUserID:  user.ClerkUserID,
			Email:        user.Email,
			Role:         user.Role,
		}

		// Add laboratory ID if user has one
		if user.LaboratoryID != nil {
			claims.LaboratoryID = *user.LaboratoryID
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
	// Get JWKS from cache or fetch from Clerk
	jwks, err := getJWKS()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return parseJWK(key)
		}
	}

	return nil, fmt.Errorf("public key with kid %s not found", kid)
}

// getJWKS fetches JWKS from Clerk with caching
func getJWKS() (*JWKS, error) {
	// Check cache first
	jwksCacheMutex.RLock()
	if jwksCache != nil && time.Since(jwksCacheTime) < jwksCacheTTL {
		defer jwksCacheMutex.RUnlock()
		return jwksCache, nil
	}
	jwksCacheMutex.RUnlock()

	// Fetch from Clerk
	jwksCacheMutex.Lock()
	defer jwksCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if jwksCache != nil && time.Since(jwksCacheTime) < jwksCacheTTL {
		return jwksCache, nil
	}

	// Get Clerk publishable key to construct JWKS URL
	clerkPublishableKey := os.Getenv("CLERK_PUBLISHABLE_KEY")
	if clerkPublishableKey == "" {
		return nil, fmt.Errorf("CLERK_PUBLISHABLE_KEY environment variable not set")
	}

	// Construct JWKS URL from publishable key
	// The publishable key is base64 encoded and contains the frontend API URL
	// Format: pk_test_[base64] or pk_live_[base64]
	var jwksURL string

	// Decode the publishable key to get the frontend API domain
	// Example: pk_test_c3RpcnJlZC1hbnQtNzEuY2xlcmsuYWNjb3VudHMuZGV2JA
	// The part after pk_test_ or pk_live_ is base64 encoded domain
	keyParts := strings.SplitN(clerkPublishableKey, "_", 3)
	if len(keyParts) != 3 {
		return nil, fmt.Errorf("invalid CLERK_PUBLISHABLE_KEY format (expected pk_test_xxx or pk_live_xxx)")
	}

	// Clerk uses URL-safe base64 encoding (without padding)
	decoded, err := base64.RawURLEncoding.DecodeString(keyParts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode CLERK_PUBLISHABLE_KEY: %w", err)
	}

	// Remove trailing $ if present
	domain := strings.TrimSuffix(string(decoded), "$")
	jwksURL = fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

	// Fetch JWKS
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("JWKS endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %w", err)
	}

	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// Update cache
	jwksCache = &jwks
	jwksCacheTime = time.Now()

	return &jwks, nil
}

// parseJWK converts a JWK to an RSA public key
func parseJWK(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode the exponent
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert bytes to big.Int
	n := new(big.Int).SetBytes(nBytes)

	// Convert exponent bytes to int
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

func GetUserClaims(r *http.Request) (*JWTClaims, bool) {
	claims, ok := r.Context().Value(UserClaimsKey).(*JWTClaims)
	return claims, ok
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, map[string]string{"error": message})
}