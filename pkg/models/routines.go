package models

import (
	"time"

	"gorm.io/gorm"
)

// Routine represents a laboratory procedure/protocol that can be scheduled and executed
type Routine struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Description  string         `json:"description"`
	ScheduleType string         `gorm:"size:50;not null;check:schedule_type IN ('one_time', 'recurring', 'template')" json:"schedule_type"`
	Deadline     *time.Time     `json:"deadline,omitempty"`
	LaboratoryID uint           `gorm:"not null;index" json:"laboratory_id"`
	Laboratory   Laboratory     `json:"laboratory,omitempty"`
	CreatedBy    uint           `gorm:"not null" json:"created_by"`
	Creator      *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Steps      []RoutineStep      `json:"steps,omitempty"`
	Materials  []RoutineMaterial  `json:"materials,omitempty"`
	Equipment  []RoutineEquipment `json:"equipment,omitempty"`
	Recurrence *RecurrenceRule    `json:"recurrence,omitempty"`
	AssignedTo []User             `gorm:"many2many:routine_assignments;" json:"assigned_to,omitempty"`
	Executions []RoutineExecution `json:"executions,omitempty"`
}

// RoutineStep represents a step in a routine procedure
type RoutineStep struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	RoutineID   uint           `gorm:"not null;index" json:"routine_id"`
	Order       int            `gorm:"not null" json:"order"`
	Description string         `gorm:"not null" json:"description"`
	Notes       string         `json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RoutineMaterial represents a product/material required for a routine
type RoutineMaterial struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	RoutineID uint           `gorm:"not null;index" json:"routine_id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	Product   *Product       `json:"product,omitempty"`
	Quantity  float64        `gorm:"not null" json:"quantity"`
	Unit      string         `gorm:"size:50;not null" json:"unit"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RoutineEquipment represents equipment required for a routine
type RoutineEquipment struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	RoutineID         uint           `gorm:"not null;index" json:"routine_id"`
	EquipmentID       uint           `gorm:"not null" json:"equipment_id"`
	Equipment         *Equipment     `json:"equipment,omitempty"`
	EstimatedDuration int            `json:"estimated_duration"` // in minutes
	Required          bool           `gorm:"default:true" json:"required"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the table name used by RoutineEquipment to `routine_equipment`
func (RoutineEquipment) TableName() string {
	return "routine_equipment"
}

// RecurrenceRule defines how a routine recurs
type RecurrenceRule struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	RoutineID  uint           `gorm:"unique;not null" json:"routine_id"`
	Frequency  string         `gorm:"size:50;not null;check:frequency IN ('daily', 'weekly', 'monthly')" json:"frequency"`
	Interval   int            `gorm:"default:1" json:"interval"` // Every N days/weeks/months
	DaysOfWeek string         `json:"days_of_week"`              // JSON array for weekly: [0,1,2,3,4,5,6]
	DayOfMonth *int           `json:"day_of_month,omitempty"`    // For monthly
	StartDate  time.Time      `gorm:"not null" json:"start_date"`
	EndDate    *time.Time     `json:"end_date,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// RoutineExecution tracks the execution of a routine
type RoutineExecution struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	RoutineID    uint           `gorm:"not null;index" json:"routine_id"`
	Routine      *Routine       `json:"routine,omitempty"`
	ExecutedBy   uint           `gorm:"not null" json:"executed_by"`
	Executor     *User          `gorm:"foreignKey:ExecutedBy" json:"executor,omitempty"`
	Status       string         `gorm:"size:50;default:'in_progress';check:status IN ('in_progress', 'completed', 'cancelled')" json:"status"`
	StartedAt    time.Time      `gorm:"not null" json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
	Notes        string         `json:"notes"`
	LaboratoryID uint           `gorm:"not null;index" json:"laboratory_id"`
	Laboratory   Laboratory     `json:"laboratory,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	StepCompletions []RoutineExecutionStep     `gorm:"foreignKey:ExecutionID" json:"step_completions,omitempty"`
	Materials       []RoutineExecutionMaterial `gorm:"foreignKey:ExecutionID" json:"materials,omitempty"`
}

// RoutineExecutionStep tracks completion of individual steps during execution
type RoutineExecutionStep struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ExecutionID uint           `gorm:"not null;index" json:"execution_id"`
	StepID      uint           `gorm:"not null" json:"step_id"`
	Step        *RoutineStep   `json:"step,omitempty"`
	Completed   bool           `gorm:"default:false" json:"completed"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Notes       string         `json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RoutineExecutionMaterial tracks materials used during routine execution
type RoutineExecutionMaterial struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	ExecutionID     uint           `gorm:"not null;index" json:"execution_id"`
	ProductID       uint           `gorm:"not null" json:"product_id"`
	Product         *Product       `json:"product,omitempty"`
	PlannedQuantity float64        `gorm:"not null" json:"planned_quantity"`
	ActualQuantity  float64        `gorm:"default:0" json:"actual_quantity"`
	Unit            string         `gorm:"size:50;not null" json:"unit"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
