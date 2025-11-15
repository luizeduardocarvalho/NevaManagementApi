package models

import (
	"time"
	"gorm.io/gorm"
)

type Laboratory struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `json:"description"`
	Address     string         `gorm:"size:200" json:"address"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relations
	Invitations []LaboratoryInvitation `json:"invitations,omitempty"`
	Users       []User                 `json:"users,omitempty"`
}

type LaboratoryInvitation struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	LaboratoryID    uint           `gorm:"not null" json:"laboratory_id"`
	Laboratory      Laboratory     `json:"laboratory,omitempty"`
	Email           string         `gorm:"size:255;not null" json:"email"`
	InvitationToken string         `gorm:"size:36;unique;not null" json:"invitation_token"`
	Role            string         `gorm:"size:50;not null" json:"role"`
	ExpiresAt       time.Time      `gorm:"not null" json:"expires_at"`
	IsAccepted      bool           `gorm:"default:false" json:"is_accepted"`
	AcceptedAt      *time.Time     `json:"accepted_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ClerkUserID  string         `gorm:"size:255;unique;not null" json:"clerk_user_id"`
	Email        string         `gorm:"size:255;unique;not null" json:"email"`
	FirstName    string         `gorm:"size:100" json:"first_name"`
	LastName     string         `gorm:"size:100" json:"last_name"`
	LaboratoryID *uint          `json:"laboratory_id"`
	Laboratory   *Laboratory    `json:"laboratory,omitempty"`
	Role         string         `gorm:"size:50" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
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
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
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

type Researcher struct {
	ID              uint              `gorm:"primarykey" json:"id"`
	Name            string            `gorm:"size:80;not null" json:"name"`
	ClerkUserID     string            `gorm:"size:255;unique;not null" json:"clerk_user_id"`
	Email           string            `gorm:"size:255;unique;not null" json:"email"`
	LaboratoryID    uint              `gorm:"not null" json:"laboratory_id"`
	Laboratory      Laboratory        `json:"laboratory,omitempty"`
	EquipmentUsages []EquipmentUsage  `json:"equipment_usages,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}

type EquipmentUsage struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ResearcherID uint           `gorm:"not null" json:"researcher_id"`
	Researcher   Researcher     `json:"researcher,omitempty"`
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