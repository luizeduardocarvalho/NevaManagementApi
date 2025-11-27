package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

// Request/Response Types

type CreateRoutineRequest struct {
	Name         string                    `json:"name" validate:"required"`
	Description  string                    `json:"description"`
	Materials    []RoutineMaterialRequest  `json:"materials"`
	Equipment    []RoutineEquipmentRequest `json:"equipment"`
	Steps        []RoutineStepRequest      `json:"steps"`
	ScheduleType string                    `json:"scheduleType" validate:"required,oneof=one_time recurring template"`
	Recurrence   *RecurrenceRuleRequest    `json:"recurrence,omitempty"`
	Deadline     *string                   `json:"deadline,omitempty"`
	AssignedTo   []uint                    `json:"assignedTo,omitempty"`
}

type RoutineMaterialRequest struct {
	ProductID uint    `json:"productId" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

type RoutineEquipmentRequest struct {
	EquipmentID       uint `json:"equipmentId" validate:"required"`
	EstimatedDuration int  `json:"estimatedDuration"`
	Required          bool `json:"required"`
}

type RoutineStepRequest struct {
	Order       int    `json:"order" validate:"required"`
	Description string `json:"description" validate:"required"`
	Notes       string `json:"notes"`
}

type RecurrenceRuleRequest struct {
	Frequency  string  `json:"frequency" validate:"required,oneof=daily weekly monthly"`
	Interval   int     `json:"interval"`
	DaysOfWeek []int   `json:"daysOfWeek,omitempty"`
	DayOfMonth *int    `json:"dayOfMonth,omitempty"`
	StartDate  string  `json:"startDate" validate:"required"`
	EndDate    *string `json:"endDate,omitempty"`
}

type RoutineResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	ScheduleType string                     `json:"scheduleType"`
	Deadline     *time.Time                 `json:"deadline,omitempty"`
	Materials    []RoutineMaterialResponse  `json:"materials"`
	Equipment    []RoutineEquipmentResponse `json:"equipment"`
	Steps        []RoutineStepResponse      `json:"steps"`
	Recurrence   *RecurrenceRuleResponse    `json:"recurrence,omitempty"`
	AssignedTo   []uint                     `json:"assignedTo"`
	LaboratoryID uint                       `json:"laboratory_id"`
	CreatedBy    uint                       `json:"created_by"`
	CreatedAt    time.Time                  `json:"createdAt"`
}

type RoutineMaterialResponse struct {
	ProductID   uint    `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

type RoutineEquipmentResponse struct {
	EquipmentID       uint   `json:"equipmentId"`
	EquipmentName     string `json:"equipmentName"`
	EstimatedDuration int    `json:"estimatedDuration"`
	Required          bool   `json:"required"`
}

type RoutineStepResponse struct {
	ID          uint   `json:"id"`
	Order       int    `json:"order"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

type RecurrenceRuleResponse struct {
	Frequency  string     `json:"frequency"`
	Interval   int        `json:"interval"`
	DaysOfWeek []int      `json:"daysOfWeek,omitempty"`
	DayOfMonth *int       `json:"dayOfMonth,omitempty"`
	StartDate  time.Time  `json:"startDate"`
	EndDate    *time.Time `json:"endDate,omitempty"`
}

// Execution request/response types

type StartExecutionRequest struct {
	Notes string `json:"notes,omitempty"`
}

type ExecutionResponse struct {
	ID              uint                        `json:"id"`
	RoutineID       uint                        `json:"routineId"`
	RoutineName     string                      `json:"routineName"`
	ExecutedBy      uint                        `json:"executedBy"`
	ExecutorName    string                      `json:"executorName"`
	Status          string                      `json:"status"`
	StartedAt       time.Time                   `json:"startedAt"`
	CompletedAt     *time.Time                  `json:"completedAt,omitempty"`
	Notes           string                      `json:"notes"`
	LaboratoryID    uint                        `json:"laboratory_id"`
	StepCompletions []ExecutionStepResponse     `json:"stepCompletions"`
	Materials       []ExecutionMaterialResponse `json:"materials"`
}

type ExecutionStepResponse struct {
	ID          uint       `json:"id"`
	StepID      uint       `json:"stepId"`
	Description string     `json:"description"`
	Order       int        `json:"order"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	Notes       string     `json:"notes"`
}

type ExecutionMaterialResponse struct {
	ID              uint    `json:"id"`
	ProductID       uint    `json:"productId"`
	ProductName     string  `json:"productName"`
	PlannedQuantity float64 `json:"plannedQuantity"`
	ActualQuantity  float64 `json:"actualQuantity"`
	Unit            string  `json:"unit"`
}

type MarkStepCompleteRequest struct {
	Notes string `json:"notes,omitempty"`
}

type CompletionMaterial struct {
	ProductID      uint    `json:"productId" validate:"required"`
	ActualQuantity float64 `json:"actualQuantity" validate:"required,gte=0"`
}

type CompleteExecutionRequest struct {
	Notes     string               `json:"notes,omitempty"`
	Materials []CompletionMaterial `json:"materials,omitempty"`
}

// Helper functions

func parseISODate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// Try multiple date formats to support different frontend inputs
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // ISO without timezone
		"2006-01-02T15:04",    // datetime-local format (HTML5)
		"2006-01-02",          // date only
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, dateStr)
		if err == nil {
			return &t, nil
		}
	}

	return nil, err
}

func buildRoutineResponse(routine models.Routine) RoutineResponse {
	response := RoutineResponse{
		ID:           routine.ID,
		Name:         routine.Name,
		Description:  routine.Description,
		ScheduleType: routine.ScheduleType,
		Deadline:     routine.Deadline,
		LaboratoryID: routine.LaboratoryID,
		CreatedBy:    routine.CreatedBy,
		CreatedAt:    routine.CreatedAt,
		Materials:    []RoutineMaterialResponse{},
		Equipment:    []RoutineEquipmentResponse{},
		Steps:        []RoutineStepResponse{},
		AssignedTo:   []uint{},
	}

	for _, m := range routine.Materials {
		productName := ""
		if m.Product != nil {
			productName = m.Product.Name
		}
		response.Materials = append(response.Materials, RoutineMaterialResponse{
			ProductID:   m.ProductID,
			ProductName: productName,
			Quantity:    m.Quantity,
			Unit:        m.Unit,
		})
	}

	for _, e := range routine.Equipment {
		equipmentName := ""
		if e.Equipment != nil {
			equipmentName = e.Equipment.Name
		}
		response.Equipment = append(response.Equipment, RoutineEquipmentResponse{
			EquipmentID:       e.EquipmentID,
			EquipmentName:     equipmentName,
			EstimatedDuration: e.EstimatedDuration,
			Required:          e.Required,
		})
	}

	for _, s := range routine.Steps {
		response.Steps = append(response.Steps, RoutineStepResponse{
			ID:          s.ID,
			Order:       s.Order,
			Description: s.Description,
			Notes:       s.Notes,
		})
	}

	if routine.Recurrence != nil {
		var daysOfWeek []int
		if routine.Recurrence.DaysOfWeek != "" {
			json.Unmarshal([]byte(routine.Recurrence.DaysOfWeek), &daysOfWeek)
		}
		response.Recurrence = &RecurrenceRuleResponse{
			Frequency:  routine.Recurrence.Frequency,
			Interval:   routine.Recurrence.Interval,
			DaysOfWeek: daysOfWeek,
			DayOfMonth: routine.Recurrence.DayOfMonth,
			StartDate:  routine.Recurrence.StartDate,
			EndDate:    routine.Recurrence.EndDate,
		}
	}

	for _, user := range routine.AssignedTo {
		response.AssignedTo = append(response.AssignedTo, user.ID)
	}

	return response
}

func buildExecutionResponse(execution models.RoutineExecution) ExecutionResponse {
	response := ExecutionResponse{
		ID:              execution.ID,
		RoutineID:       execution.RoutineID,
		RoutineName:     "",
		ExecutedBy:      execution.ExecutedBy,
		ExecutorName:    "",
		Status:          execution.Status,
		StartedAt:       execution.StartedAt,
		CompletedAt:     execution.CompletedAt,
		Notes:           execution.Notes,
		LaboratoryID:    execution.LaboratoryID,
		StepCompletions: []ExecutionStepResponse{},
		Materials:       []ExecutionMaterialResponse{},
	}

	if execution.Routine != nil {
		response.RoutineName = execution.Routine.Name
	}

	if execution.Executor != nil {
		response.ExecutorName = execution.Executor.FirstName + " " + execution.Executor.LastName
	}

	for _, stepCompletion := range execution.StepCompletions {
		stepResp := ExecutionStepResponse{
			ID:          stepCompletion.ID,
			StepID:      stepCompletion.StepID,
			Completed:   stepCompletion.Completed,
			CompletedAt: stepCompletion.CompletedAt,
			Notes:       stepCompletion.Notes,
		}

		if stepCompletion.Step != nil {
			stepResp.Description = stepCompletion.Step.Description
			stepResp.Order = stepCompletion.Step.Order
		}

		response.StepCompletions = append(response.StepCompletions, stepResp)
	}

	for _, material := range execution.Materials {
		productName := ""
		if material.Product != nil {
			productName = material.Product.Name
		}

		materialResp := ExecutionMaterialResponse{
			ID:              material.ID,
			ProductID:       material.ProductID,
			ProductName:     productName,
			PlannedQuantity: material.PlannedQuantity,
			ActualQuantity:  material.ActualQuantity,
			Unit:            material.Unit,
		}

		response.Materials = append(response.Materials, materialResp)
	}

	return response
}

// GET /api/routines
func GetAllRoutines(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	db := database.GetDB()
	query := db.Where("laboratory_id = ?", claims.LaboratoryID)

	// Filter by schedule type if provided
	scheduleType := r.URL.Query().Get("scheduleType")
	if scheduleType != "" {
		query = query.Where("schedule_type = ?", scheduleType)
	}

	var routines []models.Routine
	if err := query.
		Preload("Steps").
		Preload("Materials.Product").
		Preload("Equipment.Equipment").
		Preload("Recurrence").
		Preload("AssignedTo").
		Order("created_at DESC").
		Find(&routines).Error; err != nil {
		log.Printf("[ROUTINES] Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch routines")
		return
	}

	response := make([]RoutineResponse, 0, len(routines))
	for _, routine := range routines {
		response = append(response, buildRoutineResponse(routine))
	}

	render.JSON(w, r, response)
}

// GET /api/routines/:id
func GetRoutineByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	routineIDStr := chi.URLParam(r, "id")
	routineID, err := strconv.ParseUint(routineIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid routine ID")
		return
	}

	db := database.GetDB()
	var routine models.Routine

	if err := db.Where("id = ? AND laboratory_id = ?", routineID, claims.LaboratoryID).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Materials.Product").
		Preload("Equipment.Equipment").
		Preload("Recurrence").
		Preload("AssignedTo").
		First(&routine).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "routine not found")
			return
		}
		log.Printf("[ROUTINES] Database error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	render.JSON(w, r, buildRoutineResponse(routine))
}

// POST /api/routines
func CreateRoutine(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	var req CreateRoutineRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate schedule type
	validScheduleTypes := map[string]bool{"one_time": true, "recurring": true, "template": true}
	if !validScheduleTypes[req.ScheduleType] {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid schedule type")
		return
	}

	// Parse deadline
	deadline, err := parseISODate(stringOrEmpty(req.Deadline))
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid deadline format")
		return
	}

	db := database.GetDB()

	// Create routine in transaction
	var createdRoutine models.Routine
	err = db.Transaction(func(tx *gorm.DB) error {
		routine := models.Routine{
			Name:         req.Name,
			Description:  req.Description,
			ScheduleType: req.ScheduleType,
			Deadline:     deadline,
			LaboratoryID: claims.LaboratoryID,
			CreatedBy:    claims.UserID,
		}

		if err := tx.Create(&routine).Error; err != nil {
			return err
		}

		// Create steps
		for _, stepReq := range req.Steps {
			step := models.RoutineStep{
				RoutineID:   routine.ID,
				Order:       stepReq.Order,
				Description: stepReq.Description,
				Notes:       stepReq.Notes,
			}
			if err := tx.Create(&step).Error; err != nil {
				return err
			}
		}

		// Create materials
		for _, matReq := range req.Materials {
			// Get product unit
			var product models.Product
			if err := tx.Select("unit").Where("id = ?", matReq.ProductID).First(&product).Error; err != nil {
				return err
			}

			material := models.RoutineMaterial{
				RoutineID: routine.ID,
				ProductID: matReq.ProductID,
				Quantity:  matReq.Quantity,
				Unit:      product.Unit,
			}
			if err := tx.Create(&material).Error; err != nil {
				return err
			}
		}

		// Create equipment
		for _, eqReq := range req.Equipment {
			equipment := models.RoutineEquipment{
				RoutineID:         routine.ID,
				EquipmentID:       eqReq.EquipmentID,
				EstimatedDuration: eqReq.EstimatedDuration,
				Required:          eqReq.Required,
			}
			if err := tx.Create(&equipment).Error; err != nil {
				return err
			}
		}

		// Create recurrence rule if recurring
		if req.ScheduleType == "recurring" && req.Recurrence != nil {
			startDate, err := parseISODate(req.Recurrence.StartDate)
			if err != nil {
				return err
			}
			endDate, err := parseISODate(stringOrEmpty(req.Recurrence.EndDate))
			if err != nil {
				return err
			}

			daysOfWeekJSON, _ := json.Marshal(req.Recurrence.DaysOfWeek)

			recurrence := models.RecurrenceRule{
				RoutineID:  routine.ID,
				Frequency:  req.Recurrence.Frequency,
				Interval:   req.Recurrence.Interval,
				DaysOfWeek: string(daysOfWeekJSON),
				DayOfMonth: req.Recurrence.DayOfMonth,
				StartDate:  *startDate,
				EndDate:    endDate,
			}
			if err := tx.Create(&recurrence).Error; err != nil {
				return err
			}
		}

		// Assign users
		if len(req.AssignedTo) > 0 {
			for _, userID := range req.AssignedTo {
				if err := tx.Exec("INSERT INTO routine_assignments (routine_id, user_id) VALUES (?, ?)", routine.ID, userID).Error; err != nil {
					return err
				}
			}
		}

		// Fetch complete routine
		if err := tx.Where("id = ?", routine.ID).
			Preload("Steps", func(db *gorm.DB) *gorm.DB {
				return db.Order("\"order\" ASC")
			}).
			Preload("Materials.Product").
			Preload("Equipment.Equipment").
			Preload("Recurrence").
			Preload("AssignedTo").
			First(&createdRoutine).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("[ROUTINES] Transaction error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to create routine")
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, buildRoutineResponse(createdRoutine))
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// PUT /api/routines/:id
func UpdateRoutine(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	routineIDStr := chi.URLParam(r, "id")
	routineID, err := strconv.ParseUint(routineIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid routine ID")
		return
	}

	var req CreateRoutineRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	deadline, err := parseISODate(stringOrEmpty(req.Deadline))
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid deadline format")
		return
	}

	db := database.GetDB()

	// Update routine in transaction
	var updatedRoutine models.Routine
	err = db.Transaction(func(tx *gorm.DB) error {
		// Find routine
		var routine models.Routine
		if err := tx.Where("id = ? AND laboratory_id = ?", routineID, claims.LaboratoryID).First(&routine).Error; err != nil {
			return err
		}

		// Update basic fields
		routine.Name = req.Name
		routine.Description = req.Description
		routine.ScheduleType = req.ScheduleType
		routine.Deadline = deadline

		if err := tx.Save(&routine).Error; err != nil {
			return err
		}

		// Delete and recreate related entities
		tx.Where("routine_id = ?", routineID).Delete(&models.RoutineStep{})
		tx.Where("routine_id = ?", routineID).Delete(&models.RoutineMaterial{})
		tx.Where("routine_id = ?", routineID).Delete(&models.RoutineEquipment{})
		tx.Exec("DELETE FROM routine_assignments WHERE routine_id = ?", routineID)

		// Recreate steps
		for _, stepReq := range req.Steps {
			step := models.RoutineStep{
				RoutineID:   routine.ID,
				Order:       stepReq.Order,
				Description: stepReq.Description,
				Notes:       stepReq.Notes,
			}
			if err := tx.Create(&step).Error; err != nil {
				return err
			}
		}

		// Recreate materials
		for _, matReq := range req.Materials {
			var product models.Product
			if err := tx.Select("unit").Where("id = ?", matReq.ProductID).First(&product).Error; err != nil {
				continue
			}
			material := models.RoutineMaterial{
				RoutineID: routine.ID,
				ProductID: matReq.ProductID,
				Quantity:  matReq.Quantity,
				Unit:      product.Unit,
			}
			if err := tx.Create(&material).Error; err != nil {
				return err
			}
		}

		// Recreate equipment
		for _, eqReq := range req.Equipment {
			equipment := models.RoutineEquipment{
				RoutineID:         routine.ID,
				EquipmentID:       eqReq.EquipmentID,
				EstimatedDuration: eqReq.EstimatedDuration,
				Required:          eqReq.Required,
			}
			if err := tx.Create(&equipment).Error; err != nil {
				return err
			}
		}

		// Update recurrence
		tx.Where("routine_id = ?", routineID).Delete(&models.RecurrenceRule{})
		if req.ScheduleType == "recurring" && req.Recurrence != nil {
			startDate, err := parseISODate(req.Recurrence.StartDate)
			if err != nil {
				return err
			}
			endDate, err := parseISODate(stringOrEmpty(req.Recurrence.EndDate))
			if err != nil {
				return err
			}

			daysOfWeekJSON, _ := json.Marshal(req.Recurrence.DaysOfWeek)
			recurrence := models.RecurrenceRule{
				RoutineID:  routine.ID,
				Frequency:  req.Recurrence.Frequency,
				Interval:   req.Recurrence.Interval,
				DaysOfWeek: string(daysOfWeekJSON),
				DayOfMonth: req.Recurrence.DayOfMonth,
				StartDate:  *startDate,
				EndDate:    endDate,
			}
			if err := tx.Create(&recurrence).Error; err != nil {
				return err
			}
		}

		// Reassign users
		if len(req.AssignedTo) > 0 {
			for _, userID := range req.AssignedTo {
				tx.Exec("INSERT INTO routine_assignments (routine_id, user_id) VALUES (?, ?)", routine.ID, userID)
			}
		}

		// Fetch updated routine
		if err := tx.Where("id = ?", routine.ID).
			Preload("Steps", func(db *gorm.DB) *gorm.DB {
				return db.Order("\"order\" ASC")
			}).
			Preload("Materials.Product").
			Preload("Equipment.Equipment").
			Preload("Recurrence").
			Preload("AssignedTo").
			First(&updatedRoutine).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "routine not found")
			return
		}
		log.Printf("[ROUTINES] Update error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to update routine")
		return
	}

	render.JSON(w, r, buildRoutineResponse(updatedRoutine))
}

// DELETE /api/routines/:id
func DeleteRoutine(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	routineIDStr := chi.URLParam(r, "id")
	routineID, err := strconv.ParseUint(routineIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid routine ID")
		return
	}

	db := database.GetDB()

	var routine models.Routine
	if err := db.Where("id = ? AND laboratory_id = ?", routineID, claims.LaboratoryID).First(&routine).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "routine not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	if err := db.Delete(&routine).Error; err != nil {
		log.Printf("[ROUTINES] Delete error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to delete routine")
		return
	}

	render.Status(r, http.StatusNoContent)
	w.Write(nil)
}

// POST /api/routines/:routineId/executions - Start execution
func StartExecution(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	routineIDStr := chi.URLParam(r, "routineId")
	routineID, err := strconv.ParseUint(routineIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid routine ID")
		return
	}

	var req StartExecutionRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()

	// Verify routine exists and belongs to lab
	var routine models.Routine
	if err := db.Where("id = ? AND laboratory_id = ?", routineID, claims.LaboratoryID).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Materials.Product").
		Preload("Equipment.Equipment").
		First(&routine).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "routine not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	var createdExecution models.RoutineExecution
	err = db.Transaction(func(tx *gorm.DB) error {
		// Create execution
		execution := models.RoutineExecution{
			RoutineID:    uint(routineID),
			ExecutedBy:   claims.UserID,
			Status:       "in_progress",
			StartedAt:    time.Now(),
			Notes:        req.Notes,
			LaboratoryID: claims.LaboratoryID,
		}

		if err := tx.Create(&execution).Error; err != nil {
			return err
		}

		// Create step completion records for all steps
		for _, step := range routine.Steps {
			stepCompletion := models.RoutineExecutionStep{
				ExecutionID: execution.ID,
				StepID:      step.ID,
				Completed:   false,
			}
			if err := tx.Create(&stepCompletion).Error; err != nil {
				return err
			}
		}

		// Block equipment - create equipment_usage records
		now := time.Now()
		for _, eq := range routine.Equipment {
			endDate := now.Add(time.Duration(eq.EstimatedDuration) * time.Minute)
			equipmentUsage := models.EquipmentUsage{
				UserID:       claims.UserID,
				EquipmentID:  eq.EquipmentID,
				StartDate:    now,
				EndDate:      endDate,
				Description:  fmt.Sprintf("Routine: %s (Execution #%d)", routine.Name, execution.ID),
				LaboratoryID: claims.LaboratoryID,
			}
			if err := tx.Create(&equipmentUsage).Error; err != nil {
				return fmt.Errorf("failed to create equipment usage: %w", err)
			}
		}

		// Reserve materials and create execution materials records
		for _, mat := range routine.Materials {
			// Check availability
			var product models.Product
			if err := tx.Where("id = ?", mat.ProductID).First(&product).Error; err != nil {
				return fmt.Errorf("failed to fetch product %d: %w", mat.ProductID, err)
			}

			availableQty := product.Quantity - product.ReservedQuantity
			if availableQty < mat.Quantity {
				return fmt.Errorf("insufficient quantity for product %s: need %.2f %s, available %.2f %s",
					product.Name, mat.Quantity, mat.Unit, availableQty, product.Unit)
			}

			// Reserve the material
			if err := tx.Model(&models.Product{}).
				Where("id = ?", mat.ProductID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", mat.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to reserve material: %w", err)
			}

			// Create execution material record
			execMaterial := models.RoutineExecutionMaterial{
				ExecutionID:     execution.ID,
				ProductID:       mat.ProductID,
				PlannedQuantity: mat.Quantity,
				ActualQuantity:  0, // Will be filled on completion
				Unit:            mat.Unit,
			}
			if err := tx.Create(&execMaterial).Error; err != nil {
				return fmt.Errorf("failed to create execution material: %w", err)
			}
		}

		// Fetch complete execution with relations
		if err := tx.Where("id = ?", execution.ID).
			Preload("Routine").
			Preload("Executor").
			Preload("StepCompletions.Step", func(db *gorm.DB) *gorm.DB {
				return db.Order("\"order\" ASC")
			}).
			Preload("Materials.Product").
			First(&createdExecution).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("[ROUTINES] Start execution error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to start execution")
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, buildExecutionResponse(createdExecution))
}

// GET /api/routines/executions/:executionId - Get execution details
func GetExecution(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	executionIDStr := chi.URLParam(r, "executionId")
	executionID, err := strconv.ParseUint(executionIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid execution ID")
		return
	}

	db := database.GetDB()
	var execution models.RoutineExecution

	if err := db.Where("id = ? AND laboratory_id = ?", executionID, claims.LaboratoryID).
		Preload("Routine").
		Preload("Executor").
		Preload("StepCompletions.Step", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Materials.Product").
		First(&execution).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "execution not found")
			return
		}
		log.Printf("[ROUTINES] Get execution error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	render.JSON(w, r, buildExecutionResponse(execution))
}

// PATCH /api/routines/executions/:executionId/steps/:stepId - Mark step complete
func MarkStepComplete(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	executionIDStr := chi.URLParam(r, "executionId")
	executionID, err := strconv.ParseUint(executionIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid execution ID")
		return
	}

	stepIDStr := chi.URLParam(r, "stepId")
	stepID, err := strconv.ParseUint(stepIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid step ID")
		return
	}

	var req MarkStepCompleteRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()

	// Verify execution belongs to lab
	var execution models.RoutineExecution
	if err := db.Where("id = ? AND laboratory_id = ?", executionID, claims.LaboratoryID).First(&execution).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "execution not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	if execution.Status != "in_progress" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "execution is not in progress")
		return
	}

	// Update step completion
	var stepCompletion models.RoutineExecutionStep
	if err := db.Where("execution_id = ? AND step_id = ?", executionID, stepID).First(&stepCompletion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "step not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	now := time.Now()
	stepCompletion.Completed = true
	stepCompletion.CompletedAt = &now
	stepCompletion.Notes = req.Notes

	if err := db.Save(&stepCompletion).Error; err != nil {
		log.Printf("[ROUTINES] Mark step complete error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to mark step complete")
		return
	}

	// Fetch updated execution
	if err := db.Where("id = ?", executionID).
		Preload("Routine").
		Preload("Executor").
		Preload("StepCompletions.Step", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		First(&execution).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	render.JSON(w, r, buildExecutionResponse(execution))
}

// POST /api/routines/executions/:executionId/complete - Complete execution
func CompleteExecution(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	executionIDStr := chi.URLParam(r, "executionId")
	executionID, err := strconv.ParseUint(executionIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid execution ID")
		return
	}

	var req CompleteExecutionRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	db := database.GetDB()

	var execution models.RoutineExecution
	if err := db.Where("id = ? AND laboratory_id = ?", executionID, claims.LaboratoryID).
		Preload("Routine").
		First(&execution).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "execution not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	if execution.Status != "in_progress" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "execution is not in progress")
		return
	}

	// Process completion in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update execution status
		execution.Status = "completed"
		execution.CompletedAt = &now
		if req.Notes != "" {
			execution.Notes = req.Notes
		}

		if err := tx.Save(&execution).Error; err != nil {
			return fmt.Errorf("failed to update execution: %w", err)
		}

		// Process materials if provided
		if len(req.Materials) > 0 {
			// Create a map for quick lookup
			materialMap := make(map[uint]float64)
			for _, mat := range req.Materials {
				materialMap[mat.ProductID] = mat.ActualQuantity
			}

			// Fetch all execution materials
			var execMaterials []models.RoutineExecutionMaterial
			if err := tx.Where("execution_id = ?", execution.ID).Find(&execMaterials).Error; err != nil {
				return fmt.Errorf("failed to fetch execution materials: %w", err)
			}

			// Update each material
			for _, execMat := range execMaterials {
				actualQty, provided := materialMap[execMat.ProductID]
				if !provided {
					// Use planned quantity if not provided
					actualQty = execMat.PlannedQuantity
				}

				// Update actual quantity in execution_materials
				if err := tx.Model(&models.RoutineExecutionMaterial{}).
					Where("id = ?", execMat.ID).
					Update("actual_quantity", actualQty).Error; err != nil {
					return fmt.Errorf("failed to update execution material: %w", err)
				}

				// Consume from product inventory
				if err := tx.Model(&models.Product{}).
					Where("id = ?", execMat.ProductID).
					Update("quantity", gorm.Expr("quantity - ?", actualQty)).Error; err != nil {
					return fmt.Errorf("failed to consume product: %w", err)
				}

				// Release reserved quantity
				if err := tx.Model(&models.Product{}).
					Where("id = ?", execMat.ProductID).
					Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", execMat.PlannedQuantity)).Error; err != nil {
					return fmt.Errorf("failed to release reserved quantity: %w", err)
				}
			}
		}

		// Update equipment usage end_date to actual completion time
		if err := tx.Exec(`
			UPDATE equipment_usages
			SET end_date = ?
			WHERE description LIKE ?
			AND laboratory_id = ?
		`, now, fmt.Sprintf("%%Execution #%d%%", execution.ID), claims.LaboratoryID).Error; err != nil {
			return fmt.Errorf("failed to update equipment usage: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("[ROUTINES] Complete execution error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to complete execution")
		return
	}

	// Fetch updated execution
	if err := db.Where("id = ?", executionID).
		Preload("Routine").
		Preload("Executor").
		Preload("StepCompletions.Step", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Materials.Product").
		First(&execution).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	render.JSON(w, r, buildExecutionResponse(execution))
}

// POST /api/routines/executions/:executionId/cancel - Cancel execution
func CancelExecution(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	executionIDStr := chi.URLParam(r, "executionId")
	executionID, err := strconv.ParseUint(executionIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid execution ID")
		return
	}

	db := database.GetDB()

	var execution models.RoutineExecution
	if err := db.Where("id = ? AND laboratory_id = ?", executionID, claims.LaboratoryID).First(&execution).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "execution not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	if execution.Status != "in_progress" {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "execution is not in progress")
		return
	}

	// Cancel execution and release resources in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		// Update execution status
		execution.Status = "cancelled"
		if err := tx.Save(&execution).Error; err != nil {
			return fmt.Errorf("failed to update execution: %w", err)
		}

		// Release reserved materials
		var execMaterials []models.RoutineExecutionMaterial
		if err := tx.Where("execution_id = ?", execution.ID).Find(&execMaterials).Error; err != nil {
			return fmt.Errorf("failed to fetch execution materials: %w", err)
		}

		for _, execMat := range execMaterials {
			// Release reserved quantity
			if err := tx.Model(&models.Product{}).
				Where("id = ?", execMat.ProductID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", execMat.PlannedQuantity)).Error; err != nil {
				return fmt.Errorf("failed to release reserved quantity: %w", err)
			}
		}

		// Delete execution materials records
		if err := tx.Where("execution_id = ?", execution.ID).Delete(&models.RoutineExecutionMaterial{}).Error; err != nil {
			return fmt.Errorf("failed to delete execution materials: %w", err)
		}

		// Delete equipment usage records
		if err := tx.Exec(`
			DELETE FROM equipment_usages
			WHERE description LIKE ?
			AND laboratory_id = ?
		`, fmt.Sprintf("%%Execution #%d%%", execution.ID), claims.LaboratoryID).Error; err != nil {
			return fmt.Errorf("failed to delete equipment usage: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("[ROUTINES] Cancel execution error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to cancel execution")
		return
	}

	// Fetch updated execution
	if err := db.Where("id = ?", executionID).
		Preload("Routine").
		Preload("Executor").
		Preload("StepCompletions.Step", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Materials.Product").
		First(&execution).Error; err != nil {
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	render.JSON(w, r, buildExecutionResponse(execution))
}

// GET /api/routines/:routineId/availability - Check equipment/resource availability
func GetRoutineAvailability(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	routineIDStr := chi.URLParam(r, "routineId")
	routineID, err := strconv.ParseUint(routineIDStr, 10, 32)
	if err != nil {
		middleware.ErrorResponse(w, r, http.StatusBadRequest, "invalid routine ID")
		return
	}

	db := database.GetDB()

	// Verify routine exists and belongs to lab
	var routine models.Routine
	if err := db.Where("id = ? AND laboratory_id = ?", routineID, claims.LaboratoryID).
		Preload("Equipment.Equipment").
		Preload("Materials.Product").
		First(&routine).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(w, r, http.StatusNotFound, "routine not found")
			return
		}
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
		return
	}

	// Check equipment availability (check if equipment is currently in use)
	equipmentAvailability := make([]map[string]interface{}, 0)
	now := time.Now()
	for _, eq := range routine.Equipment {
		if eq.Equipment != nil {
			// Check if equipment has any active usage
			var activeUsageCount int64
			db.Model(&models.EquipmentUsage{}).
				Where("equipment_id = ? AND start_date <= ? AND end_date >= ?", eq.EquipmentID, now, now).
				Count(&activeUsageCount)

			available := activeUsageCount == 0
			status := "available"
			if !available {
				status = "in_use"
			}

			equipmentAvailability = append(equipmentAvailability, map[string]interface{}{
				"equipmentId":   eq.EquipmentID,
				"equipmentName": eq.Equipment.Name,
				"available":     available,
				"status":        status,
				"required":      eq.Required,
			})
		}
	}

	// Check material availability (stock levels)
	materialAvailability := make([]map[string]interface{}, 0)
	for _, mat := range routine.Materials {
		if mat.Product != nil {
			sufficient := mat.Product.Quantity >= mat.Quantity
			materialAvailability = append(materialAvailability, map[string]interface{}{
				"productId":    mat.ProductID,
				"productName":  mat.Product.Name,
				"requiredQty":  mat.Quantity,
				"availableQty": mat.Product.Quantity,
				"sufficient":   sufficient,
				"unit":         mat.Unit,
			})
		}
	}

	allAvailable := true
	for _, eq := range equipmentAvailability {
		if eq["required"].(bool) && !eq["available"].(bool) {
			allAvailable = false
			break
		}
	}
	if allAvailable {
		for _, mat := range materialAvailability {
			if !mat["sufficient"].(bool) {
				allAvailable = false
				break
			}
		}
	}

	render.JSON(w, r, map[string]interface{}{
		"routineId":   routineID,
		"routineName": routine.Name,
		"available":   allAvailable,
		"equipment":   equipmentAvailability,
		"materials":   materialAvailability,
	})
}

// GET /api/routines/upcoming - Get upcoming routines based on schedules
func GetUpcomingRoutines(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	db := database.GetDB()

	// Get routines with deadlines or recurrence
	var routines []models.Routine
	now := time.Now()
	if err := db.Where("laboratory_id = ? AND (deadline IS NOT NULL OR schedule_type = 'recurring')", claims.LaboratoryID).
		Preload("Recurrence").
		Preload("AssignedTo").
		Order("deadline ASC NULLS LAST").
		Find(&routines).Error; err != nil {
		log.Printf("[ROUTINES] Get upcoming error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch upcoming routines")
		return
	}

	type UpcomingRoutine struct {
		ID           uint       `json:"id"`
		Name         string     `json:"name"`
		ScheduleType string     `json:"scheduleType"`
		NextDue      *time.Time `json:"nextDue,omitempty"`
		Deadline     *time.Time `json:"deadline,omitempty"`
		AssignedTo   []uint     `json:"assignedTo"`
	}

	upcoming := make([]UpcomingRoutine, 0)
	for _, routine := range routines {
		ur := UpcomingRoutine{
			ID:           routine.ID,
			Name:         routine.Name,
			ScheduleType: routine.ScheduleType,
			Deadline:     routine.Deadline,
			AssignedTo:   []uint{},
		}

		// For one-time routines, use deadline
		if routine.ScheduleType == "one_time" && routine.Deadline != nil && routine.Deadline.After(now) {
			ur.NextDue = routine.Deadline
		}

		// For recurring routines, calculate next occurrence
		if routine.ScheduleType == "recurring" && routine.Recurrence != nil {
			// Simple calculation - for production, use proper recurrence library
			if routine.Recurrence.Frequency == "daily" {
				nextDue := now.AddDate(0, 0, routine.Recurrence.Interval)
				ur.NextDue = &nextDue
			} else if routine.Recurrence.Frequency == "weekly" {
				nextDue := now.AddDate(0, 0, 7*routine.Recurrence.Interval)
				ur.NextDue = &nextDue
			} else if routine.Recurrence.Frequency == "monthly" {
				nextDue := now.AddDate(0, routine.Recurrence.Interval, 0)
				ur.NextDue = &nextDue
			}
		}

		for _, user := range routine.AssignedTo {
			ur.AssignedTo = append(ur.AssignedTo, user.ID)
		}

		if ur.NextDue != nil || ur.Deadline != nil {
			upcoming = append(upcoming, ur)
		}
	}

	render.JSON(w, r, upcoming)
}

// GET /api/routines/executions - Get execution history with filters
func GetExecutionHistory(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	db := database.GetDB()
	query := db.Where("laboratory_id = ?", claims.LaboratoryID)

	// Filter by routine ID
	routineID := r.URL.Query().Get("routineId")
	if routineID != "" {
		query = query.Where("routine_id = ?", routineID)
	}

	// Filter by status
	status := r.URL.Query().Get("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by executed by
	executedBy := r.URL.Query().Get("executedBy")
	if executedBy != "" {
		query = query.Where("executed_by = ?", executedBy)
	}

	// Date range filters
	startDate := r.URL.Query().Get("startDate")
	if startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			query = query.Where("started_at >= ?", t)
		}
	}

	endDate := r.URL.Query().Get("endDate")
	if endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			query = query.Where("started_at <= ?", t)
		}
	}

	var executions []models.RoutineExecution
	if err := query.
		Preload("Routine").
		Preload("Executor").
		Preload("StepCompletions.Step").
		Order("started_at DESC").
		Find(&executions).Error; err != nil {
		log.Printf("[ROUTINES] Get execution history error: %v", err)
		middleware.ErrorResponse(w, r, http.StatusInternalServerError, "failed to fetch executions")
		return
	}

	response := make([]ExecutionResponse, 0, len(executions))
	for _, execution := range executions {
		response = append(response, buildExecutionResponse(execution))
	}

	render.JSON(w, r, response)
}

// GET /api/routines/statistics - Get routine execution statistics
func GetRoutineStatistics(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		middleware.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
		return
	}

	if claims.LaboratoryID == 0 {
		middleware.ErrorResponse(w, r, http.StatusForbidden, "user does not belong to a laboratory")
		return
	}

	db := database.GetDB()

	// Total routines count
	var totalRoutines int64
	db.Model(&models.Routine{}).Where("laboratory_id = ?", claims.LaboratoryID).Count(&totalRoutines)

	// Routines by schedule type
	type ScheduleTypeCount struct {
		ScheduleType string `json:"scheduleType"`
		Count        int64  `json:"count"`
	}
	var scheduleTypeCounts []ScheduleTypeCount
	db.Model(&models.Routine{}).
		Select("schedule_type, COUNT(*) as count").
		Where("laboratory_id = ?", claims.LaboratoryID).
		Group("schedule_type").
		Scan(&scheduleTypeCounts)

	// Total executions
	var totalExecutions int64
	db.Model(&models.RoutineExecution{}).Where("laboratory_id = ?", claims.LaboratoryID).Count(&totalExecutions)

	// Executions by status
	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var statusCounts []StatusCount
	db.Model(&models.RoutineExecution{}).
		Select("status, COUNT(*) as count").
		Where("laboratory_id = ?", claims.LaboratoryID).
		Group("status").
		Scan(&statusCounts)

	// Average completion time for completed executions
	type CompletionTimeResult struct {
		AvgDurationSeconds float64
	}
	var avgResult CompletionTimeResult
	db.Model(&models.RoutineExecution{}).
		Select("AVG(EXTRACT(EPOCH FROM (completed_at - started_at))) as avg_duration_seconds").
		Where("laboratory_id = ? AND status = 'completed' AND completed_at IS NOT NULL", claims.LaboratoryID).
		Scan(&avgResult)

	// Most executed routines (top 5)
	type RoutineExecutionCount struct {
		RoutineID   uint   `json:"routineId"`
		RoutineName string `json:"routineName"`
		Count       int64  `json:"count"`
	}
	var topRoutines []RoutineExecutionCount
	db.Model(&models.RoutineExecution{}).
		Select("routine_executions.routine_id, routines.name as routine_name, COUNT(*) as count").
		Joins("JOIN routines ON routines.id = routine_executions.routine_id").
		Where("routine_executions.laboratory_id = ?", claims.LaboratoryID).
		Group("routine_executions.routine_id, routines.name").
		Order("count DESC").
		Limit(5).
		Scan(&topRoutines)

	// Recent activity (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var recentExecutions int64
	db.Model(&models.RoutineExecution{}).
		Where("laboratory_id = ? AND started_at >= ?", claims.LaboratoryID, sevenDaysAgo).
		Count(&recentExecutions)

	// Completion rate
	var completedCount int64
	db.Model(&models.RoutineExecution{}).
		Where("laboratory_id = ? AND status = 'completed'", claims.LaboratoryID).
		Count(&completedCount)

	completionRate := 0.0
	if totalExecutions > 0 {
		completionRate = float64(completedCount) / float64(totalExecutions) * 100
	}

	render.JSON(w, r, map[string]interface{}{
		"totalRoutines":          totalRoutines,
		"routinesByScheduleType": scheduleTypeCounts,
		"totalExecutions":        totalExecutions,
		"executionsByStatus":     statusCounts,
		"completionRate":         completionRate,
		"avgCompletionTime":      avgResult.AvgDurationSeconds,
		"topRoutines":            topRoutines,
		"recentActivity": map[string]interface{}{
			"last7Days": recentExecutions,
		},
	})
}

// RegisterRoutineRoutes registers all routine-related routes
func RegisterRoutineRoutes(r chi.Router) {
	r.Route("/routines", func(r chi.Router) {
		// CRUD endpoints
		r.Get("/", GetAllRoutines)
		r.Post("/", CreateRoutine)

		// Availability and scheduling endpoints (must come before /{id})
		r.Get("/upcoming", GetUpcomingRoutines)
		r.Get("/executions", GetExecutionHistory)
		r.Get("/statistics", GetRoutineStatistics)

		r.Get("/{id}", GetRoutineByID)
		r.Put("/{id}", UpdateRoutine)
		r.Delete("/{id}", DeleteRoutine)
		r.Get("/{routineId}/availability", GetRoutineAvailability)

		// Execution endpoints
		r.Post("/{routineId}/executions", StartExecution)
		r.Get("/executions/{executionId}", GetExecution)
		r.Patch("/executions/{executionId}/steps/{stepId}", MarkStepComplete)
		r.Post("/executions/{executionId}/complete", CompleteExecution)
		r.Post("/executions/{executionId}/cancel", CancelExecution)
	})
}
