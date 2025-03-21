package osquery

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"version-backend/internal/db/models"
	"version-backend/pkg/logger"

	"github.com/osquery/osquery-go"
	"github.com/sirupsen/logrus"
)

// Client represents an osquery client
type Client struct {
	instance *osquery.ExtensionManagerClient
	logger   *logrus.Logger
}

// NewClient creates a new osquery client
func NewClient(socketPath string) (*Client, error) {
	client, err := osquery.NewClient(socketPath, 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error creating osquery client: %w", err)
	}

	return &Client{
		instance: client,
		logger:   logger.GetLogger(),
	}, nil
}

// Close closes the osquery client connection
func (c *Client) Close() {
	if c.instance != nil {
		c.instance.Close()
	}
}

// GetSystemInfo retrieves system information using osquery
func (c *Client) GetSystemInfo(ctx context.Context) (*models.SystemInfo, error) {
	// Get OS version
	osVersionRows, err := c.instance.QueryContext(ctx, Queries.GetOSVersion)
	if err != nil {
		return nil, fmt.Errorf("error querying OS version: %w", err)
	}
	if len(osVersionRows.Response) == 0 {
		return nil, fmt.Errorf("no OS version information found")
	}

	// Get osquery version
	osqueryVersionRows, err := c.instance.QueryContext(ctx, Queries.GetOsqueryVersion)
	if err != nil {
		return nil, fmt.Errorf("error querying osquery version: %w", err)
	}
	if len(osqueryVersionRows.Response) == 0 {
		return nil, fmt.Errorf("no osquery version information found")
	}

	// Get installed applications
	appsRows, err := c.instance.QueryContext(ctx, Queries.GetInstalledApps)
	if err != nil {
		return nil, fmt.Errorf("error querying installed applications: %w", err)
	}

	// Create system info
	sysInfo := &models.SystemInfo{
		OSName:         osVersionRows.Response[0]["name"],
		OSVersion:      osVersionRows.Response[0]["version"],
		OSPlatform:     osVersionRows.Response[0]["platform"],
		OsqueryVersion: osqueryVersionRows.Response[0]["version"],
		InstalledApps:  make([]models.InstalledApp, 0, len(appsRows.Response)),
	}

	// Add installed applications
	for _, row := range appsRows.Response {
		var lastOpenedTime float64
		if timestamp, err := strconv.ParseFloat(row["last_opened_time"], 64); err == nil {
			lastOpenedTime = timestamp
		}

		app := models.InstalledApp{
			Name:                 row["name"],
			Path:                 row["path"],
			BundleIdentifier:     row["bundle_identifier"],
			BundleName:           row["bundle_name"],
			BundleShortVersion:   row["bundle_short_version"],
			DisplayName:          row["display_name"],
			MinimumSystemVersion: row["minimum_system_version"],
			LastOpenedTime:       lastOpenedTime,
		}
		sysInfo.InstalledApps = append(sysInfo.InstalledApps, app)
	}

	return sysInfo, nil
}
