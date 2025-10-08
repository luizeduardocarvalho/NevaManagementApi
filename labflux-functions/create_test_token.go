package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Test token generator for development
func main() {
	// Use the same secret as your app (for testing only)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-local-jwt-secret-key"
	}

	// Create test claims (simulating a Clerk user)
	claims := jwt.MapClaims{
		"sub":   "user_test123",           // Clerk User ID
		"email": "test@example.com",       // User email
		"iat":   time.Now().Unix(),        // Issued at
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expires in 24 hours
		"iss":   "https://test.clerk.accounts.dev", // Issuer (Clerk)
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("Error creating token: %v\n", err)
		return
	}

	fmt.Println("🔑 Test JWT Token:")
	fmt.Println(tokenString)
	fmt.Println()
	fmt.Println("📋 Usage:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8080/api/auth/me\n", tokenString)
	fmt.Println()
	fmt.Println("⚠️  Note: This is for testing only. In production, use real Clerk tokens.")
}