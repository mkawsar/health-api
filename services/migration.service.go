package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// RunMigrations executes all pending migrations
func RunMigrations(db *gorm.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Load all migration files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Execute pending migrations
	for _, migration := range migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d_%s already applied, skipping", migration.Version, migration.Name)
			continue
		}

		log.Printf("Running migration %d_%s...", migration.Version, migration.Name)

		// Execute migration in a transaction
		tx := db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to begin transaction: %w", tx.Error)
		}

		// Execute the migration SQL
		if err := tx.Exec(migration.Up).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		// Record migration as applied
		if err := recordMigration(tx, migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		log.Printf("Migration %d_%s applied successfully", migration.Version, migration.Name)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the schema_migrations table
func createMigrationsTable(db *gorm.DB) error {
	sql := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return db.Exec(sql).Error
}

// getAppliedMigrations returns a map of applied migration versions
func getAppliedMigrations(db *gorm.DB) (map[int]bool, error) {
	type MigrationRecord struct {
		Version int
	}

	var records []MigrationRecord
	if err := db.Table("schema_migrations").Select("version").Find(&records).Error; err != nil {
		return nil, err
	}

	applied := make(map[int]bool)
	for _, record := range records {
		applied[record.Version] = true
	}

	return applied, nil
}

// recordMigration records a migration as applied
func recordMigration(db *gorm.DB, version int, name string) error {
	sql := `INSERT INTO schema_migrations (version, name) VALUES (?, ?)`
	return db.Exec(sql, version, name).Error
}

// loadMigrations loads all migration files from the filesystem
func loadMigrations() ([]Migration, error) {
	migrationsDir := "migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory not found: %s", migrationsDir)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Parse filename: {version}_{name}.sql
		parts := strings.Split(strings.TrimSuffix(entry.Name(), ".sql"), "_")
		if len(parts) < 2 {
			log.Printf("Warning: skipping invalid migration file: %s", entry.Name())
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Warning: skipping migration file with invalid version: %s", entry.Name())
			continue
		}

		name := strings.Join(parts[1:], "_")

		// Read migration file
		filePath := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		// Split into UP and DOWN sections
		sql := string(content)
		upSQL, downSQL := parseMigrationSQL(sql)

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Up:      upSQL,
			Down:    downSQL,
		})
	}

	return migrations, nil
}

// parseMigrationSQL parses migration SQL into UP and DOWN sections
func parseMigrationSQL(sql string) (string, string) {
	// Look for -- +goose Up and -- +goose Down markers, or just use the whole SQL as UP
	// For simplicity, we'll use the whole SQL as UP migration
	// You can enhance this to support DOWN migrations if needed
	upSQL := strings.TrimSpace(sql)
	downSQL := "" // Can be implemented later if rollback is needed

	// Check for goose-style markers
	if strings.Contains(sql, "-- +goose Up") {
		parts := strings.Split(sql, "-- +goose Down")
		if len(parts) > 0 {
			upPart := strings.Split(parts[0], "-- +goose Up")
			if len(upPart) > 1 {
				upSQL = strings.TrimSpace(upPart[1])
			}
		}
		if len(parts) > 1 {
			downSQL = strings.TrimSpace(parts[1])
		}
	}

	return upSQL, downSQL
}

// RollbackMigrations rolls back the last N migrations
func RollbackMigrations(db *gorm.DB, steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	// Get applied migrations in reverse order
	type MigrationRecord struct {
		Version int
		Name    string
	}

	var records []MigrationRecord
	if err := db.Table("schema_migrations").Order("version DESC").Limit(steps).Find(&records).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(records) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Load all migrations to get DOWN SQL
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Create a map for quick lookup
	migrationMap := make(map[int]Migration)
	for _, m := range migrations {
		migrationMap[m.Version] = m
	}

	// Rollback migrations in reverse order
	for _, record := range records {
		migration, exists := migrationMap[record.Version]
		if !exists {
			log.Printf("Warning: migration file for version %d not found, skipping", record.Version)
			continue
		}

		if migration.Down == "" {
			log.Printf("Warning: no DOWN migration for %d_%s, skipping", record.Version, record.Name)
			continue
		}

		log.Printf("Rolling back migration %d_%s...", record.Version, record.Name)

		// Execute rollback in a transaction
		tx := db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to begin transaction: %w", tx.Error)
		}

		// Execute the DOWN migration SQL
		if err := tx.Exec(migration.Down).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute rollback for %d_%s: %w", record.Version, record.Name, err)
		}

		// Remove migration record
		if err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", record.Version).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record: %w", err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit rollback: %w", err)
		}

		log.Printf("Migration %d_%s rolled back successfully", record.Version, record.Name)
	}

	return nil
}

// FreshMigrations drops all tables and re-runs all migrations
func FreshMigrations(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	// Get all table names
	var tables []string
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}

	// Disable foreign key checks
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// Drop all tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		log.Printf("Dropped table: %s", table)
	}

	// Re-enable foreign key checks
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("failed to enable foreign key checks: %w", err)
	}

	log.Println("All tables dropped. Running migrations...")

	// Run all migrations
	return RunMigrations(db)
}

// MigrationStatus shows the status of all migrations
func MigrationStatus(db *gorm.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Load all migration files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	fmt.Println("\nMigration Status:")
	fmt.Println("==================")
	fmt.Printf("%-6s %-40s %-10s\n", "Version", "Name", "Status")
	fmt.Println(strings.Repeat("-", 60))

	for _, migration := range migrations {
		status := "Pending"
		if applied[migration.Version] {
			status = "✓ Applied"
		}
		fmt.Printf("%-6d %-40s %-10s\n", migration.Version, migration.Name, status)
	}

	fmt.Println()
	return nil
}
