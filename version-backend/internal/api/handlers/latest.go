package handlers

import (
	"encoding/json"
	"net/http"

	"version-backend/internal/api/middleware"
	"version-backend/internal/db"
)

// LatestDataResponse represents the structure of the response from the /latest_data endpoint
type LatestDataResponse struct {
	OSVersion struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Platform string `json:"platform"`
	} `json:"os_version"`
	OsqueryVersion string    `json:"osquery_version"`
	InstalledApps  []AppInfo `json:"installed_apps"`
	LastUpdated    string    `json:"last_updated"`
}

// AppInfo represents an installed application in the API response
type AppInfo struct {
	Name                 string  `json:"name"`
	Path                 string  `json:"path"`
	BundleIdentifier     string  `json:"bundle_identifier,omitempty"`
	BundleName           string  `json:"bundle_name,omitempty"`
	BundleShortVersion   string  `json:"bundle_short_version,omitempty"`
	DisplayName          string  `json:"display_name,omitempty"`
	MinimumSystemVersion string  `json:"minimum_system_version,omitempty"`
	LastOpenedTime       float64 `json:"last_opened_time,omitempty"`
	EndTime              float64 `json:"end_time,omitempty"`
}

// GetLatestData handles the GET /latest_data endpoint
// It retrieves the most recent system information from the database
// and returns it in JSON format
func GetLatestData(w http.ResponseWriter, r *http.Request) {
	// Get database instance from context using the correct key
	dbInstance, ok := r.Context().Value(middleware.DBKey{}).(*db.DB)
	if !ok {
		http.Error(w, "Database connection not found", http.StatusInternalServerError)
		return
	}

	// Get latest system info from database
	sysInfo, err := dbInstance.GetLatestSystemInfo()
	if err != nil {
		if err.Error() == "no system information available yet - waiting for first osquery data collection" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "initializing",
				"message": "System information is being collected. Please try again in a few seconds.",
			})
			return
		}
		http.Error(w, "Error retrieving system information: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := LatestDataResponse{
		OsqueryVersion: sysInfo.OsqueryVersion,
		LastUpdated:    sysInfo.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Set OS version info
	response.OSVersion.Name = sysInfo.OSName
	response.OSVersion.Version = sysInfo.OSVersion
	response.OSVersion.Platform = sysInfo.OSPlatform

	// Convert installed apps
	response.InstalledApps = make([]AppInfo, len(sysInfo.InstalledApps))
	for i, app := range sysInfo.InstalledApps {
		var endTime float64
		if app.EndTime != nil {
			endTime = float64(app.EndTime.Unix())
		}
		response.InstalledApps[i] = AppInfo{
			Name:                 app.Name,
			Path:                 app.Path,
			BundleIdentifier:     app.BundleIdentifier,
			BundleName:           app.BundleName,
			BundleShortVersion:   app.BundleShortVersion,
			DisplayName:          app.DisplayName,
			MinimumSystemVersion: app.MinimumSystemVersion,
			LastOpenedTime:       app.LastOpenedTime,
			EndTime:              endTime,
		}
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	// Write response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
