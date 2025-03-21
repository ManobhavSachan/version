package db

import (
	"fmt"

	"version-backend/internal/config"
	"version-backend/internal/db/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB represents the database connection
type DB struct {
	*sqlx.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &DB{db}, nil
}

// SaveSystemInfo saves or updates system information in the database
func (db *DB) SaveSystemInfo(info *models.SystemInfo) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// First, try to find an existing system info record for this snapshot
	var existingID int64
	query := `
		SELECT id FROM system_info 
		WHERE os_name = ? AND os_version = ? AND os_platform = ? AND osquery_version = ?
		ORDER BY created_at DESC LIMIT 1
	`
	err = tx.Get(&existingID, query,
		info.OSName,
		info.OSVersion,
		info.OSPlatform,
		info.OsqueryVersion,
	)

	var systemInfoID int64
	if err == nil {
		// System info exists, check if apps have changed
		systemInfoID = existingID

		// Get existing apps for comparison
		var existingApps []models.InstalledApp
		query = `
			SELECT name, path, bundle_identifier, bundle_name, 
			       bundle_short_version, display_name, minimum_system_version, 
			       last_opened_time
			FROM installed_apps 
			WHERE system_info_id = ?
		`
		if err := tx.Select(&existingApps, query, systemInfoID); err != nil {
			return fmt.Errorf("error getting existing apps: %w", err)
		}

		// Compare apps and only update if there are changes
		if appsHaveChanged(existingApps, info.InstalledApps) {
			// Update system_info timestamp to mark the change
			query = `
				UPDATE system_info 
				SET updated_at = CURRENT_TIMESTAMP
				WHERE id = ?
			`
			if _, err := tx.Exec(query, systemInfoID); err != nil {
				return fmt.Errorf("error updating system info: %w", err)
			}

			// Archive old apps data with end_time
			query = `
				UPDATE installed_apps 
				SET end_time = CURRENT_TIMESTAMP
				WHERE system_info_id = ? AND end_time IS NULL
			`
			if _, err := tx.Exec(query, systemInfoID); err != nil {
				return fmt.Errorf("error archiving old apps: %w", err)
			}

			// Insert new apps as current snapshot
			if err := insertApps(tx, systemInfoID, info.InstalledApps); err != nil {
				return fmt.Errorf("error inserting new apps: %w", err)
			}
		}
	} else {
		// Insert new system info record
		query = `
			INSERT INTO system_info (
				os_name, os_version, os_platform, osquery_version
			) VALUES (?, ?, ?, ?)
		`
		result, err := tx.Exec(query,
			info.OSName,
			info.OSVersion,
			info.OSPlatform,
			info.OsqueryVersion,
		)
		if err != nil {
			return fmt.Errorf("error inserting system info: %w", err)
		}

		systemInfoID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID: %w", err)
		}

		// Insert initial apps snapshot
		if err := insertApps(tx, systemInfoID, info.InstalledApps); err != nil {
			return fmt.Errorf("error inserting initial apps: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// appsHaveChanged compares two sets of apps to detect changes
func appsHaveChanged(existing, new []models.InstalledApp) bool {
	if len(existing) != len(new) {
		return true
	}

	// Create maps for easier comparison
	existingMap := make(map[string]models.InstalledApp)
	for _, app := range existing {
		key := fmt.Sprintf("%s:%s:%s", app.Name, app.Path, app.BundleIdentifier)
		existingMap[key] = app
	}

	// Check if any app has changed
	for _, app := range new {
		key := fmt.Sprintf("%s:%s:%s", app.Name, app.Path, app.BundleIdentifier)
		if existing, ok := existingMap[key]; !ok {
			return true // New app found
		} else {
			// Check if any attributes changed
			if app.BundleName != existing.BundleName ||
				app.BundleShortVersion != existing.BundleShortVersion ||
				app.DisplayName != existing.DisplayName ||
				app.MinimumSystemVersion != existing.MinimumSystemVersion {
				return true
			}
		}
	}
	return false
}

// insertApps handles inserting a batch of apps
func insertApps(tx *sqlx.Tx, systemInfoID int64, apps []models.InstalledApp) error {
	query := `
		INSERT INTO installed_apps (
			system_info_id, name, path, bundle_identifier, 
			bundle_name, bundle_short_version, display_name,
			minimum_system_version, last_opened_time
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, app := range apps {
		_, err := tx.Exec(query,
			systemInfoID,
			app.Name,
			app.Path,
			app.BundleIdentifier,
			app.BundleName,
			app.BundleShortVersion,
			app.DisplayName,
			app.MinimumSystemVersion,
			app.LastOpenedTime,
		)
		if err != nil {
			return fmt.Errorf("error inserting app %s: %w", app.Name, err)
		}
	}
	return nil
}

// GetLatestSystemInfo retrieves the most recent system information
func (db *DB) GetLatestSystemInfo() (*models.SystemInfo, error) {
	var info models.SystemInfo

	// Get latest system info
	query := `
		SELECT 
			id, os_name, os_version, os_platform, osquery_version,
			created_at, updated_at
		FROM system_info
		ORDER BY updated_at DESC, created_at DESC
		LIMIT 1
	`
	if err := db.Get(&info, query); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("no system information available yet - waiting for first osquery data collection")
		}
		return nil, fmt.Errorf("error getting system info: %w", err)
	}

	// Get active installed apps for this system info
	query = `
		SELECT 
			id, system_info_id, name, path, bundle_identifier,
			bundle_name, bundle_short_version, display_name,
			minimum_system_version, last_opened_time, created_at
		FROM installed_apps
		WHERE system_info_id = ? AND end_time IS NULL
		ORDER BY last_opened_time DESC
	`
	if err := db.Select(&info.InstalledApps, query, info.ID); err != nil {
		return nil, fmt.Errorf("error getting installed apps: %w", err)
	}

	return &info, nil
}
