package main

import (
	"flag"
	_ "health/seeders" // Import seeders to register them
	"health/services"
	"log"
)

func main() {
	services.LoadConfig()
	services.InitMongoDB()

	// Parse command
	seederName := flag.String("seeder", "", "Specific seeder name to run (optional, runs all if not specified)")
	flag.Parse()

	if *seederName != "" {
		// Run specific seeder
		if err := services.RunSeeder(*seederName); err != nil {
			log.Fatalf("Seeder failed: %v", err)
		}
		log.Printf("✓ Seeder '%s' completed successfully\n", *seederName)
	} else {
		// Run all seeders
		if err := services.RunSeeders(); err != nil {
			log.Fatalf("Seeders failed: %v", err)
		}
		log.Println("✓ All seeders completed successfully")
	}
}
