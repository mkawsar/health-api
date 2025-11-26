# Migration & Seeder Commands

This project includes Laravel-style migration and seeder commands for managing your database.

## Migration Commands

### Run Migrations
```bash
make migrate
# or
go run cmd/migrate/main.go -command migrate
```

Runs all pending migrations.

### Rollback Migrations
```bash
# Rollback last migration
make rollback

# Rollback N migrations
make rollback-steps=3
# or
go run cmd/migrate/main.go -command rollback -steps 3
```

### Fresh Migrations
```bash
make fresh
# or
go run cmd/migrate/main.go -command fresh
```

Drops all tables and re-runs all migrations from scratch. **Warning: This will delete all data!**

### Migration Status
```bash
make status
# or
go run cmd/migrate/main.go -command status
```

Shows which migrations have been applied and which are pending.

## Seeder Commands

### Run All Seeders
```bash
make seed
# or
go run cmd/seed/main.go
```

Runs all registered seeders in order.

### Run Specific Seeder
```bash
make seed-specific name=user_seeder
# or
go run cmd/seed/main.go -seeder user_seeder
```

Runs a specific seeder by name.

## Creating New Migrations

1. Create a new SQL file in the `migrations/` directory
2. Follow the naming pattern: `{next_number}_{description}.sql`
3. Include both UP and DOWN sections:

```sql
-- +goose Up
CREATE TABLE example_table (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS example_table;
```

## Creating New Seeders

1. Create a new Go file in the `seeders/` directory
2. Follow the naming pattern: `{number}_{name}_seeder.go`
3. Register the seeder in the `init()` function:

```go
package seeders

import (
    "health/models/db"
    "health/services"
    "gorm.io/gorm"
)

func init() {
    services.RegisterSeeder(1, "my_seeder", func(db *gorm.DB) error {
        // Your seeding logic here
        return nil
    })
}
```

## Examples

### Complete Setup (Fresh Database)
```bash
# Drop all tables and re-run migrations
make fresh

# Seed the database
make seed
```

### Development Workflow
```bash
# Check migration status
make status

# Run new migrations
make migrate

# If something goes wrong, rollback
make rollback
```

## Notes

- Migrations are automatically run on application startup (in `main.go`)
- Seeders must be registered using `RegisterSeeder()` in their `init()` functions
- The `fresh` command will delete all data - use with caution!
- Rollback requires DOWN migrations to be defined in your SQL files

