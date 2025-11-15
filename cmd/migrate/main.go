package main

import (
	"fmt"
	"log"
	"os"

	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/migrations"
)

func main() {
	action := "up"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	if err := database.Initialize(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	migrator := migrations.NewMigrator(database.GetDB())

	var err error
	switch action {
	case "up":
		err = migrator.Migrate()
	case "down":
		err = migrator.RollbackLast()
	case "goto":
		if len(os.Args) < 3 {
			log.Fatalf("missing migration ID for goto command")
		}
		err = migrator.MigrateTo(os.Args[2])
	case "redo":
		if err = migrator.RollbackLast(); err == nil {
			err = migrator.Migrate()
		}
	default:
		log.Fatalf("unknown command %q. Supported commands: up, down, goto, redo", action)
	}

	if err != nil {
		log.Fatalf("migration command %q failed: %v", action, err)
	}

	fmt.Printf("migration command %q completed successfully\n", action)
}
