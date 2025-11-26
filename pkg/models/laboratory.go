package models

import (
	"fmt"
	"strings"
	"time"
	"gorm.io/gorm"
)

type Organization struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Laboratories []Laboratory `json:"laboratories,omitempty"`
	Users        []User       `json:"users,omitempty"`
}

type Laboratory struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	OrganizationID uint           `gorm:"not null;index" json:"organization_id"`
	Organization   *Organization  `json:"organization,omitempty"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Description    string         `json:"description"`
	Address        string         `gorm:"size:200" json:"address"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Invitations []LaboratoryInvitation `json:"invitations,omitempty"`
	Users       []User                 `json:"users,omitempty"`
}

type LaboratoryInvitation struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	LaboratoryID    uint           `gorm:"not null;index:idx_lab_email_status" json:"laboratory_id"`
	Laboratory      Laboratory     `json:"laboratory,omitempty"`
	Email           string         `gorm:"size:255;not null;index;index:idx_lab_email_status" json:"email"`
	Role            string         `gorm:"size:50;not null;check:role IN ('lab-coordinator', 'technician', 'student')" json:"role"`
	FirstName       string         `gorm:"size:100" json:"first_name"`
	LastName        string         `gorm:"size:100" json:"last_name"`
	InvitedBy       uint           `gorm:"not null" json:"invited_by"`
	InvitedByUser   *User          `gorm:"foreignKey:InvitedBy" json:"invited_by_user,omitempty"`
	InvitationToken string         `gorm:"size:255;unique;not null;index" json:"invitation_token"`
	Status          string         `gorm:"size:20;not null;default:'pending';check:status IN ('pending', 'accepted', 'expired', 'cancelled');index;index:idx_lab_email_status" json:"status"`
	ExpiresAt       time.Time      `gorm:"not null;index" json:"expires_at"`
	AcceptedAt      *time.Time     `json:"accepted_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to normalize email to lowercase
func (i *LaboratoryInvitation) BeforeCreate(tx *gorm.DB) error {
	i.Email = strings.ToLower(i.Email)
	return nil
}

// BeforeUpdate hook to normalize email to lowercase
func (i *LaboratoryInvitation) BeforeUpdate(tx *gorm.DB) error {
	i.Email = strings.ToLower(i.Email)
	return nil
}

type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	ClerkUserID    string         `gorm:"size:255;unique;not null" json:"clerk_user_id"`
	Email          string         `gorm:"size:255;unique;not null;index" json:"email"`
	FirstName      string         `gorm:"size:100" json:"first_name"`
	LastName       string         `gorm:"size:100" json:"last_name"`
	OrganizationID *uint          `json:"organization_id"`
	Organization   *Organization  `json:"organization,omitempty"`
	LaboratoryID   *uint          `json:"laboratory_id"`
	Laboratory     *Laboratory    `json:"laboratory,omitempty"`
	Role           string         `gorm:"size:50;check:role IN ('org-coordinator', 'lab-coordinator', 'technician', 'student')" json:"role"`
	Status         string         `gorm:"size:20;default:'active';check:status IN ('pending', 'active', 'suspended');index" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeSave hook to enforce mutual exclusivity between organization and laboratory
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.OrganizationID != nil && u.LaboratoryID != nil {
		return fmt.Errorf("user cannot belong to both organization and laboratory")
	}
	return nil
}

type Location struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Description   string         `json:"description"`
	SubLocationID *uint          `json:"sub_location_id"`
	SubLocation   *Location      `gorm:"foreignKey:SubLocationID" json:"sub_location,omitempty"`
	LaboratoryID  uint           `gorm:"not null" json:"laboratory_id"`
	Laboratory    Laboratory     `json:"laboratory,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Description    string         `json:"description"`
	LocationID     uint           `gorm:"not null" json:"location_id"`
	Location       Location       `json:"location,omitempty"`
	Quantity       float64        `gorm:"not null" json:"quantity"`
	Formula        string         `json:"formula"`
	Unit           string         `gorm:"size:50;not null" json:"unit"`
	ExpirationDate time.Time      `json:"expiration_date"`
	LaboratoryID   uint           `gorm:"not null" json:"laboratory_id"`
	Laboratory     Laboratory     `json:"laboratory,omitempty"`
	ProductUsages  []ProductUsage `json:"product_usages,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductUsage struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ProductID    uint           `gorm:"not null;index" json:"product_id"`
	Product      Product        `json:"product,omitempty"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	User         User           `json:"user,omitempty"`
	QuantityUsed float64        `gorm:"not null" json:"quantity_used"`
	Unit         string         `gorm:"size:50;not null" json:"unit"`
	UsedAt       time.Time      `gorm:"not null;index" json:"used_at"`
	Notes        string         `json:"notes"`
	LaboratoryID uint           `gorm:"not null;index" json:"laboratory_id"`
	Laboratory   Laboratory     `json:"laboratory,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Equipment struct {
	ID             uint              `gorm:"primarykey" json:"id"`
	Name           string            `gorm:"size:100;not null" json:"name"`
	Description    string            `json:"description"`
	LocationID     uint              `gorm:"not null" json:"location_id"`
	Location       Location          `json:"location,omitempty"`
	PropertyNumber string            `gorm:"size:100" json:"property_number"`
	LaboratoryID   uint              `gorm:"not null" json:"laboratory_id"`
	Laboratory     Laboratory        `json:"laboratory,omitempty"`
	EquipmentUsages []EquipmentUsage `json:"equipment_usages,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"-"`
}

type EquipmentUsage struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	User         User           `json:"user,omitempty"`
	EquipmentID  uint           `gorm:"not null" json:"equipment_id"`
	Equipment    Equipment      `json:"equipment,omitempty"`
	StartDate    time.Time      `gorm:"not null" json:"start_date"`
	EndDate      time.Time      `gorm:"not null" json:"end_date"`
	Description  string         `json:"description"`
	LaboratoryID uint           `gorm:"not null" json:"laboratory_id"`
	Laboratory   Laboratory     `json:"laboratory,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}