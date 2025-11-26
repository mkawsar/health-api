package main

import (
	"flag"
	"fmt"
	"health/services"
	"log"
	"os"
)

func main() {
	services.LoadConfig()

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		services.Config.MySQLUser,
		services.Config.MySQLPassword,
		services.Config.MySQLHost,
		services.Config.MySQLPort,
		services.Config.MySQLDatabase,
		services.Config.MySQLCharset,
	)

	var err error
	services.DB, err = services.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Parse command
	command := flag.String("command", "", "Migration command: migrate, rollback, fresh, status")
	steps := flag.Int("steps", 1, "Number of migration steps to rollback (for rollback command)")
	flag.Parse()

	if *command == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch *command {
	case "migrate":
		if err := services.RunMigrations(services.DB); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✓ Migrations completed successfully")

	case "rollback":
		if err := services.RollbackMigrations(services.DB, *steps); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Printf("✓ Rolled back %d migration(s) successfully\n", *steps)

	case "fresh":
		if err := services.FreshMigrations(services.DB); err != nil {
			log.Fatalf("Fresh migration failed: %v", err)
		}
		fmt.Println("✓ Database refreshed successfully")

	case "status":
		if err := services.MigrationStatus(services.DB); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: migrate, rollback, fresh, status")
		os.Exit(1)
	}
}

