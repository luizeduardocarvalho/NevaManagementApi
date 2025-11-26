package models

import (
	"time"

	"gorm.io/gorm"
)

type Sample struct {
	ID                     uint           `gorm:"primarykey" json:"id"`
	Name                   string         `gorm:"size:255;not null" json:"name"`
	Description            string         `json:"description"`
	Origin                 string         `gorm:"size:255" json:"origin"`
	IsolationDate          *time.Time     `json:"isolation_date,omitempty"`
	Latitude               *float64       `json:"latitude,omitempty"`
	Longitude              *float64       `json:"longitude,omitempty"`
	SubcultureMedium       string         `gorm:"size:255" json:"subculture_medium"`
	SubcultureIntervalDays *int           `json:"subculture_interval_days,omitempty"`
	ResearcherID           *uint          `json:"researcher_id,omitempty"`
	Researcher             *User          `gorm:"foreignKey:ResearcherID" json:"researcher,omitempty"`
	LocationID             *uint          `json:"location_id,omitempty"`
	Location               *Location      `json:"location,omitempty"`
	Tags                   string         `gorm:"type:jsonb" json:"tags"`
	LaboratoryID           uint           `gorm:"not null;index" json:"laboratory_id"`
	Laboratory             Laboratory     `json:"laboratory,omitempty"`
	Replicas               []Replica      `json:"replicas,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

type Replica struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	Name               string         `gorm:"size:255;not null" json:"name"`
	SampleID           uint           `gorm:"not null;index" json:"sample_id"`
	Sample             *Sample        `json:"sample,omitempty"`
	LocationID         *uint          `json:"location_id,omitempty"`
	Location           *Location      `json:"location,omitempty"`
	Status             string         `gorm:"size:50;default:'active'" json:"status"`
	LastSubcultureDate *time.Time     `json:"last_subculture_date,omitempty"`
	NextSubcultureDate *time.Time     `json:"next_subculture_date,omitempty"`
	LaboratoryID       uint           `gorm:"not null;index" json:"laboratory_id"`
	Laboratory         Laboratory     `json:"laboratory,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
