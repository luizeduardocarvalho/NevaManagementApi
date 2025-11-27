package migrations

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/models"
)

var migrationList = []*gormigrate.Migration{
	{
		ID: "20250215_initial_schema",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&models.Laboratory{},
				&models.LaboratoryInvitation{},
				&models.User{},
				&models.Location{},
				&models.Product{},
				&models.Equipment{},
				&models.EquipmentUsage{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&models.EquipmentUsage{},
				&models.Equipment{},
				&models.Product{},
				&models.Location{},
				&models.User{},
				&models.LaboratoryInvitation{},
				&models.Laboratory{},
			)
		},
	},
	{
		ID: "20250217_enhanced_invitations",
		Migrate: func(tx *gorm.DB) error {
			// Add new columns to laboratory_invitations table
			if !tx.Migrator().HasColumn(&models.LaboratoryInvitation{}, "first_name") {
				if err := tx.Migrator().AddColumn(&models.LaboratoryInvitation{}, "first_name"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(&models.LaboratoryInvitation{}, "last_name") {
				if err := tx.Migrator().AddColumn(&models.LaboratoryInvitation{}, "last_name"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(&models.LaboratoryInvitation{}, "invited_by") {
				if err := tx.Migrator().AddColumn(&models.LaboratoryInvitation{}, "invited_by"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(&models.LaboratoryInvitation{}, "status") {
				if err := tx.Migrator().AddColumn(&models.LaboratoryInvitation{}, "status"); err != nil {
					return err
				}
			}

			// Migrate is_accepted to status
			if tx.Migrator().HasColumn(&models.LaboratoryInvitation{}, "is_accepted") {
				// Update existing records: is_accepted = true -> status = 'accepted', else 'pending'
				if err := tx.Exec(`UPDATE laboratory_invitations SET status = CASE WHEN is_accepted = true THEN 'accepted' ELSE 'pending' END WHERE status IS NULL OR status = ''`).Error; err != nil {
					return err
				}
				// Drop the old is_accepted column
				if err := tx.Migrator().DropColumn(&models.LaboratoryInvitation{}, "is_accepted"); err != nil {
					return err
				}
			}

			// Add status column to users table
			if !tx.Migrator().HasColumn(&models.User{}, "status") {
				if err := tx.Migrator().AddColumn(&models.User{}, "status"); err != nil {
					return err
				}
				// Set default status to 'active' for existing users
				if err := tx.Exec(`UPDATE users SET status = 'active' WHERE status IS NULL OR status = ''`).Error; err != nil {
					return err
				}
			}

			// Create indexes manually
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_lab_email_status ON laboratory_invitations(laboratory_id, email, status)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_email ON laboratory_invitations(email)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_token ON laboratory_invitations(invitation_token)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_status ON laboratory_invitations(status)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_invitations_expires ON laboratory_invitations(expires_at)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`)

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Add back is_accepted column
			type TempInvitation struct {
				IsAccepted bool `gorm:"default:false"`
			}
			if err := tx.Migrator().AddColumn(&models.LaboratoryInvitation{}, "is_accepted"); err != nil {
				return err
			}

			// Migrate status back to is_accepted
			if err := tx.Exec(`UPDATE laboratory_invitations SET is_accepted = (status = 'accepted')`).Error; err != nil {
				return err
			}

			// Remove new columns
			if err := tx.Migrator().DropColumn(&models.LaboratoryInvitation{}, "status"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&models.LaboratoryInvitation{}, "invited_by"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&models.LaboratoryInvitation{}, "first_name"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&models.LaboratoryInvitation{}, "last_name"); err != nil {
				return err
			}

			// Remove status from users
			if err := tx.Migrator().DropColumn(&models.User{}, "status"); err != nil {
				return err
			}

			return nil
		},
	},
	{
		ID: "20250223_add_product_usages",
		Migrate: func(tx *gorm.DB) error {
			// Create product_usages table manually
			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS product_usages (
					id BIGSERIAL PRIMARY KEY,
					product_id BIGINT NOT NULL,
					user_id BIGINT NOT NULL,
					quantity_used DOUBLE PRECISION NOT NULL,
					unit VARCHAR(50) NOT NULL,
					used_at TIMESTAMP NOT NULL,
					notes TEXT,
					laboratory_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (product_id) REFERENCES products(id),
					FOREIGN KEY (user_id) REFERENCES users(id),
					FOREIGN KEY (laboratory_id) REFERENCES laboratories(id)
				);
				CREATE INDEX IF NOT EXISTS idx_product_usages_product_id ON product_usages(product_id);
				CREATE INDEX IF NOT EXISTS idx_product_usages_used_at ON product_usages(used_at);
				CREATE INDEX IF NOT EXISTS idx_product_usages_laboratory_id ON product_usages(laboratory_id);
				CREATE INDEX IF NOT EXISTS idx_product_usages_deleted_at ON product_usages(deleted_at);
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS product_usages`).Error
		},
	},
	{
		ID: "20251124_add_organization_hierarchy",
		Migrate: func(tx *gorm.DB) error {
			// Step 1: Create organizations table
			if err := tx.AutoMigrate(&models.Organization{}); err != nil {
				return fmt.Errorf("failed to create organizations table: %w", err)
			}

			// Step 2: Create default organization for each existing lab
			type TempLab struct {
				ID   uint
				Name string
			}
			var labs []TempLab
			if err := tx.Table("laboratories").Select("id, name").Where("deleted_at IS NULL").Find(&labs).Error; err != nil {
				return fmt.Errorf("failed to fetch laboratories: %w", err)
			}

			// Map to store lab_id -> org_id
			labToOrgMap := make(map[uint]uint)

			for _, lab := range labs {
				org := models.Organization{
					Name:        fmt.Sprintf("%s Organization", lab.Name),
					Description: "Auto-created organization during hierarchy migration",
				}
				if err := tx.Create(&org).Error; err != nil {
					return fmt.Errorf("failed to create organization for lab %d: %w", lab.ID, err)
				}
				labToOrgMap[lab.ID] = org.ID
			}

			// Step 3: Add organization_id column to laboratories (nullable first)
			if !tx.Migrator().HasColumn(&models.Laboratory{}, "organization_id") {
				if err := tx.Exec(`ALTER TABLE laboratories ADD COLUMN organization_id BIGINT`).Error; err != nil {
					return fmt.Errorf("failed to add organization_id to laboratories: %w", err)
				}
			}

			// Step 4: Populate organization_id for existing labs
			for labID, orgID := range labToOrgMap {
				if err := tx.Exec("UPDATE laboratories SET organization_id = ? WHERE id = ?", orgID, labID).Error; err != nil {
					return fmt.Errorf("failed to update laboratory %d with organization_id: %w", labID, err)
				}
			}

			// Step 5: Make organization_id NOT NULL and add foreign key
			if err := tx.Exec(`ALTER TABLE laboratories ALTER COLUMN organization_id SET NOT NULL`).Error; err != nil {
				return fmt.Errorf("failed to set organization_id as NOT NULL: %w", err)
			}

			// Step 6: Add organization_id column to users (nullable, as org users won't have lab_id)
			if !tx.Migrator().HasColumn(&models.User{}, "organization_id") {
				if err := tx.Exec(`ALTER TABLE users ADD COLUMN organization_id BIGINT`).Error; err != nil {
					return fmt.Errorf("failed to add organization_id to users: %w", err)
				}
			}

			// Step 7: Rename 'coordinator' to 'lab-coordinator' in users
			if err := tx.Exec(`UPDATE users SET role = 'lab-coordinator' WHERE role = 'coordinator'`).Error; err != nil {
				return fmt.Errorf("failed to rename coordinator role in users: %w", err)
			}

			// Step 8: Rename 'coordinator' to 'lab-coordinator' in invitations
			if err := tx.Exec(`UPDATE laboratory_invitations SET role = 'lab-coordinator' WHERE role = 'coordinator'`).Error; err != nil {
				return fmt.Errorf("failed to rename coordinator role in invitations: %w", err)
			}

			// Step 9: Create indexes
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_laboratories_organization_id ON laboratories(organization_id)`)
			tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id)`)

			// Step 10: Add foreign key constraints
			tx.Exec(`ALTER TABLE laboratories ADD CONSTRAINT fk_laboratories_organization FOREIGN KEY (organization_id) REFERENCES organizations(id)`)
			tx.Exec(`ALTER TABLE users ADD CONSTRAINT fk_users_organization FOREIGN KEY (organization_id) REFERENCES organizations(id)`)

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Revert role changes
			tx.Exec(`UPDATE users SET role = 'coordinator' WHERE role = 'lab-coordinator'`)
			tx.Exec(`UPDATE laboratory_invitations SET role = 'coordinator' WHERE role = 'lab-coordinator'`)

			// Drop foreign key constraints
			tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_organization`)
			tx.Exec(`ALTER TABLE laboratories DROP CONSTRAINT IF EXISTS fk_laboratories_organization`)

			// Remove indexes
			tx.Exec(`DROP INDEX IF EXISTS idx_users_organization_id`)
			tx.Exec(`DROP INDEX IF EXISTS idx_laboratories_organization_id`)

			// Remove columns
			if err := tx.Migrator().DropColumn(&models.Laboratory{}, "organization_id"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&models.User{}, "organization_id"); err != nil {
				return err
			}

			// Drop organizations table
			return tx.Migrator().DropTable(&models.Organization{})
		},
	},
	{
		ID: "20251124_remove_researcher_table",
		Migrate: func(tx *gorm.DB) error {
			// Step 1: Check if researchers table exists
			if tx.Migrator().HasTable("researchers") {
				// Step 2: Add user_id column to equipment_usages (nullable first)
				if !tx.Migrator().HasColumn(&models.EquipmentUsage{}, "user_id") {
					if err := tx.Exec(`ALTER TABLE equipment_usages ADD COLUMN user_id BIGINT`).Error; err != nil {
						return fmt.Errorf("failed to add user_id column: %w", err)
					}
				}

				// Step 3: Migrate data from researcher_id to user_id
				// For existing data, we need to find matching users by clerk_user_id and email
				if err := tx.Exec(`
					UPDATE equipment_usages eu
					SET user_id = u.id
					FROM researchers r
					JOIN users u ON r.clerk_user_id = u.clerk_user_id AND r.email = u.email
					WHERE eu.researcher_id = r.id
				`).Error; err != nil {
					return fmt.Errorf("failed to migrate researcher_id to user_id: %w", err)
				}

				// Step 4: Drop the researcher_id column
				if tx.Migrator().HasColumn(&models.EquipmentUsage{}, "researcher_id") {
					if err := tx.Migrator().DropColumn(&models.EquipmentUsage{}, "researcher_id"); err != nil {
						return fmt.Errorf("failed to drop researcher_id column: %w", err)
					}
				}

				// Step 5: Make user_id NOT NULL
				if err := tx.Exec(`ALTER TABLE equipment_usages ALTER COLUMN user_id SET NOT NULL`).Error; err != nil {
					return fmt.Errorf("failed to set user_id as NOT NULL: %w", err)
				}

				// Step 6: Add foreign key constraint
				if err := tx.Exec(`ALTER TABLE equipment_usages ADD CONSTRAINT fk_equipment_usages_user FOREIGN KEY (user_id) REFERENCES users(id)`).Error; err != nil {
					return fmt.Errorf("failed to add foreign key: %w", err)
				}

				// Step 7: Drop researchers table
				if err := tx.Migrator().DropTable("researchers"); err != nil {
					return fmt.Errorf("failed to drop researchers table: %w", err)
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// This rollback is complex and may result in data loss
			// For now, we'll just log a warning
			return fmt.Errorf("rollback not supported for removing researcher table - manual intervention required")
		},
	},
	{
		ID: "20251124_remove_admin_role",
		Migrate: func(tx *gorm.DB) error {
			// Step 1: Rename 'Admin' to 'lab-coordinator' in users table
			if err := tx.Exec(`UPDATE users SET role = 'lab-coordinator' WHERE role = 'Admin'`).Error; err != nil {
				return fmt.Errorf("failed to rename Admin role in users: %w", err)
			}

			// Step 2: Drop old role constraint and add new one without Admin
			if err := tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check`).Error; err != nil {
				return fmt.Errorf("failed to drop old role constraint: %w", err)
			}

			if err := tx.Exec(`ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('org-coordinator', 'lab-coordinator', 'technician', 'student'))`).Error; err != nil {
				return fmt.Errorf("failed to add new role constraint: %w", err)
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Revert role changes
			tx.Exec(`UPDATE users SET role = 'Admin' WHERE role = 'lab-coordinator'`)

			// Add back Admin to constraint
			tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check`)
			tx.Exec(`ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('org-coordinator', 'lab-coordinator', 'technician', 'student', 'Admin'))`)

			return nil
		},
	},
	{
		ID: "20251125_add_samples_and_replicas",
		Migrate: func(tx *gorm.DB) error {
			// Create samples table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS samples (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					origin VARCHAR(255),
					isolation_date TIMESTAMP,
					latitude DOUBLE PRECISION,
					longitude DOUBLE PRECISION,
					subculture_medium VARCHAR(255),
					subculture_interval_days INTEGER,
					researcher_id BIGINT,
					location_id BIGINT,
					tags JSONB,
					laboratory_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (researcher_id) REFERENCES users(id),
					FOREIGN KEY (location_id) REFERENCES locations(id),
					FOREIGN KEY (laboratory_id) REFERENCES laboratories(id)
				);
				CREATE INDEX IF NOT EXISTS idx_samples_laboratory_id ON samples(laboratory_id);
				CREATE INDEX IF NOT EXISTS idx_samples_deleted_at ON samples(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create samples table: %w", err)
			}

			// Create replicas table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS replicas (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					sample_id BIGINT NOT NULL,
					location_id BIGINT,
					status VARCHAR(50) DEFAULT 'active',
					last_subculture_date TIMESTAMP,
					next_subculture_date TIMESTAMP,
					laboratory_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (sample_id) REFERENCES samples(id),
					FOREIGN KEY (location_id) REFERENCES locations(id),
					FOREIGN KEY (laboratory_id) REFERENCES laboratories(id)
				);
				CREATE INDEX IF NOT EXISTS idx_replicas_sample_id ON replicas(sample_id);
				CREATE INDEX IF NOT EXISTS idx_replicas_laboratory_id ON replicas(laboratory_id);
				CREATE INDEX IF NOT EXISTS idx_replicas_deleted_at ON replicas(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create replicas table: %w", err)
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP TABLE IF EXISTS replicas`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE IF EXISTS samples`).Error; err != nil {
				return err
			}
			return nil
		},
	},
	{
		ID: "20251126_add_routines",
		Migrate: func(tx *gorm.DB) error {
			// Create routines table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routines (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					schedule_type VARCHAR(50) NOT NULL CHECK (schedule_type IN ('one_time', 'recurring', 'template')),
					deadline TIMESTAMP,
					laboratory_id BIGINT NOT NULL,
					created_by BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (laboratory_id) REFERENCES laboratories(id),
					FOREIGN KEY (created_by) REFERENCES users(id)
				);
				CREATE INDEX IF NOT EXISTS idx_routines_laboratory_id ON routines(laboratory_id);
				CREATE INDEX IF NOT EXISTS idx_routines_deleted_at ON routines(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routines table: %w", err)
			}

			// Create routine_steps table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_steps (
					id BIGSERIAL PRIMARY KEY,
					routine_id BIGINT NOT NULL,
					"order" INTEGER NOT NULL,
					description TEXT NOT NULL,
					notes TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE
				);
				CREATE INDEX IF NOT EXISTS idx_routine_steps_routine_id ON routine_steps(routine_id);
				CREATE INDEX IF NOT EXISTS idx_routine_steps_deleted_at ON routine_steps(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_steps table: %w", err)
			}

			// Create routine_materials table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_materials (
					id BIGSERIAL PRIMARY KEY,
					routine_id BIGINT NOT NULL,
					product_id BIGINT NOT NULL,
					quantity DOUBLE PRECISION NOT NULL,
					unit VARCHAR(50) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE,
					FOREIGN KEY (product_id) REFERENCES products(id)
				);
				CREATE INDEX IF NOT EXISTS idx_routine_materials_routine_id ON routine_materials(routine_id);
				CREATE INDEX IF NOT EXISTS idx_routine_materials_deleted_at ON routine_materials(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_materials table: %w", err)
			}

			// Create routine_equipment table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_equipment (
					id BIGSERIAL PRIMARY KEY,
					routine_id BIGINT NOT NULL,
					equipment_id BIGINT NOT NULL,
					estimated_duration INTEGER,
					required BOOLEAN DEFAULT true,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE,
					FOREIGN KEY (equipment_id) REFERENCES equipment(id)
				);
				CREATE INDEX IF NOT EXISTS idx_routine_equipment_routine_id ON routine_equipment(routine_id);
				CREATE INDEX IF NOT EXISTS idx_routine_equipment_deleted_at ON routine_equipment(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_equipment table: %w", err)
			}

			// Create recurrence_rules table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS recurrence_rules (
					id BIGSERIAL PRIMARY KEY,
					routine_id BIGINT UNIQUE NOT NULL,
					frequency VARCHAR(50) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),
					"interval" INTEGER DEFAULT 1,
					days_of_week TEXT,
					day_of_month INTEGER,
					start_date TIMESTAMP NOT NULL,
					end_date TIMESTAMP,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE
				);
				CREATE INDEX IF NOT EXISTS idx_recurrence_rules_deleted_at ON recurrence_rules(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create recurrence_rules table: %w", err)
			}

			// Create routine_assignments table (many-to-many for assigned users)
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_assignments (
					routine_id BIGINT NOT NULL,
					user_id BIGINT NOT NULL,
					PRIMARY KEY (routine_id, user_id),
					FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_assignments table: %w", err)
			}

			// Create routine_executions table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_executions (
					id BIGSERIAL PRIMARY KEY,
					routine_id BIGINT NOT NULL,
					executed_by BIGINT NOT NULL,
					status VARCHAR(50) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'cancelled')),
					started_at TIMESTAMP NOT NULL,
					completed_at TIMESTAMP,
					notes TEXT,
					laboratory_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (routine_id) REFERENCES routines(id),
					FOREIGN KEY (executed_by) REFERENCES users(id),
					FOREIGN KEY (laboratory_id) REFERENCES laboratories(id)
				);
				CREATE INDEX IF NOT EXISTS idx_routine_executions_routine_id ON routine_executions(routine_id);
				CREATE INDEX IF NOT EXISTS idx_routine_executions_laboratory_id ON routine_executions(laboratory_id);
				CREATE INDEX IF NOT EXISTS idx_routine_executions_deleted_at ON routine_executions(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_executions table: %w", err)
			}

			// Create routine_execution_steps table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_execution_steps (
					id BIGSERIAL PRIMARY KEY,
					execution_id BIGINT NOT NULL,
					step_id BIGINT NOT NULL,
					completed BOOLEAN DEFAULT false,
					completed_at TIMESTAMP,
					notes TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (execution_id) REFERENCES routine_executions(id) ON DELETE CASCADE,
					FOREIGN KEY (step_id) REFERENCES routine_steps(id)
				);
				CREATE INDEX IF NOT EXISTS idx_routine_execution_steps_execution_id ON routine_execution_steps(execution_id);
				CREATE INDEX IF NOT EXISTS idx_routine_execution_steps_deleted_at ON routine_execution_steps(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_execution_steps table: %w", err)
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			tables := []string{
				"routine_execution_steps",
				"routine_executions",
				"routine_assignments",
				"recurrence_rules",
				"routine_equipment",
				"routine_materials",
				"routine_steps",
				"routines",
			}
			for _, table := range tables {
				if err := tx.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s CASCADE`, table)).Error; err != nil {
					return err
				}
			}
			return nil
		},
	},
	{
		ID: "20251127_add_routine_execution_materials",
		Migrate: func(tx *gorm.DB) error {
			// Create routine_execution_materials table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS routine_execution_materials (
					id BIGSERIAL PRIMARY KEY,
					execution_id BIGINT NOT NULL,
					product_id BIGINT NOT NULL,
					planned_quantity DOUBLE PRECISION NOT NULL,
					actual_quantity DOUBLE PRECISION DEFAULT 0,
					unit VARCHAR(50) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (execution_id) REFERENCES routine_executions(id) ON DELETE CASCADE,
					FOREIGN KEY (product_id) REFERENCES products(id)
				);
				CREATE INDEX IF NOT EXISTS idx_routine_execution_materials_execution_id ON routine_execution_materials(execution_id);
				CREATE INDEX IF NOT EXISTS idx_routine_execution_materials_deleted_at ON routine_execution_materials(deleted_at);
			`).Error; err != nil {
				return fmt.Errorf("failed to create routine_execution_materials table: %w", err)
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS routine_execution_materials CASCADE`).Error
		},
	},
	{
		ID: "20251127_add_reserved_quantity_to_products",
		Migrate: func(tx *gorm.DB) error {
			// Add reserved_quantity column to products table
			if !tx.Migrator().HasColumn(&models.Product{}, "reserved_quantity") {
				if err := tx.Exec(`ALTER TABLE products ADD COLUMN reserved_quantity DOUBLE PRECISION DEFAULT 0 NOT NULL`).Error; err != nil {
					return fmt.Errorf("failed to add reserved_quantity column: %w", err)
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(&models.Product{}, "reserved_quantity") {
				return tx.Migrator().DropColumn(&models.Product{}, "reserved_quantity")
			}
			return nil
		},
	},
}

// NewMigrator creates a configured gormigrate instance bound to the provided DB.
func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	opts := *gormigrate.DefaultOptions
	opts.TableName = "schema_migrations"
	return gormigrate.New(db, &opts, migrationList)
}
