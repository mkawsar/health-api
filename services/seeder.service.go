package services

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Seeder represents a database seeder
type Seeder struct {
	Version int
	Name    string
	Run     func(*gorm.DB) error
}

var seeders []Seeder

// RegisterSeeder registers a seeder function
func RegisterSeeder(version int, name string, run func(*gorm.DB) error) {
	seeders = append(seeders, Seeder{
		Version: version,
		Name:    name,
		Run:     run,
	})
}

// RunSeeders executes all registered seeders
func RunSeeders(db *gorm.DB) error {
	if len(seeders) == 0 {
		log.Println("No seeders registered")
		return nil
	}

	// Sort seeders by version
	sort.Slice(seeders, func(i, j int) bool {
		return seeders[i].Version < seeders[j].Version
	})

	log.Printf("Running %d seeder(s)...\n", len(seeders))

	for _, seeder := range seeders {
		log.Printf("Running seeder %d_%s...", seeder.Version, seeder.Name)

		if err := seeder.Run(db); err != nil {
			return fmt.Errorf("seeder %d_%s failed: %w", seeder.Version, seeder.Name, err)
		}

		log.Printf("Seeder %d_%s completed successfully", seeder.Version, seeder.Name)
	}

	log.Println("All seeders completed successfully")
	return nil
}

// RunSeeder executes a specific seeder by name
func RunSeeder(db *gorm.DB, name string) error {
	for _, seeder := range seeders {
		if seeder.Name == name {
			log.Printf("Running seeder %d_%s...", seeder.Version, seeder.Name)
			if err := seeder.Run(db); err != nil {
				return fmt.Errorf("seeder %d_%s failed: %w", seeder.Version, seeder.Name, err)
			}
			log.Printf("Seeder %d_%s completed successfully", seeder.Version, seeder.Name)
			return nil
		}
	}
	return fmt.Errorf("seeder '%s' not found", name)
}

// GetSeeders returns list of registered seeders
func GetSeeders() []Seeder {
	return seeders
}

// LoadSeedersFromDirectory loads seeder files from a directory
func LoadSeedersFromDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("seeders directory not found: %s", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read seeders directory: %w", err)
	}

	var seederFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			seederFiles = append(seederFiles, entry.Name())
		}
	}

	if len(seederFiles) == 0 {
		log.Println("No seeder files found")
		return nil
	}

	// Note: In a real implementation, you might want to use go plugins or
	// a different approach to dynamically load seeders. For now, seeders
	// should be registered manually in the seeder files.
	log.Printf("Found %d seeder file(s). Make sure they register their seeders using RegisterSeeder()", len(seederFiles))

	return nil
}
