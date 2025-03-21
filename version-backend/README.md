# Version Backend

A Go service that integrates with osquery to collect system information, stores it in a MariaDB database, and exposes it through a REST API.

## Features

- **System Information Collection**:
  - OS Version (via osquery's `os_version` table)
  - Osquery Version (via `osquery_info` table)
  - Installed Applications (via `apps` table)
- **Real-time Data Collection**:
  - Initial data collection on startup
  - Configurable periodic updates
  - Change detection and versioning
- **Data Storage**:
  - MariaDB for persistent storage
  - Efficient schema design
  - Data versioning support
- **REST API**:
  - Clean JSON responses
  - Error handling
  - CORS support for frontend integration
- **Monitoring**:
  - Health check endpoint
  - Detailed status information
  - Comprehensive logging

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Osquery installed on your system
- MacOS or Windows (project requirement)

## Quick Start

1. Install osquery:
```bash
# MacOS
brew install osquery

# Windows (using Chocolatey)
choco install osquery
```

2. Start the database:
```bash
docker-compose up -d
```

3. Set up environment variables (create .env file):
```env
SERVER_HOST=localhost
SERVER_PORT=7070
DB_HOST=localhost
DB_PORT=3306
DB_USER=osquery
DB_PASSWORD=osquery_password
DB_NAME=osquery_data
OSQUERY_SOCKET=/var/osquery/osquery.em
QUERY_INTERVAL=300
```

4. Run the application:
```bash
go run cmd/server/main.go
```

## API Endpoints

### GET /api/latest_data

Returns the most recent system information.

Response format:
```json
{
    "os_version": {
        "name": "macOS",
        "version": "14.0.0",
        "platform": "darwin"
    },
    "osquery_version": "5.10.2",
    "installed_apps": [
        {
            "name": "Chrome",
            "path": "/Applications/Google Chrome.app",
            "bundle_identifier": "com.google.Chrome",
            "bundle_name": "Google Chrome",
            "bundle_short_version": "120.0.6099.129",
            "display_name": "Google Chrome",
            "minimum_system_version": "10.13",
            "last_opened_time": 1678901234
        }
    ],
    "last_updated": "2024-03-15T10:30:00Z"
}
```

### GET /health

Basic health check endpoint.

### GET /status

Detailed system status information.

## Project Structure

```
version-backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # HTTP server and handlers
│   ├── config/         # Configuration management
│   ├── db/             # Database operations
│   └── osquery/        # Osquery client
├── pkg/
│   └── logger/         # Logging package
├── scripts/
│   └── init.sql        # Database initialization
└── docker-compose.yml  # Docker configuration
```

## Development

### Building

```bash
go build -o bin/server cmd/server/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/api/...
```

### Database Management

The application uses MariaDB for data storage. Schema migrations are handled through the `init.sql` script.

Access phpMyAdmin:
- URL: http://localhost:6060
- Username: root
- Password: rootpassword

## Monitoring

The application provides several monitoring endpoints:

- `/health` - Basic health check
- `/status` - Detailed system status including:
  - Uptime
  - Memory usage
  - Database connection status
  - Last data collection timestamp

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License 