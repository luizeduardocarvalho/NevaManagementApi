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