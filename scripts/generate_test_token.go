package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	LaboratoryID uint   `json:"laboratory_id"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// Create test claims
	claims := &JWTClaims{
		UserID:       1,
		Email:        "test@labflux.com",
		LaboratoryID: 1,
		Role:         "Admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "1",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := "your-local-jwt-secret-key"
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatal("Error creating token:", err)
	}

	fmt.Println("🔑 Test JWT Token Generated:")
	fmt.Println("==========================")
	fmt.Println(tokenString)
	fmt.Println("")
	fmt.Println("📋 Token Claims:")
	fmt.Printf("   User ID: %d\n", claims.UserID)
	fmt.Printf("   Email: %s\n", claims.Email)
	fmt.Printf("   Laboratory ID: %d\n", claims.LaboratoryID)
	fmt.Printf("   Role: %s\n", claims.Role)
	fmt.Printf("   Expires: %s\n", claims.ExpiresAt.Time.Format(time.RFC3339))
	fmt.Println("")
	fmt.Println("🧪 Example Usage:")
	fmt.Println("curl -H \"Authorization: Bearer " + tokenString + "\" \\")
	fmt.Println("     http://localhost:8080/api/laboratories/1")
}