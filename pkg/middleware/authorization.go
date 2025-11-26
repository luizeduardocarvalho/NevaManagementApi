package middleware

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// IsOrganizationCoordinator checks if user is an organization-level coordinator
func IsOrganizationCoordinator(claims *JWTClaims) bool {
	return claims.OrganizationID != 0 && claims.Role == "org-coordinator"
}

// IsLabCoordinator checks if user is a lab-level coordinator
func IsLabCoordinator(claims *JWTClaims) bool {
	return claims.LaboratoryID != 0 && claims.Role == "lab-coordinator"
}

// CanManageLab checks if user can manage a specific lab (coordinator of that lab)
func CanManageLab(claims *JWTClaims, labID uint) bool {
	if claims.LaboratoryID == labID && claims.Role == "lab-coordinator" {
		return true
	}
	return false
}

// CanViewLab checks if user can view a specific lab (includes org coordinators with labs in their org)
func CanViewLab(db *gorm.DB, claims *JWTClaims, labID uint) bool {
	// Lab-level users can only view their own lab
	if claims.LaboratoryID != 0 {
		return claims.LaboratoryID == labID
	}

	// Org coordinators can view any lab in their organization
	if IsOrganizationCoordinator(claims) {
		var count int64
		db.Table("laboratories").
			Where("id = ? AND organization_id = ? AND deleted_at IS NULL", labID, claims.OrganizationID).
			Count(&count)
		return count > 0
	}

	return false
}

// GetAccessibleLabIDs returns all lab IDs the user can access
func GetAccessibleLabIDs(db *gorm.DB, claims *JWTClaims) ([]uint, error) {
	var labIDs []uint

	// Lab-level users can only access their own lab
	if claims.LaboratoryID != 0 {
		return []uint{claims.LaboratoryID}, nil
	}

	// Org coordinators can access all labs in their organization
	if IsOrganizationCoordinator(claims) {
		if err := db.Table("laboratories").
			Where("organization_id = ? AND deleted_at IS NULL", claims.OrganizationID).
			Pluck("id", &labIDs).Error; err != nil {
			return nil, fmt.Errorf("failed to get accessible labs: %w", err)
		}
		return labIDs, nil
	}

	return nil, fmt.Errorf("user has no lab access")
}

// RequireOrganizationCoordinator middleware enforces org coordinator access
func RequireOrganizationCoordinator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r)
		if !ok || !IsOrganizationCoordinator(claims) {
			ErrorResponse(w, r, http.StatusForbidden, "Organization coordinator access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireLabCoordinator middleware enforces lab coordinator access
func RequireLabCoordinator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r)
		if !ok || !IsLabCoordinator(claims) {
			ErrorResponse(w, r, http.StatusForbidden, "Lab coordinator access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
