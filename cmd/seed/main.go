package main

import (
	"flag"
	"fmt"
	"health/services"
	_ "health/seeders" // Import seeders to register them
	"log"
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
	seederName := flag.String("seeder", "", "Specific seeder name to run (optional, runs all if not specified)")
	flag.Parse()

	if *seederName != "" {
		// Run specific seeder
		if err := services.RunSeeder(services.DB, *seederName); err != nil {
			log.Fatalf("Seeder failed: %v", err)
		}
		fmt.Printf("✓ Seeder '%s' completed successfully\n", *seederName)
	} else {
		// Run all seeders
		if err := services.RunSeeders(services.DB); err != nil {
			log.Fatalf("Seeders failed: %v", err)
		}
		fmt.Println("✓ All seeders completed successfully")
	}
}

