package migrate

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Run applies embedded SQL migrations in filename order, tracking progress in
// a simple schema_migrations table. It is intentionally minimal and has no
// down/rollback support.
func Run(ctx context.Context, db *gorm.DB) error {
    // Ensure tracking table exists
    if err := db.WithContext(ctx).Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
        )
    `).Error; err != nil {
        return fmt.Errorf("ensure schema_migrations: %w", err)
    }

    // Read all migration files
    entries, err := migrationsFS.ReadDir("migrations")
    if err != nil {
        return fmt.Errorf("read migrations dir: %w", err)
    }
    // Sort by name to guarantee deterministic order
    sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
            continue
        }
        version := entry.Name()

        // Check if applied
        var count int64
        if err := db.WithContext(ctx).
            Raw(`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, version).
            Scan(&count).Error; err != nil {
            return fmt.Errorf("check migration %s: %w", version, err)
        }
        if count > 0 {
            continue
        }

        // Read SQL
        sqlBytes, err := migrationsFS.ReadFile("migrations/" + version)
        if err != nil {
            return fmt.Errorf("read migration %s: %w", version, err)
        }

        // Apply in a transaction
        if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
            if err := tx.Exec(string(sqlBytes)).Error; err != nil {
                return fmt.Errorf("apply migration %s: %w", version, err)
            }
            if err := tx.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, version).Error; err != nil {
                return fmt.Errorf("record migration %s: %w", version, err)
            }
            return nil
        }); err != nil {
            return err
        }
    }

    return nil
}


