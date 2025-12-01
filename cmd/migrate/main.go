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
	services.InitMongoDB()

	// Parse command
	command := flag.String("command", "", "Migration command: migrate, rollback, fresh, status")
	flag.Parse()

	if *command == "" {
		fmt.Println("MongoDB Migration Tool")
		fmt.Println("Note: MongoDB is schema-less and doesn't require SQL migrations.")
		fmt.Println("Collections are created automatically when first used.")
		fmt.Println("\nAvailable commands:")
		fmt.Println("  migrate  - No-op for MongoDB (collections created automatically)")
		fmt.Println("  rollback - Not applicable for MongoDB")
		fmt.Println("  fresh    - Drop all collections and recreate")
		fmt.Println("  status   - Show collection status")
		os.Exit(0)
	}

	switch *command {
	case "migrate":
		fmt.Println("✓ MongoDB is schema-less - no migrations needed")
		fmt.Println("Collections will be created automatically when first used.")

	case "rollback":
		fmt.Println("⚠ Rollback is not applicable for MongoDB (schema-less database)")

	case "fresh":
		if err := services.FreshMongoDB(); err != nil {
			log.Fatalf("Fresh operation failed: %v", err)
		}
		fmt.Println("✓ All collections dropped successfully")

	case "status":
		if err := services.MongoDBStatus(); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: migrate, rollback, fresh, status")
		os.Exit(1)
	}
}

