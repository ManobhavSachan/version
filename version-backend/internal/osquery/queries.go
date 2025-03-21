package osquery

// Queries contains all the SQL queries used to retrieve system information
var Queries = struct {
	// GetOSVersion retrieves the operating system version information
	GetOSVersion string

	// GetOsqueryVersion retrieves the installed osquery version
	GetOsqueryVersion string

	// GetInstalledApps retrieves the list of installed applications
	GetInstalledApps string
}{
	GetOSVersion: `
		SELECT
			name,
			version,
			platform
		FROM os_version
		LIMIT 1;
	`,

	GetOsqueryVersion: `
		SELECT
			version
		FROM osquery_info
		LIMIT 1;
	`,

	GetInstalledApps: `
		SELECT
			name,
			path,
			bundle_identifier,
			bundle_name,
			bundle_short_version,
			display_name,
			minimum_system_version,
			last_opened_time
		FROM apps
		WHERE 
			bundle_identifier IS NOT NULL
			AND path LIKE '/Applications/%'
		ORDER BY last_opened_time DESC;
	`,
}
