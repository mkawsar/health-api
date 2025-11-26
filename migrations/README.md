# Database Migrations

This directory contains SQL migration files for the database schema.

## Migration Files

Migrations are numbered sequentially and follow the naming pattern:
```
{version}_{description}.sql
```

For example:
- `001_create_users_table.sql`
- `002_create_doctors_table.sql`
- `003_create_notes_table.sql`
- `004_create_tokens_table.sql`

## Migration Format

Each migration file should follow the goose migration format:

```sql
-- +goose Up
-- SQL statements to apply the migration
CREATE TABLE ...

-- +goose Down
-- SQL statements to rollback the migration
DROP TABLE ...
```

The `-- +goose Up` section contains the SQL to apply the migration.
The `-- +goose Down` section contains the SQL to rollback the migration (optional).

## How Migrations Work

1. **Automatic Execution**: Migrations are automatically executed when the application starts via `services.InitPostgreSQL()`

2. **Migration Tracking**: The system uses a `schema_migrations` table to track which migrations have been applied

3. **Idempotent**: Migrations are only applied once. If a migration has already been applied, it will be skipped

4. **Transaction Safety**: Each migration runs in a transaction, so if it fails, it will be rolled back

## Creating New Migrations

1. Create a new SQL file in this directory with the next sequential number
2. Follow the naming pattern: `{next_number}_{description}.sql`
3. Include both `-- +goose Up` and `-- +goose Down` sections
4. Use `IF NOT EXISTS` clauses where appropriate to make migrations idempotent

## Example Migration

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS example_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS example_table;
```

## Manual Migration Execution

Migrations are automatically run on application startup. If you need to run them manually, you can call:

```go
services.RunMigrations(services.DB)
```

## Notes

- Always test migrations on a development database first
- Use `IF NOT EXISTS` for CREATE statements to make migrations safe to re-run
- Include proper indexes and foreign key constraints
- Consider data migration scripts for complex schema changes

