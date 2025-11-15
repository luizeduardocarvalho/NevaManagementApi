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
}

// NewMigrator creates a configured gormigrate instance bound to the provided DB.
func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	opts := *gormigrate.DefaultOptions
	opts.TableName = "schema_migrations"
	return gormigrate.New(db, &opts, migrationList)
}
