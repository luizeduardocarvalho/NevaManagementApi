package migrations

import (
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
				&models.Researcher{},
				&models.EquipmentUsage{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&models.EquipmentUsage{},
				&models.Researcher{},
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
}

// NewMigrator creates a configured gormigrate instance bound to the provided DB.
func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	opts := *gormigrate.DefaultOptions
	opts.TableName = "schema_migrations"
	return gormigrate.New(db, &opts, migrationList)
}
