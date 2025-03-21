package models

import (
	"time"
)

// SystemInfo represents the system information stored in the database
type SystemInfo struct {
	ID             int64     `db:"id"`
	OSName         string    `db:"os_name"`
	OSVersion      string    `db:"os_version"`
	OSPlatform     string    `db:"os_platform"`
	OsqueryVersion string    `db:"osquery_version"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	InstalledApps  []InstalledApp
}

// InstalledApp represents an installed application in the system
type InstalledApp struct {
	ID                   int64      `db:"id"`
	SystemInfoID         int64      `db:"system_info_id"`
	Name                 string     `db:"name"`
	Path                 string     `db:"path"`
	BundleIdentifier     string     `db:"bundle_identifier"`
	BundleName           string     `db:"bundle_name"`
	BundleShortVersion   string     `db:"bundle_short_version"`
	DisplayName          string     `db:"display_name"`
	MinimumSystemVersion string     `db:"minimum_system_version"`
	LastOpenedTime       float64    `db:"last_opened_time"`
	CreatedAt            time.Time  `db:"created_at"`
	EndTime              *time.Time `db:"end_time"`
}

// // ToResponse converts the model to an API response format
// func (s *SystemInfo) ToResponse() map[string]interface{} {
// 	return nil
// }
